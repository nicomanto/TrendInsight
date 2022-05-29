package support

import (
	"github.com/dghubble/oauth1"
)

// getOauth1Token return oauth1 token from access_token and access_token secret
func getOauth1Token(token, tokenSecret string) *oauth1.Token {
	return oauth1.NewToken(token, tokenSecret)
}

// getOauth1Config return oauth1 config from consumer_key and consumer_secret
func getOauth1Config(consumerKey, consumerSecret string) *oauth1.Config {
	return oauth1.NewConfig(consumerKey, consumerSecret)
}
