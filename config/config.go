package config

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/duoflow/yc-snapshot/loggers"
)

// VirtualMachine - struct for VM description
type VirtualMachine struct {
	VMname  string
	VMid    string
	VMhddid string
}

// Configuration - struct to load config
// from a json file
type Configuration struct {
	Token             string
	Folderid          string
	KeyID             string
	ServiceAccountID  string
	PrivateRSAKeyFile string
	StartTime         string
	CleanUpTime       string
	TelegramBotToken  string
}

// ReadConfig - function to read config from yaml file
func ReadConfig(ctx context.Context) (Configuration, []VirtualMachine, error) {
	var configuration Configuration
	vms := make([]VirtualMachine, 0)
	//
	loggers.Info.Printf("ReadConfig() starts")
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()
	// --------
	// Open main config file
	configFile, err1 := os.Open("./ycsd.conf.json")
	defer configFile.Close()
	// Open virtual machines config file
	vmsFile, err2 := os.Open("./ycsd.vms.conf.json")
	defer vmsFile.Close()
	// if we os.Open returns an error then handle it
	if err1 != nil || err2 != nil {
		loggers.Info.Println(err1)
		loggers.Info.Println(err2)
		os.Exit(1)
	}
	loggers.Info.Printf("ReadConfig() Successfully opened сonfig files")
	// read our opened xmlFile as a byte array.
	configValue, _ := ioutil.ReadAll(configFile)
	vmsValue, _ := ioutil.ReadAll(vmsFile)
	// we initialize our Users array
	json.Unmarshal(configValue, &configuration)
	json.Unmarshal(vmsValue, &vms)
	// defer the closing of our jsonFile so that we can parse it later on
	//fmt.Printf("%#v", configuration)
	//fmt.Printf("\r\n\r\n\r\nVMs:\r\n")
	//fmt.Printf("%#v", vms)
	// --------
	return configuration, vms, nil
}
