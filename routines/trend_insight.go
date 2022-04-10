package routines

import (
	"TrendInsight/controllers"
	"TrendInsight/models"
	"context"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// InitTrendInsightRoutine setup trend insight goroutine with the increment of waiting group
func InitTrendInsightRoutine(ctx context.Context, group *sync.WaitGroup, interval time.Duration) {
	group.Add(1)
	go trendInsight(ctx, group, interval)
}

// trendInsight perform trend insight logic (fetch most popular hashtag -> fetch most popular tweet on given hashtag -> post insight)
func trendInsight(ctx context.Context, group *sync.WaitGroup, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer func() {
		group.Done()
		logrus.Warn("Trend insight has stopped")
	}()
	for {
		select {
		case <-ticker.C:
			// get most popular hashtag
			mostTrend, err := controllers.GetMostPopularTrend(controllers.WoieIDItaly, nil)
			if err != nil {
				logrus.Errorln(err)
				continue
			}
			// get most popular tweet
			mostTweet, err := controllers.GetMostTweet(mostTrend.Name, controllers.ResultTypePopular)
			if err != nil {
				logrus.Errorln(err)
				continue
			}
			// create tweet msg
			timestampTweet, err := mostTweet.CreatedAtTime()
			if err != nil {
				logrus.Errorln(err)
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
			if err := controllers.PostTweet(tweetString, nil); err != nil {
				logrus.Errorln(err)
			} else {
				logrus.Infoln("Twitted successfully")
			}
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}
