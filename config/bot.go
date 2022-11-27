package config

import (
	"fmt"
	"strconv"
)

const (
	// bot config key
	BOOT_CONFIG_ROOT        = "/trend_insight/bot/"
	BOT_NEED_TWEET_LANG_KEY = "need_tweet_lang"
	BOT_TWEET_LANG_KEY      = "tweet_lang"
)

type BotConfigParam struct {
	TweetLang     string
	NeedTweetLang bool
}

// NewBotConfig setup bot configuration from parameter store
func NewBotConfig(parameterStore ParameterStore) (*BotConfigParam, error) {
	botClientConfigParam, err := parameterStore.GetParametersByPath(BOOT_CONFIG_ROOT, false, 2)
	if err != nil {
		return nil, err
	}
	config := BotConfigParam{}
	if val, ok := botClientConfigParam[BOT_NEED_TWEET_LANG_KEY]; ok {
		if v, e := strconv.ParseBool(val); e == nil {
			config.NeedTweetLang = v
		} else {
			return nil, fmt.Errorf("cannot parse bool value of %s: %v", BOT_NEED_TWEET_LANG_KEY, e)
		}
	} else {
		return nil, fmt.Errorf("failed to find parameter %s", BOT_NEED_TWEET_LANG_KEY)
	}
	if val, ok := botClientConfigParam[BOT_TWEET_LANG_KEY]; ok {
		config.TweetLang = val
	} else {
		return nil, fmt.Errorf("failed to find parameter %s", BOT_TWEET_LANG_KEY)
	}
	return &config, nil
}
