package main

import (
	"fmt"
	"log"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"encoding/json"
	"os"
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
	user, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}
	log.Printf("User's account: %+v\n", user)
	return client, nil
}
func main() {
	fmt.Println("im a bot")
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

	fmt.Printf("%+v\n", client)
	params := &twitter.UserShowParams{ScreenName: "nobody_stop_me"}
	paras := &twitter.UserTimelineParams{
		ScreenName: "byyourlogic",
		Count:		10,
	}
	user, _, err := client.Users.Show(params)
	println(user.FollowersCount)
	twts, _, err := client.Timelines.UserTimeline(paras)
	println(twts[0].Text)
	ent := twts[0].ExtendedEntities
	if ent == nil {
		log.Println("no media")
		return
	}
	println(ent.Media[0].MediaURL)
}
