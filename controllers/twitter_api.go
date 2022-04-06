package controllers

import (
	"TrendInsight/support"

	"github.com/dghubble/go-twitter/twitter"
)

type SearchResultType string

const (
	ResultTypePopular SearchResultType = "popular"
	ResultTypeMixed   SearchResultType = "mixed"
	ResultTypeRecent  SearchResultType = "recent"
)

// PostTweet post a new tweet with the given message and params
func PostTweet(msg string, params *twitter.StatusUpdateParams) error {
	if _, _, err := support.TwitterClient.Statuses.Update(msg, params); err != nil {
		return err
	}
	return nil
}

// GetMostTrendInsight return the most populare trend by location
func GetMostPopularTrendInsight(woeid int64, params *twitter.TrendsPlaceParams) (*twitter.Trend, error) {
	trendList, _, err := support.TwitterClient.Trends.Place(woeid, params)
	if err != nil {
		return nil, err
	}

	higherVolumes := new(int64)
	var mostPopularTrend twitter.Trend
	for _, trend := range trendList[0].Trends {
		if trend.TweetVolume > *higherVolumes {
			higherVolumes = &trend.TweetVolume
			mostPopularTrend = trend
		}
	}
	return &mostPopularTrend, nil
}

func GetMostTweetByHashtag(query string, resultType SearchResultType) (*twitter.Tweet, error) {
	tweetSearch, _, err := support.TwitterClient.Search.Tweets(&twitter.SearchTweetParams{
		ResultType: string(resultType),
		Query:      query,
	})
	if err != nil {
		return nil, err
	}

	higherLikes := new(int)
	var mostPopularTweet twitter.Tweet
	for _, tweet := range tweetSearch.Statuses {
		if tweet.FavoriteCount > *higherLikes {
			higherLikes = &tweet.FavoriteCount
			mostPopularTweet = tweet
		}
	}
	return &mostPopularTweet, nil
}
