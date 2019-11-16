package disk

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/duoflow/yc-snapshot/config"
)

// Disk - struct for yc disk operations
type Disk struct {
	Token    string
	Folderid string
}

// New - constructor function for Disk
func New(conf config.Configuration) Disk {
	d := Disk{conf.Token, conf.Folderid}
	return d
}

// List - function for listing of all disks
func (d Disk) List(ctx context.Context) {
	log.Println("List of function starts")
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
	log.Println("Function -GetDiskInfo- starts")
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

	fmt.Println(resp.Status)
	fmt.Println(string(respBody))
}
