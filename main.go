package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required env " + name)
	}

	return v
}

func setupEnv() {
	err := godotenv.Load("twitter.env")
	if err != nil {
		panic("Error loading twitter.env file")
	}
}

func main() {
	setupEnv()
	var (
		consumerKey       = getenv("TWITTER_CONSUMER_KEY")
		consumerSecret    = getenv("TWITTER_CONSUMER_SECRET")
		accessToken       = getenv("TWITTER_ACCESS_TOKEN")
		accessTokenSecret = getenv("TWITTER_ACCESS_TOKEN_SECRET")
	)

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)

	log := &logger{logrus.New()}
	api.SetLogger(log)

	stream := api.PublicStreamFilter(url.Values{
		"track": []string{"#golang"},
	})
	defer stream.Stop()

	for v := range stream.C {
		t, ok := v.(anaconda.Tweet)
		if !ok {
			logrus.Warningf("Received unespected value of type %T", v)
			continue
		}

		if t.RetweetedStatus != nil {
			continue
		}

		_, err := api.Retweet(t.Id, false)
		if err != nil {
			logrus.Errorf("could not retweet %d:", t.Id)
		}

		fmt.Printf("%s\n", t.Text)
	}
}

type logger struct {
	*logrus.Logger
}

func (log *logger) Critical(args ...interface{})                 { log.Error(args...) }
func (log *logger) Criticalf(format string, args ...interface{}) { log.Errorf(format, args...) }
func (log *logger) Notice(args ...interface{})                   { log.Info(args...) }
func (log *logger) Noticef(format string, args ...interface{})   { log.Infof(format, args...) }
