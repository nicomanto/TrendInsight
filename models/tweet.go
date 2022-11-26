package models

import (
	"fmt"
	"time"
)

type TrendingTweet struct {
	Hashtag string
	Text    string
	Author  string
	Lang    string
	Likes   int
	Ts      time.Time
	URL     string
}

type PostedTweet struct {
	TrendingTweet
	TsPosted time.Time
}

// PostedTweetToString convert PostedTweet struct to string
func (p *PostedTweet) PostedTweetToString() string {
	return fmt.Sprintf("%s\nTweetino: %s\nAuthor: %s\nLikes: %d\nOriginal tweet: %s\nCreated at: %s",
		p.Hashtag,
		p.Text[:(len(p.Text)/2)]+"...",
		p.Author,
		p.Likes,
		p.URL,
		p.Ts.Format(time.RFC822))
}
