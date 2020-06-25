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

	nick := "nobody_stop_me"
	count := 150
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
	m := limit.Resources.Statuses
	for k, v := range m {
		fmt.Println("Key: ", k, "value: ", v)
	}
	//log.Printf("limit %d\n", limit.Resources.Users["limit"])
	twts, _, err := client.Timelines.UserTimeline(paras)
	for _, twt := range twts {
		if twt.ExtendedEntities != nil {
			for j, u := range twt.ExtendedEntities.Media {
				fmt.Printf("%d %s\n", j, u.MediaURL)
			}
		}
	}
}
