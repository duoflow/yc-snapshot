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

// New - constructor function for Disk
func New(conf config.Configuration) Disk {
	d := Disk{conf.Token, conf.Folderid}
	return d
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
func (d Disk) GetDiskInfo(ctx context.Context, diskid string) {
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
		fmt.Println("Errored when sending request to the server")
		return
	}
	// ---------
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	// parse disk info
	var diskinfo Diskinfo
	parsestatus := json.Unmarshal(respBody, &diskinfo)
	if parsestatus != nil {
		loggers.Error.Printf("Disk GetDiskInfo() Parsing error: %s", parsestatus.Error())
	}
	// ----------
	loggers.Info.Println(diskinfo)
}
