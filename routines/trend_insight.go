package routines

import (
	"TrendInsight/models"
	"TrendInsight/support"
	"context"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// InitTrendInsightRoutine setup trend insight goroutine with the increment of waiting group
func InitTrendInsightRoutine(ctx context.Context, group *sync.WaitGroup, interval time.Duration, userToNotify []string, mostPopularTweetSearchLang *string) {
	logrus.Infoln("Start new trend insight routine")
	group.Add(1)
	go trendInsight(ctx, group, interval, userToNotify, mostPopularTweetSearchLang)
}

// trendInsight perform trend insight logic (fetch most popular hashtag -> fetch most popular tweet on given hashtag -> post insight)
func trendInsight(ctx context.Context, group *sync.WaitGroup, interval time.Duration, userToNotify []string, mostPopularTweetSearchLang *string) {
	ticker := time.NewTicker(interval)
	defer func() {
		group.Done()
		logrus.Warn("Trend insight has stopped")
	}()
	for {
		select {
		case <-ticker.C:
			logrus.Infoln("Trend insight routine run at " + time.Now().UTC().Format(time.RFC822))
			// get most popular hashtag
			mostTrend, err := support.GetMostPopularTrend(support.WoieIDItaly, nil)
			if err != nil {
				logrus.Errorln(err)
				support.SendErrorMail(userToNotify, err.Error())
				continue
			}
			// get most popular tweet
			mostTweet, err := support.GetMostTweet(mostTrend.Name, mostPopularTweetSearchLang, support.ResultTypePopular, true)
			if err != nil {
				logrus.Errorln(err)
				support.SendErrorMail(userToNotify, err.Error())
				continue
			}
			// create tweet msg
			timestampTweet, err := mostTweet.CreatedAtTime()
			if err != nil {
				logrus.Errorln(err)
				support.SendErrorMail(userToNotify, err.Error())
				continue
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
			if err := support.PostTweet(tweetString, nil); err != nil {
				logrus.Errorln(err)
				support.SendErrorMail(userToNotify, err.Error())
			} else {
				logrus.Infoln("Twitted successfully")
			}
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}
