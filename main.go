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
var image_total = 0
var page int
var page_total int
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
	page = 0
	page_total = image_total/4
	if image_total%4 > 0 {page_total++}
	fmt.Printf("images %d pages %d", len(img), page_total)
	grid := fyne.NewContainerWithLayout(layout.NewGridLayout(2),
		img[0], img[1], img[2], img[3])
	w.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		if (string(ev.Name) == "Q") {app.Quit()}
		if (string(ev.Name) == "Space") {
			fmt.Println("1 "+id_list[image_count])
			fmt.Println("2 "+id_list[image_count+1])
			fmt.Println("3 "+id_list[image_count+2])
			fmt.Println("4 "+id_list[image_count+3])
		}
		onPress(ev, w,grid, img)
		if (image_count+8) > (len(img)-1) {
			fmt.Printf("get more tweets. count %d len-1 %d\n", image_count, len(img)-1)
			img, maxid, err = get_twts(client, nick, count, maxid, img)
		}
	})
	w.SetContent(grid)
	w.ShowAndRun()
}
func onPress(ev *fyne.KeyEvent, w fyne.Window, grid *fyne.Container, img []*canvas.Image) {
	fmt.Println("KeyDown: "+string(ev.Name))
	if (ev.Name == "Right") {
			if (image_count+4)<(len(img)-1) {
				image_count+=4
				fmt.Printf("page %d image_count %d\n", page, image_count)
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

	fmt.Printf("Tweets: %d\n", len(twts))
	maxid = twts[0].ID
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
				image_total += 1
			}
		}
	}
	return img, maxid, nil
}
func min(x, y int64) int64 {
 if x < y {
   return x
 }
 return y
}
