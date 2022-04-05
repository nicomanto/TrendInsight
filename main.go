package main

import (
	"TrendInsight/controllers"
	"TrendInsight/support"
	"flag"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func setupConfig(cfgFile string) {
	viper.SetConfigType("json")
	viper.SetConfigName(cfgFile)
	viper.AddConfigPath(".")
	if cfgErr := viper.ReadInConfig(); cfgErr != nil {
		logrus.Fatalln(cfgErr)
	}
}

func main() {
	// add config from config file
	cfgFile := flag.String("cfg", "config", "Specify the viper config file to be used.")
	flag.Parse()
	setupConfig(*cfgFile)

	// create twitter client
	support.SetupOauth1TwitterClient()

	// post tweet
	if err := controllers.PostTweet("Tweet!", nil); err != nil {
		logrus.Errorln(err)
	} else {
		logrus.Infoln("Twitted successfully")
	}
}
