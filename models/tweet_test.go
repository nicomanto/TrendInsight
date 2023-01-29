package models

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCorrectTweetStructure(t *testing.T) {
	r := require.New(t)
	now := time.Now()
	tweet := PostedTweet{
		TrendingTweet: TrendingTweet{
			Hashtag: "dummyHastagh",
			Text:    "dummyText",
			Author:  "dummyAuthor",
			Lang:    "dummyLang",
			Likes:   104,
			Ts:      now,
			URL:     "dummyURL",
		},
		TsPosted: now.Add(time.Hour * 24),
	}
	postedString := tweet.PostedTweetToString()
	likesS := strconv.Itoa(tweet.Likes)
	r.Contains(postedString, tweet.Hashtag)
	r.Contains(postedString, tweet.Text[:(len(tweet.Text)/2)]+"...")
	r.Contains(postedString, tweet.Author)
	r.Contains(postedString, "#TrendInsight")
	r.NotContains(postedString, tweet.Lang)
	r.Contains(postedString, likesS)
	r.Contains(postedString, tweet.Ts.Format(time.RFC822))
	r.Contains(postedString, tweet.URL)
}
