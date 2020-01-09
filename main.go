package main

import (
	"context"
	"io/ioutil"
	"os"
	"os/signal"
	"time"

	"github.com/duoflow/yc-snapshot/config"
	"github.com/duoflow/yc-snapshot/loggers"
	"github.com/duoflow/yc-snapshot/snapshot"
	"github.com/duoflow/yc-snapshot/telegrambot"
	"github.com/duoflow/yc-snapshot/token"
	"github.com/robfig/cron"
)

func main() {
	// make loggers initialization
	loggers.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	//-----------
	ctx := context.Background()
	conf, vms, _ := config.ReadConfig(ctx)
	// init telegram bot
	telegrambot.Initv3(conf.TelegramBotToken)
	// get new IAM token
	token.GetIAMToken(&conf)
	// create
	snap := snapshot.New(&conf, vms)
	//
	c := cron.New()
	// "35 23 */2 * *"
	c.AddFunc(conf.StartTime, func() { loggers.Warning.Printf("Make snapshot f() %s", snap.Folderid) /*snap.MakeSnapshot(ctx)/**/ })
	c.AddFunc(conf.CleanUpTime, func() { loggers.Warning.Printf("Cleanup snapshot f()") /*snap.CleanUpOldSnapshots(ctx) /**/ })
	c.Start()

	// start listening for terminate signals
	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt)
	// run goroutine for listening interrupt signals
	go func() {
		select {
		case sig := <-channel:
			loggers.Trace.Printf("YCSD Aborting. Reason: %s signal was received.\n", sig)
			os.Exit(1)
		}
	}()

	// plan keepalive job
	keepaliveTicker := time.NewTicker(30 * time.Minute)
	for {
		select {
		case <-keepaliveTicker.C:
			loggers.Info.Printf("YCSD Daemon Keepalive - I'm still alive!")
		}
	}
}
