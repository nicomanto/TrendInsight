package controllers

import (
	"TrendInsight/support"

	"github.com/dghubble/go-twitter/twitter"
)

// PostTweet post a new tweet with the given message and params
func PostTweet(msg string, params *twitter.StatusUpdateParams) error {
	if _, _, err := support.TwitterClient.Statuses.Update(msg, params); err != nil {
		return err
	}
	return nil
}
