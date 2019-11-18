package main

import (
	"context"
	"fmt"

	"github.com/duoflow/yc-snapshot/config"
	"github.com/duoflow/yc-snapshot/disk"
	"github.com/duoflow/yc-snapshot/instance"
	"github.com/duoflow/yc-snapshot/token"
)

func main() {
	ctx := context.Background()
	conf, _ := config.ReadConfiguration(ctx)
	fmt.Println(conf)
	// refresh token
	signedtoken := token.GetSignedToken(&conf)
	fmt.Println("SignedToken: " + signedtoken)
	// get new IAM token
	conf.Token = token.GetIAMToken(&conf)
	fmt.Println("New IAMToken: " + conf.Token)
	//
	disks := disk.New(conf)
	disks.GetDiskInfo(ctx, "epdsh9cmbecta6mnfuja")
	//
	vm01 := instance.New(&conf)
	vm01.Get(ctx, "epdnvonl4vd6ik8ngood")

}
