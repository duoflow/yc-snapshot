package disk

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/duoflow/yc-snapshot/config"
	"github.com/duoflow/yc-snapshot/loggers"
)

var (
	// Client - Client for API requests
	Client Disk
)

// Disk - struct for yc disk operations
type Disk struct {
	Token    string
	Folderid string
}

// Diskinfo  - struct for disk info representation
type Diskinfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Size string `json:"size"`
}

// Init - constructor function for Disk
func Init(conf *config.Configuration) {
	Client.Token = conf.Token
	Client.Folderid = conf.Folderid
}

// List - function for listing of all disks
func (d Disk) List(ctx context.Context) {
	loggers.Info.Println("List of function starts")
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()
	// ---------
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://compute.api.cloud.yandex.net/compute/v1/disks", nil)
	// add Auth Header
	req.Header.Add("Authorization", "Bearer "+d.Token)
	// add query params
	q := req.URL.Query()
	q.Add("folderId", d.Folderid)
	req.URL.RawQuery = q.Encode()
	// make request
	resp, err := client.Do(req)
	// ----------
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	// ---------
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(resp.Status)
	fmt.Println(string(respBody))
}

// GetDiskInfo - function for listing of all disks
func (d Disk) GetDiskInfo(ctx context.Context, diskid string) Diskinfo {
	// diskinfo struct
	diskinfo := Diskinfo{"", "", ""}
	loggers.Info.Println("Disk GetDiskInfo() starts")
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()
	// ---------
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://compute.api.cloud.yandex.net/compute/v1/disks/"+diskid, nil)
	// add Auth Header
	req.Header.Add("Authorization", "Bearer "+d.Token)
	// make request
	resp, err := client.Do(req)
	// ----------
	if err != nil {
		loggers.Error.Printf("Errored when sending request to the server: %s", err.Error())
	} else {
		loggers.Info.Printf("Disk GetDiskInfo() Request status = %s", resp.Status)
		respBody, _ := ioutil.ReadAll(resp.Body)
		loggers.Info.Printf(string(respBody))
		// parse disk info
		parsestatus := json.Unmarshal(respBody, &diskinfo)
		if parsestatus != nil {
			loggers.Error.Printf("Disk GetDiskInfo() Parsing error: %s", parsestatus.Error())
		} else {
			loggers.Info.Println("Disk info: ", diskinfo)
		}
	}
	defer resp.Body.Close()
	return diskinfo
}
