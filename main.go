package main

import (
	"fmt"
	"flag"
	"log"
	"time"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"encoding/json"
	"os"
	"io"
	"io/ioutil"
	"net/http"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
)
type Credentials struct {
	ConsumerKey			string
	ConsumerSecret		string
	AccessToken			string
	AccessTokenSecret	string
}

func getClient(creds *Credentials) (*twitter.Client, error) {
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:		twitter.Bool(true),
		IncludeEmail:	twitter.Bool(true),
	}
	_, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}
	return client, nil
}
var image_count = 0
var id_list []string
func main() {
	cfg, err := os.Open("tokens.json")
	if err != nil {
		log.Println(err)
		return
	}
	defer cfg.Close()
	decoder := json.NewDecoder(cfg)
	creds := Credentials{}
	err = decoder.Decode(&creds)
	if err != nil {
		log.Println(err)
		return
	}
	client, err := getClient(&creds)
	if err != nil {
		log.Println("Error getting Twitter Client")
		log.Println(err)
		return
	}

	var nick string
	var count int
	flag.StringVar(&nick, "nick", "gatorsdaily", "twitter nickname to grab")
	flag.IntVar(&count, "count", 40, "number of tweets to get at a time")
	flag.Parse()
	threshold := 12
	var img[]*canvas.Image
	var maxid int64
	img, maxid, err = get_twts(client, nick, count, 0, img)
	fmt.Printf("maxid %d\n", maxid)
	if err != nil {
		log.Println(err)
		return
	}

	app := app.New()
	w := app.NewWindow("fehtgo")
	grid := fyne.NewContainerWithLayout(layout.NewGridLayout(2),
		img[0], img[1], img[2], img[3])
	growing := false
	w.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		if (string(ev.Name) == "Q") {app.Quit()} else 
		if (string(ev.Name) == "Space") {
			fmt.Println("1 "+id_list[image_count])
			fmt.Println("2 "+id_list[image_count+1])
			fmt.Println("3 "+id_list[image_count+2])
			fmt.Println("4 "+id_list[image_count+3])
		} else
		if (string(ev.Name) == "R") {
			fmt.Println("manual get more tweets at image_count %d len-1 %d", image_count, len(img)-1)
			go func() {img, maxid, err = get_twts(client, nick, count, maxid, img)}() //I GUESS
		} else {
			onPress(ev, w,grid, img)
		}
		if (image_count+threshold) > (len(img)-1) {
			fmt.Println("at the threshold")
			if(!growing) {
				fmt.Printf("threshold: auto get more tweets at image_count %d len-1 %d\n", image_count, len(img)-1)
				growing = true
				go func() {
					img, maxid, err = get_twts(client, nick, count, maxid, img)
					growing = false
				}() //I GUESS
			}
		}
	})
	w.SetContent(grid)
	w.ShowAndRun()
}
func onPress(ev *fyne.KeyEvent, w fyne.Window, grid *fyne.Container, img []*canvas.Image) {
	if (ev.Name == "Right") {
			if (image_count+4)<(len(img)-1) {
				image_count+=4
				fmt.Printf("image_count %d\n", image_count)
				if (image_count+4)<(len(img)) {
					grid = fyne.NewContainerWithLayout(layout.NewGridLayout(2),
						img[image_count], img[image_count+1], img[image_count+2], img[image_count+3])
				} else {
					rem := len(img)-image_count
					grid = fyne.NewContainerWithLayout(layout.NewGridLayout(2),
							img[image_count])
					for i:=1;i < rem; i++ {
						fmt.Println(image_count+i)
						grid.AddObject(img[image_count+i])
					}
				}
				w.SetContent(grid)
			}
	} else if (ev.Name == "Left") {
		image_count-=4
		fmt.Printf("image_count %d\n", image_count)
		if (image_count < 0) {image_count = 0}
		grid = fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			img[image_count], img[image_count+1], img[image_count+2], img[image_count+3])
		w.SetContent(grid)
	}
}
func get_twts(client *twitter.Client, nick string, count int, maxid int64, img []*canvas.Image)([]*canvas.Image,int64, error){

		paras := &twitter.UserTimelineParams{
			ScreenName: nick,
			Count:		count,
		}
	if maxid != 0 {
		paras = &twitter.UserTimelineParams{
			ScreenName: nick,
			Count:		count,
			MaxID:		maxid,
		}
	}

	rate_paras := &twitter.RateLimitParams{
		Resources: []string{"statuses"},
	}
	limit, _, err := client.RateLimits.Status(rate_paras)
	if err != nil {
		log.Println("error getting rate limit")
		log.Println(err)
		return img, -1, err
	}
	t := limit.Resources.Statuses["/statuses/user_timeline"]
	fmt.Printf("remaining: %d\n", t.Remaining)
	if t.Remaining == 0 {
		log.Println("Resources are limited")
		log.Println("Reset at", time.Unix(int64(t.Reset),0).String())
	}
	twts, _, err := client.Timelines.UserTimeline(paras)

	maxid = twts[0].ID
	i := 0
	for _, twt := range twts {
		maxid = min(maxid, twt.ID)
		if twt.ExtendedEntities != nil {
			for _, u := range twt.ExtendedEntities.Media {
				resp, e := http.Get(u.MediaURL)
				if e != nil {log.Fatal(e)}
				defer resp.Body.Close()
				f, err := ioutil.TempFile("", "twimg.jpg")
				if e != nil {log.Fatal(e)}
				defer f.Close()
				_,e = io.Copy(f, resp.Body)
				if e != nil {log.Fatal(err)}
				img = append(img, canvas.NewImageFromFile(f.Name()))
				img[len(img)-1].FillMode = canvas.ImageFillContain
				s := fmt.Sprintf("https://twitter.com/%s/status/%d", nick, twt.ID)
				id_list = append(id_list, s)
				i+=1
			}
		}
	}
	fmt.Printf("Got %d more images for total %d\n", i, len(img))
	return img, maxid, nil
}
func min(x, y int64) int64 {
 if x < y {
   return x
 }
 return y
}
