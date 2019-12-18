package snapshot

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"regexp"

	"github.com/duoflow/yc-snapshot/config"
	"github.com/duoflow/yc-snapshot/instance"
)

// Snapshot - struct for yc snapshot operations
type Snapshot struct {
	Token    string
	Folderid string
	vms      []config.VirtualMachine
	instance instance.Instance
}

// New - constructor function for Snapshot
func New(conf *config.Configuration, vms []config.VirtualMachine) Snapshot {
	snap := Snapshot{conf.Token, conf.Folderid, vms, instance.New(conf)}
	return snap
}

// List - function for listing of all Snapshots
func (snap Snapshot) List(ctx context.Context) {
	log.Println("Snapshot List() starts")
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()
	// ---------
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://compute.api.cloud.yandex.net/compute/v1/snapshots", nil)
	// add Auth Header
	req.Header.Add("Authorization", "Bearer "+snap.Token)
	// add query params
	q := req.URL.Query()
	q.Add("folderId", snap.Folderid)
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
	//
	log.Printf("Snapshot List() status: %s", resp.Status)
	log.Println(string(respBody))
}

// Get - function for listing of partucular snapshot
func (snap Snapshot) Get(ctx context.Context, snapshotid string) {
	log.Println("Function -Instance -> Get- starts")
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()
	// ---------
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://compute.api.cloud.yandex.net/compute/v1/snapshots/"+snapshotid, nil)
	// add Auth Header
	req.Header.Add("Authorization", "Bearer "+snap.Token)
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
func (snap Snapshot) Create(ctx context.Context, Diskid string, SnapshotName string, SnapshotDesc string) {
	log.Println("Function -Snapshot -> Create- starts")
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()
	// ---------
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://compute.api.cloud.yandex.net/compute/v1/snapshots", nil)
	// add Auth Header
	req.Header.Add("Authorization", "Bearer "+snap.Token)
	// add query params
	q := req.URL.Query()
	q.Add("folderId", snap.Folderid)
	q.Add("diskId", Diskid)
	q.Add("name", SnapshotName)
	q.Add("description", SnapshotDesc)
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
	// log events
	log.Println(resp.Status)
	log.Println(string(respBody))
}

// MakeSnapshot - function for create snapshot
func (snap Snapshot) MakeSnapshot(ctx context.Context) {
	log.Println("MakeSnapshot() starts")
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()
	// ---------
	for _, vm := range snap.vms {
		go func(vmi config.VirtualMachine) {
			log.Printf("MakeSnapshot(): Discovered VMs: VMid=%s", vmi.VMid)
			snap.instance.Get(ctx, vmi.VMid)

		}(vm)
	}
	// ---------
	log.Println("MakeSnapshot() action")
}
