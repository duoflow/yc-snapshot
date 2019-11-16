package config

import (
	"context"
	"fmt"
	"time"

	gonfig "github.com/tkanos/gonfig"
)

// Configuration - struct to load config
// from a json file
type Configuration struct {
	Token    string
	Folderid string
}

// ReadConfiguration - funtion to
// read service config from json file
func ReadConfiguration(ctx context.Context) (Configuration, error) {
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()
	// --------
	configuration := Configuration{}
	err := gonfig.GetConf("./yc-snapshot.conf.json", &configuration)
	if err != nil {
		fmt.Println("Reading of configuration file failed with error: ", err)
		return configuration, err
	}
	return configuration, err
}
