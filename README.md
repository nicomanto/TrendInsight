# TrendInsight
Twitter bot that shows trends insights about popular tweets

## Introduction
You can find bot at this page: [https://twitter.com/insight_trend](https://twitter.com/insight_trend)

Bot posts info of the most popular tweets about the most popular trend. Posted tweets have this format:
```
#DummyHashtag #TrendInsight
Tweet: trucated tweet's text
Author: Dummy
Likes: 2513
Original tweet: tweet's link
Created at: 28 May 22 21:05 UTC
```

## Folders structure
- **models**: models that bot used
- **config**: services configurations (bot,gomail,parameter_store,twitter client)

## Technology
- [Go](https://go.dev/)
- [Gomail](https://github.com/go-gomail/gomail)
- [Go-twitter](https://github.com/dghubble/go-twitter)
- [OAuth1](https://github.com/dghubble/oauth1)
- [Twitter API documentation](https://developer.twitter.com/en/docs/twitter-api)
- [AWS Lambda](https://docs.aws.amazon.com/lambda/latest/dg/welcome.html)
- [AWS System Manager Parameter Store](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-parameter-store.html)
- [AWS Event Bridge](https://aws.amazon.com/it/eventbridge/)


## Configuration
In order to run bot, you need to set these configurations:
```
{
    "auth":{
        "api_key": "api key for twitter bot",
        "api_key_secret": "api secret key for twitter bot"",
        "bearer_token": "brearer token generated from twitter",
        "access_token":"access token for twitter bot",
        "access_token_secret":"secret access token for twitter bot"
    }, 
    "bot":{
        "trend_insight_post_day_interval": 1, #interval that bot used to post tweets
        "need_most_popular_tweet_lang": true, #if bot search for most popular tweet based on language
        "most_popular_tweet_lang": "it" #language used if need_most_popular_tweet_lang is true
    },
    "mail":{
        "sender": "dummyEmail", # email that send errors reports
        "sender_pwd": "pwd", # email's password
        "recipiens": ["dummyEmail"] # recipiens that will receive errors reports
    }
}
```
For `most_popular_tweet_lang` see [https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes](https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes)

## AWS setup
Bot is a **Lambda** function that is invoked by a **crono scheduler in Event Bridge service**. In order to run lambda must have access to **System Manager** service, in order to pull *twitter client*, *gomail client* and *bot* **configurations** from **Parameter Store**.

## Deploy
New lambda source code will be deployed by a *push* on the *main* branch from a GitHub Action (if code pass all tests!). 

GitHub Action required:
 - AWS_ACCESS_KEY_ID
 - AWS_SECRET_ACCESS_KEY
 - AWS_REGION
