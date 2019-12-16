package snapshot

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/duoflow/yc-snapshot/config"
)

// Snapshot - struct for yc snapshot operations
type Snapshot struct {
	Token       string
	Folderid    string
	Diskid      string
	Name        string
	Description string
	Labels      string
}

// New - constructor function for Snapshot
func New(conf *config.Configuration) Snapshot {
	i := Snapshot{conf.Token, conf.Folderid, "", "-" + time.Now().Format("2019-11-22-01-00-51"), "", ""}
	return i
}

// List - function for listing of all Snapshots
func (i Snapshot) List(ctx context.Context) {
	log.Println("Function -Snapshot -> List- starts")
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()
	// ---------
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://compute.api.cloud.yandex.net/compute/v1/snapshots", nil)
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
func (i Snapshot) Get(ctx context.Context, snapshotid string) {
	log.Println("Function -Instance -> Get- starts")
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()
	// ---------
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://compute.api.cloud.yandex.net/compute/v1/snapshots/"+snapshotid, nil)
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

// Create - function for create snapshot
func (i Snapshot) Create(ctx context.Context) {
	log.Println("Function -Snapshot -> Create- starts")
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()
	// ---------
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://compute.api.cloud.yandex.net/compute/v1/snapshots", nil)
	// add Auth Header
	req.Header.Add("Authorization", "Bearer "+i.Token)
	// add query params
	q := req.URL.Query()
	q.Add("folderId", i.Folderid)
	q.Add("diskId", i.Diskid)
	q.Add("name", i.Name)
	q.Add("description", i.Description)
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
