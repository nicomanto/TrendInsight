package main

import (
	"TrendInsight/config"
	"TrendInsight/models"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
)

var twitterClient *config.Twitter
var gomailClient *config.GoMail
var parameterStoreClient *config.ParameterStore
var botConfig *config.BotConfigParam
var userToNotifyInError string

func main() {
	// LAMBDA EXECUTION
	lambda.Start(handleLamdaEvent)
}

func init() {
	var err error
	// setup parameter store client
	parameterStoreClient, err = config.NewParameterStore()
	if err != nil {
		panic(err)
	}
	// setup twitter client
	twitterConfig, err := config.NewTwitterClientConfig(*parameterStoreClient)
	if err != nil {
		panic(err)
	}
	twitterClient = config.NewOauth1TwitterClient(*twitterConfig)
	// setup gomail client
	gomailConfig, err := config.NewMailClientConfig(*parameterStoreClient)
	if err != nil {
		panic(err)
	}
	gomailClient = config.NewSetupMailMsgAndDialer(*gomailConfig)
	// setup bot config
	botConfig, err = config.NewBotConfig(*parameterStoreClient)
	if err != nil {
		panic(err)
	}
	// get error recipient
	userToNotifyInError, err = parameterStoreClient.GetParameterValue(config.MAIL_CONFIG_ROOT+config.MAIL_RECIPIENTS_KEY, true)
	if err != nil {
		panic(err)
	}
}

func handleLamdaEvent() (*models.LambdaResponse, error) {
	// tweet!
	logrus.Infoln("testino")
	logrus.Infoln("TrendInsight run at " + time.Now().UTC().Format(time.RFC822))
	// get most popular hashtag
	mostTrend, err := twitterClient.GetMostPopularTrend(config.WoieIDItaly, nil)
	if err != nil {
		gomailClient.SendErrorMail([]string{userToNotifyInError}, err.Error())
		return nil, err
	}
	// get most popular tweet
	var mostPopularTweetLang *string
	if botConfig.NeedTweetLang {
		mostPopularTweetLang = &botConfig.TweetLang
	}
	mostTweet, err := twitterClient.GetMostTweet(mostTrend.Name, mostPopularTweetLang, config.ResultTypePopular, true)
	if err != nil {
		gomailClient.SendErrorMail([]string{userToNotifyInError}, err.Error())
		return nil, err
	}
	// create tweet msg
	timestampTweet, err := mostTweet.CreatedAtTime()
	if err != nil {
		gomailClient.SendErrorMail([]string{userToNotifyInError}, err.Error())
		return nil, err
	}
	trendName := mostTrend.Name
	if !strings.HasPrefix(trendName, "#") {
		trendName = "#" + trendName
	}
	// create posted tweet
	tweet := models.PostedTweet{
		TrendingTweet: models.TrendingTweet{
			Hashtag: trendName,
			Text:    mostTweet.Text,
			Lang:    mostTweet.Lang,
			Author:  mostTweet.User.Name,
			Ts:      timestampTweet.UTC(),
			Likes:   mostTweet.FavoriteCount,
			URL:     "URL not available",
		},
		TsPosted: time.Now().UTC(),
	}
	// check if URL is present
	if len(mostTweet.Entities.Urls) > 0 {
		tweet.URL = mostTweet.Entities.Urls[0].ExpandedURL
	}
	if tweet.URL == "URL not available" && len(mostTweet.Entities.Media) > 0 {
		tweet.URL = mostTweet.Entities.Media[0].URL
	}
	tweetString := tweet.PostedTweetToString()
	// post tweet
	if err := twitterClient.PostTweet(tweetString, nil); err != nil {
		gomailClient.SendErrorMail([]string{userToNotifyInError}, err.Error())
		return nil, err
	}
	return &models.LambdaResponse{Message: "Twitted successfully", Status: 200}, nil
}
