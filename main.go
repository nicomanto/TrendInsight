package main

import (
	"TrendInsight/routines"
	"TrendInsight/support"
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"time"

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
	support.SetupOauth1TwitterClient(
		viper.GetString("auth.api_key"),
		viper.GetString("auth.api_key_secret"),
		viper.GetString("auth.access_token"),
		viper.GetString("auth.access_token_secret"))
	// setup email sender
	support.SetupMailMsgAndDialer(viper.GetString("mail.sender"), viper.GetString("mail.sender_pwd"))
	// init trend insight go routine
	waitingGroup := &sync.WaitGroup{}
	ctxTrendInsight, ctxTrendInsightCF := context.WithCancel(context.Background())
	routines.InitTrendInsightRoutine(ctxTrendInsight, waitingGroup, viper.GetDuration("trend_insight_post_minute_interval")*time.Minute, viper.GetStringSlice("mail.recipiens"))
	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logrus.Infoln("Shutting down")
	// stop trend insight and wait
	ctxTrendInsightCF()
	waitingGroup.Wait()
}
