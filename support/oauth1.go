package support

import (
	"github.com/dghubble/oauth1"
	"github.com/spf13/viper"
)

// getOauth1Token() return oauth1 token from access_token and access_token secret
func getOauth1Token() *oauth1.Token {
	return oauth1.NewToken(viper.GetString("auth.access_token"), viper.GetString("auth.access_token_secret"))
}

// getOauth1Config() return oauth1 config from api_key and api_key_secret
func getOauth1Config() *oauth1.Config {
	return oauth1.NewConfig(viper.GetString("auth.api_key"), viper.GetString("auth.api_key_secret"))
}
