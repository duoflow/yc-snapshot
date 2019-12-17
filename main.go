package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/duoflow/yc-snapshot/config"
	"github.com/duoflow/yc-snapshot/snapshot"
	"github.com/duoflow/yc-snapshot/token"
	"github.com/robfig/cron"
)

func main() {
	ctx := context.Background()
	conf, vms, _ := config.ReadConfig(ctx)
	//fmt.Println(conf)
	// get new IAM token
	token.GetIAMToken(&conf)
	// create
	snap := snapshot.New(&conf, &vms)
	//
	c := cron.New()
	// "35 23 */2 * *"
	c.AddFunc("*/1 * * * *", func() { snap.List(ctx) })
	c.Start()
	// start listening for terminate signals
	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt)
	// run goroutine for listening interrupt signals
	go func() {
		select {
		case sig := <-channel:
			log.Printf("YCSD Aborting. Reason: %s signal was received.\n", sig)
			os.Exit(1)
		}
	}()

	// plan keepalive job
	keepaliveTicker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-keepaliveTicker.C:
			log.Printf("YCSD Keepalive - I'm still alive!")
		}
	}
}
