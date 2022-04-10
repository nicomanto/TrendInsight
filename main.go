package main

import (
	"TrendInsight/controllers"
	"TrendInsight/models"
	"TrendInsight/support"
	"flag"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func setupConfig(cfgFile string) {
	viper.SetConfigType("json")
	viper.SetConfigName(cfgFile)
	viper.AddConfigPath(".")
	if cfgErr := viper.ReadInConfig(); cfgErr != nil {
		logrus.Fatalln(cfgErr)
	}
}

func main() {
	// add config from config file
	cfgFile := flag.String("cfg", "config", "Specify the viper config file to be used.")
	flag.Parse()
	setupConfig(*cfgFile)
	// create twitter client
	support.SetupOauth1TwitterClient()
	// get most popular hashtag
	mostTrend, err := controllers.GetMostPopularTrend(controllers.WoieIDItaly, nil)
	if err != nil {
		logrus.Errorln(err)
	}
	// get most popular tweet
	mostTweet, err := controllers.GetMostTweet(mostTrend.Name, controllers.ResultTypePopular)
	if err != nil {
		logrus.Errorln(err)
	}
	// create tweet msg
	timestampTweet, err := mostTweet.CreatedAtTime()
	if err != nil {
		logrus.Errorln(err)
	}
	trendName := mostTrend.Name
	if !strings.HasPrefix(trendName, "#") {
		trendName = "#" + trendName
	}
	tweet := models.PostedTweet{
		TrendingTweet: models.TrendingTweet{
			Hashtag: trendName,
			Text:    mostTweet.Text,
			Lang:    mostTweet.Lang,
			Author:  mostTweet.User.Name,
			Ts:      timestampTweet.UTC(),
			Likes:   mostTweet.FavoriteCount,
			URL:     mostTweet.Entities.Urls[0].ExpandedURL,
		},
		TsPosted: time.Now().UTC(),
	}
	tweetString := tweet.PostedTweetToString()
	// post tweet
	if err := controllers.PostTweet(tweetString, nil); err != nil {
		logrus.Errorln(err)
	} else {
		logrus.Infoln("Twitted successfully")
	}
}
