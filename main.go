package main

import (
	"context"
	"fmt"

	"github.com/duoflow/yc-snapshot/config"
	"github.com/duoflow/yc-snapshot/disk"
	"github.com/duoflow/yc-snapshot/instance"
)

func main() {
	ctx := context.Background()
	conf, _ := config.ReadConfiguration(ctx)
	fmt.Println(conf.Token)
	//
	disks := disk.New(conf)
	disks.GetDiskInfo(ctx, "epdsh9cmbecta6mnfuja")
	//
	vm01 := instance.New(conf)
	vm01.List(ctx)
}
