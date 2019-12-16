package config

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	gonfig "github.com/tkanos/gonfig"
)

// VirtualMachine - struct for VM description
type VirtualMachine struct {
	VMname  string
	VMid    string
	VMhddid string
}

// VirtualMachines - struct for VMs array
type VirtualMachines struct {
	VirtualMachines []VirtualMachine
}

// Configuration - struct to load config
// from a json file
type Configuration struct {
	Token             string
	Folderid          string
	KeyID             string
	ServiceAccountID  string
	PrivateRSAKeyFile string
}

// ReadConfiguration - funtion to
// read service config from json file
func ReadConfiguration(ctx context.Context) (Configuration, error) {
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()
	// --------
	configuration := Configuration{}
	err := gonfig.GetConf("./ycsd.conf.json", &configuration)
	if err != nil {
		fmt.Println("Reading of configuration file failed with error: ", err)
		return configuration, err
	}
	return configuration, err
}

// ReadConfigurationV2 - function to read config from yaml file
func ReadConfigurationV2(ctx context.Context) (Configuration, error) {
	var configuration Configuration
	vms := make([]VirtualMachine, 0)
	//
	log.Println("f() ReadConfigurationV2 starts")
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
		log.Println(err1)
		log.Println(err2)
		os.Exit(1)
	}
	log.Println("Successfully opened сonfig files")
	// read our opened xmlFile as a byte array.
	configValue, _ := ioutil.ReadAll(configFile)
	vmsValue, _ := ioutil.ReadAll(vmsFile)
	// we initialize our Users array
	json.Unmarshal(configValue, &configuration)
	json.Unmarshal(vmsValue, &vms)
	// defer the closing of our jsonFile so that we can parse it later on
	fmt.Printf("%#v", configuration)
	fmt.Printf("\r\n\r\n\r\nVMs:\r\n")
	fmt.Printf("%#v", vms)
	// --------
	os.Exit(3)
	return configuration, nil

}
