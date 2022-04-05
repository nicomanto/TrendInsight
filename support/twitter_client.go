package support

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

var TwitterClient *twitter.Client

// SetupOauth1TwitterClient() return twitter client with oauth1 configuration
func SetupOauth1TwitterClient() {
	config := getOauth1Config()
	token := getOauth1Token()
	httpClient := config.Client(oauth1.NoContext, token)
	TwitterClient = twitter.NewClient(httpClient)
}
