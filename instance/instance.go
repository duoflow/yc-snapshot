package instance

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/duoflow/yc-snapshot/config"
)

// Instance - struct for yc disk operations
type Instance struct {
	Token    string
	Folderid string
}

// New - constructor function for Disk
func New(conf config.Configuration) Instance {
	i := Instance{conf.Token, conf.Folderid}
	return i
}

// List - function for listing of all disks
func (i Instance) List(ctx context.Context) {
	log.Println("Function -Instance -> List- starts")
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()
	// ---------
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://compute.api.cloud.yandex.net/compute/v1/instances", nil)
	// add Auth Header
	req.Header.Add("Authorization", "Bearer "+i.Token)
	// add query params
	q := req.URL.Query()
	q.Add("folderId", i.Folderid)
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

// Get - function for listing of all disks
func (i Instance) Get(ctx context.Context, instanceid string) {
	log.Println("Function -Instance -> Get- starts")
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()
	// ---------
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://compute.api.cloud.yandex.net/compute/v1/instances/"+instanceid, nil)
	// add Auth Header
	req.Header.Add("Authorization", "Bearer "+i.Token)
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

// Stop - function for listing of all disks
func (i Instance) Stop(ctx context.Context, instanceid string) {
	log.Println("Function -Instance -> Stop- starts")
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()
	// ---------
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://compute.api.cloud.yandex.net/compute/v1/instances/"+instanceid+":stop", nil)
	// add Auth Header
	req.Header.Add("Authorization", "Bearer "+i.Token)
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

// Start - function for listing of all disks
func (i Instance) Start(ctx context.Context, instanceid string) {
	log.Println("Function -Instance -> Stop- starts")
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()
	// ---------
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://compute.api.cloud.yandex.net/compute/v1/instances/"+instanceid+":start", nil)
	// add Auth Header
	req.Header.Add("Authorization", "Bearer "+i.Token)
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
