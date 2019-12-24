package snapshot

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/duoflow/yc-snapshot/config"
	"github.com/duoflow/yc-snapshot/instance"
	"github.com/duoflow/yc-snapshot/loggers"
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
	loggers.Info.Printf("Snapshot List() starts")
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
	loggers.Info.Printf("Snapshot List() status: %s", resp.Status)
	loggers.Info.Printf(string(respBody))
}

// Get - function for listing of partucular snapshot
func (snap Snapshot) Get(ctx context.Context, snapshotid string) {
	loggers.Info.Printf("Function -Instance -> Get- starts")
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
	loggers.Info.Printf("Snapshot Create() starts")
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
	//q.Add("name", SnapshotName)
	q.Add("description", SnapshotDesc)
	req.URL.RawQuery = q.Encode()
	// make request
	resp, err := client.Do(req)
	// ----------
	if err != nil {
		loggers.Error.Printf("Snapshot Create() Errored when sending request to the server")
		return
	}
	// ---------
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	// log events
	loggers.Info.Printf("Snapshot Create() REST API Request status = %s", resp.Status)
	loggers.Trace.Printf("Snapshot Create() EST API Request body = \n%s", string(respBody))
}

// MakeSnapshot - function for create snapshot
func (snap Snapshot) MakeSnapshot(ctx context.Context) {
	loggers.Info.Printf("MakeSnapshot() starts")
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()
	// ---------
	for _, vm := range snap.vms {
		go func(vmi config.VirtualMachine) {
			loggers.Info.Printf("MakeSnapshot(): Discovered VMs: VMid=%s", vmi.VMid)
			// get vm status
			vmstatus := snap.instance.Get(ctx, vmi.VMid)
			// if VM status = RUNNING shutdown the VM
			if vmstatus == "RUNNING" {
				vmstopstate := snap.StopVM(ctx, vmi.VMid)
				loggers.Info.Printf("MakeSnapshot(): VM stop operation state = %d", vmstopstate)
				//------
				if vmstopstate == 1 {
					loggers.Info.Printf("MakeSnapshot(): Start creating snapshot for VM=%s, Disk=%s", vmi.VMid, vmi.VMhddid)
					snapname := "snapshot-date"
					snapdesc := "snap-description"
					snap.Create(ctx, vmi.VMhddid, snapname, snapdesc)
				}

			} else {
				loggers.Error.Printf("MakeSnapshot(): VM with VMid=%s is not in RUNNING state", vmi.VMid)
				loggers.Error.Printf("MakeSnapshot(): SEND EMAIL NOTIFICATION HERE")
			}

		}(vm)
	}
	// ---------
	loggers.Info.Printf("MakeSnapshot() action")
}

// StopVM - function for create snapshot
func (snap Snapshot) StopVM(ctx context.Context, vmid string) int {
	loggers.Info.Printf("Snapshot StopVM() starts")
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()
	// Call instance stop REST function
	snap.instance.Stop(ctx, vmid)
	// Check status of VM after sleep timer
	loggers.Info.Printf("Snapshot StopVM() Start sleep timer")
	time.Sleep(120 * time.Second)
	loggers.Info.Printf("Snapshot StopVM() Check VM running status")
	vmstatus := snap.instance.Get(ctx, vmid)
	loggers.Info.Printf("Snapshot StopVM() VM status after shutdown = %s", vmstatus)
	if vmstatus == "STOPPED" {
		loggers.Info.Printf("Snapshot StopVM() VM with VMid=%s has stopped in sleep timer", vmid)
		return 1
	}
	// ----
	loggers.Error.Printf("Snapshot StopVM() VM with VMid=%s hasn't stopped in sleep timer", vmid)
	return 0
}
