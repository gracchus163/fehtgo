package main

import (
	"fmt"
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
var page_total int
var image_total int
var page int
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
	//log.Printf("User's account: %+v\n", user)
	return client, nil
}
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

	nick := "gatorsdaily"
	count := 200
	page_total = 7
	image_count := 0
	image_total = 4*page_total
	paras := &twitter.UserTimelineParams{
		ScreenName: nick,
		Count:		count,
	}

	rate_paras := &twitter.RateLimitParams{
		Resources: []string{"statuses"},
	}
	limit, _, err := client.RateLimits.Status(rate_paras)
	if err != nil {
		log.Println("error getting rate limit")
		log.Println(err)
		return
	}
	t := limit.Resources.Statuses["/statuses/user_timeline"]
	fmt.Printf("remaining: %d\n", t.Remaining)
	if t.Remaining == 0 {
		log.Println("Resources are limited")
		log.Println("Reset at", time.Unix(int64(t.Reset),0).String())
	}
	twts, _, err := client.Timelines.UserTimeline(paras)
	img := make([]*canvas.Image, image_total, image_total)
	for _, twt := range twts {
		if twt.ExtendedEntities != nil {
			for _, u := range twt.ExtendedEntities.Media {
				//fmt.Printf("%d %s\n", j, u.MediaURL)
				resp, e := http.Get(u.MediaURL)
				if e != nil {log.Fatal(e)}
				defer resp.Body.Close()
				f, err := ioutil.TempFile("", "twimg.jpg")
				if e != nil {log.Fatal(e)}
				defer f.Close()
				_,e = io.Copy(f, resp.Body)
				if e != nil {log.Fatal(err)}
				img[image_count] = canvas.NewImageFromFile(f.Name())
				img[image_count].FillMode = canvas.ImageFillContain
				image_count += 1
				if image_count >= image_total {
					goto End
				}
			}
		}
	}
	End:
	app := app.New()
	w := app.NewWindow("fehtgo")
	page = 0
	grid := fyne.NewContainerWithLayout(layout.NewGridLayout(2),
		img[0], img[1], img[2], img[3])
	w.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		if (string(ev.Name) == "Q") {app.Quit()}
		onPress(ev, w,grid, img, &page)
	})
	w.SetContent(grid)
	w.ShowAndRun()
}
func onPress(ev *fyne.KeyEvent, w fyne.Window, grid *fyne.Container, img []*canvas.Image, page *int) {
	fmt.Println("KeyDown: "+string(ev.Name))
	if (ev.Name == "Right") {
		*page += 1
		if (*page >= page_total) {*page = page_total-1}
		grid = fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			img[4**page], img[1+(4**page)], img[2+(4**page)], img[3+(4**page)])
		w.SetContent(grid)
	} else if (ev.Name == "Left") {
		*page -= 1
		if (*page < 0) {*page = 0}
		grid = fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			img[4**page], img[1+(4**page)], img[2+(4**page)], img[3+(4**page)])
		w.SetContent(grid)
	}
}
