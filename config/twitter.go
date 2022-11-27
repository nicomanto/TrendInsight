package config

import (
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

const (
	// twitter client key
	TWITTER_CONFIG_ROOT             = "/trend_insight/twitter_client/"
	TWITTER_ACCESS_TOKEN_KEY        = "access_token"
	TWITTER_ACCESS_TOKEN_SECRET_KEY = "access_token_secret"
	TWITTER_API_KEY_KEY             = "api_key"
	TWITTER_API_KEY_SECRET_KEY      = "api_key_secret"
	TWITTER_BEARER_TOKEN_KEY        = "bearer_token"
)

type TwitterClientConfigParam struct {
	AccessToken      string
	AccessTokeSecret string
	APIKey           string
	APIKeySecret     string
	BearerToken      string
}

type SearchResultType string

const (
	ResultTypePopular SearchResultType = "popular"
	ResultTypeMixed   SearchResultType = "mixed"
	ResultTypeRecent  SearchResultType = "recent"
)

type WoeIDType int64

const (
	WoieIDItaly WoeIDType = 23424853
	WoieIDWorld WoeIDType = 1
)

type Twitter struct {
	client *twitter.Client
}

// NewTwitterClientConfig setup twitter client configuration from parameter store
func NewTwitterClientConfig(parameterStore ParameterStore) (*TwitterClientConfigParam, error) {
	twitterClientParams, err := parameterStore.GetParametersByPath(TWITTER_CONFIG_ROOT, true, 5)
	if err != nil {
		return nil, err
	}
	config := TwitterClientConfigParam{}
	if val, ok := twitterClientParams[TWITTER_ACCESS_TOKEN_KEY]; ok {
		config.AccessToken = val
	} else {
		return nil, fmt.Errorf("failed to find parameter %s", TWITTER_ACCESS_TOKEN_KEY)
	}
	if val, ok := twitterClientParams[TWITTER_ACCESS_TOKEN_SECRET_KEY]; ok {
		config.AccessTokeSecret = val
	} else {
		return nil, fmt.Errorf("failed to find parameter %s", TWITTER_ACCESS_TOKEN_SECRET_KEY)
	}
	if val, ok := twitterClientParams[TWITTER_API_KEY_KEY]; ok {
		config.APIKey = val
	} else {
		return nil, fmt.Errorf("failed to find parameter %s", TWITTER_API_KEY_KEY)
	}
	if val, ok := twitterClientParams[TWITTER_API_KEY_SECRET_KEY]; ok {
		config.APIKeySecret = val
	} else {
		return nil, fmt.Errorf("failed to find parameter %s", TWITTER_API_KEY_SECRET_KEY)
	}
	if val, ok := twitterClientParams[TWITTER_BEARER_TOKEN_KEY]; ok {
		config.BearerToken = val
	} else {
		return nil, fmt.Errorf("failed to find parameter %s", TWITTER_BEARER_TOKEN_KEY)
	}
	return &config, nil
}

// NewOauth1TwitterClient setup twitter client with oauth1 configuration
func NewOauth1TwitterClient(cfgParam TwitterClientConfigParam) *Twitter {
	config := oauth1.NewConfig(cfgParam.APIKey, cfgParam.APIKeySecret)
	t := oauth1.NewToken(cfgParam.AccessToken, cfgParam.AccessTokeSecret)
	httpClient := config.Client(oauth1.NoContext, t)
	return &Twitter{
		client: twitter.NewClient(httpClient),
	}
}

// PostTweet post a new tweet with the given message and params
func (tc *Twitter) PostTweet(msg string, params *twitter.StatusUpdateParams) error {
	if _, _, err := tc.client.Statuses.Update(msg, params); err != nil {
		return err
	}
	return nil
}

// GetMostPopularTrendInsight return the most popular trend by location
func (tc *Twitter) GetMostPopularTrend(woeid WoeIDType, params *twitter.TrendsPlaceParams) (*twitter.Trend, error) {
	trendList, _, err := tc.client.Trends.Place(int64(woeid), params)
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

// GetMostTweet return the most liked tweet by popularity, recently or mixed in the given language if needed
func (tc *Twitter) GetMostTweet(query string, resultLang *string, resultType SearchResultType, includeEntities bool) (*twitter.Tweet, error) {
	searchParams := twitter.SearchTweetParams{
		ResultType:      string(resultType),
		Query:           query,
		IncludeEntities: &includeEntities,
	}
	if resultLang != nil {
		searchParams.Lang = *resultLang
	}
	tweetSearch, _, err := tc.client.Search.Tweets(&searchParams)
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
