package snapshot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/duoflow/yc-snapshot/disk"

	"github.com/duoflow/yc-snapshot/config"
	"github.com/duoflow/yc-snapshot/instance"
	"github.com/duoflow/yc-snapshot/loggers"
	"github.com/duoflow/yc-snapshot/telegrambot"
	"github.com/olekukonko/tablewriter"
)

// Status - status for Disk snapshot creation
type Status struct {
	VMname     string
	DiskID     string
	SnapshotID string
	Status     string
	VMstatus   string
}

var (
	// StatusRegister - Register of snapshot creation jobs
	StatusRegister []Status
)

// Snapshot - struct for yc snapshot operations
type Snapshot struct {
	Token    string
	Folderid string
	vms      []config.VirtualMachine
	instance instance.Instance
}

// SnapshotsUnits - one struct for snapshot description
type SnapshotsUnits struct {
	ID           string   `json:"id"`
	FolderID     string   `json:"folderId"`
	CreatedAt    string   `json:"createdAt"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Labels       string   `json:"labels"`
	StorageSize  string   `json:"storageSize"`
	DiskSize     string   `json:"diskSize"`
	ProductIds   []string `json:"productIds"`
	Status       string   `json:"status"`
	SourceDiskID string   `json:"sourceDiskId"`
}

// Snapshotlist - list of all snapshots available
type Snapshotlist struct {
	SnapshotsArray []SnapshotsUnits `json:"snapshots"`
	NextPageToken  string           `json:"nextPageToken"`
}

// New - constructor function for Snapshot
func New(conf *config.Configuration, vms []config.VirtualMachine) Snapshot {
	snap := Snapshot{conf.Token, conf.Folderid, vms, instance.New(conf)}
	return snap
}

// ListSnapshots - function for listing of all Snapshots
func (snap Snapshot) ListSnapshots(ctx context.Context) Snapshotlist {
	var SnapList Snapshotlist
	// ----
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
	q.Add("pageSize", "1000")
	req.URL.RawQuery = q.Encode()
	// make request
	resp, err := client.Do(req)
	// ----------
	if err != nil {
		loggers.Info.Printf("Errored when sending request to the server")
	}
	// ---------
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	//
	loggers.Info.Printf("Snapshot List() Responce status = %s", resp.Status)
	//loggers.Info.Printf(string(respBody))
	// parse responce body
	parsestatus := json.Unmarshal(respBody, &SnapList)
	if parsestatus != nil {
		loggers.Error.Printf("Snapshot List() Parsing error: %s", parsestatus.Error())
	}
	//loggers.Info.Printf("Snapshot List() List parsed info:\n")
	//loggers.Info.Println(SnapList.SnapshotsArray)
	return SnapList
}

// GetSnapStatusByID - function for listing of partucular snapshot
func (snap Snapshot) GetSnapStatusByID(ctx context.Context, snapshotid string) string {
	loggers.Info.Printf("Snapshot GetSnapStatusByID() starts")
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
		loggers.Error.Printf("Snapshot GetSnapStatusByID() Error when sending request to the server: %s", err.Error())
		return "ERROR"
	}
	// ---------
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	// log statistics
	loggers.Info.Printf("Snapshot GetSnapStatusByID() Responce status = %s", resp.Status)
	loggers.Info.Printf("Snapshot GetSnapStatusByID() Responce body: \n%s", string(respBody))
	// read snapshot status
	statusRegexp := regexp.MustCompile(`(?mi)("status"..)"(.*)",`)
	matchResult := statusRegexp.FindStringSubmatch(string(respBody))
	if matchResult != nil {
		loggers.Info.Printf("Instance GetSnapStatusByID() VM status = %s\n", matchResult[2])
		return matchResult[2]
	}
	// return responce status
	return "UNKNOWN"
}

// GetSnapStatusByName - function for listing of partucular snapshot
func (snap Snapshot) GetSnapStatusByName(ctx context.Context, snapshotname string) (status string, snapid string) {
	loggers.Info.Printf("Snapshot GetSnapStatusByName() starts")
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()
	// ---------
	snapshotstatus := "ERROR"
	snapshotid := "ERROR"
	snaplist := snap.ListSnapshots(ctx)
	for _, snapshot := range snaplist.SnapshotsArray {
		if snapshotname == snapshot.Name {
			snapshotstatus = snapshot.Status
			snapshotid = snapshot.ID
		}
	}
	//
	return snapshotstatus, snapshotid
}

// Create - function for create snapshot
func (snap Snapshot) Create(ctx context.Context, Diskid string, SnapName string, SnapshotDesc string) string {
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
	q.Add("name", SnapName)
	q.Add("description", SnapshotDesc)
	req.URL.RawQuery = q.Encode()
	// make request
	resp, err := client.Do(req)
	// ----------
	if err != nil {
		loggers.Error.Printf("Snapshot Create() Error when sending request to the server: %s", err.Error())
		return err.Error()
	}
	// ---------
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	// log events
	loggers.Info.Printf("Snapshot Create() REST API Responce status = %s", resp.Status)
	loggers.Info.Printf("Snapshot Create() REST API Responce body:\n %s", string(respBody))
	return resp.Status
}

// MakeSnapshot - function for create snapshot
func (snap Snapshot) MakeSnapshot(ctx context.Context) {
	loggers.Info.Printf("MakeSnapshot() starts")
	telegrambot.Tgbot.SendMessage("Snapshot tasks have started")
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()
	// ---------
	for index, vm := range snap.vms {
		go func(vmi config.VirtualMachine, i int) {
			loggers.Info.Printf("Index = %d", i)
			loggers.Info.Printf("MakeSnapshot(): Discovered VMs: VMid=%s", vmi.VMid)
			// define register data
			var registerstatusunit Status
			registerstatusunit.VMname = vmi.VMname
			registerstatusunit.DiskID = vmi.VMhddid
			// get vm status
			vmstatus := snap.instance.Get(ctx, vmi.VMid)
			// if VM status = RUNNING shutdown the VM
			if vmstatus == "RUNNING" {
				// define register data
				registerstatusunit.VMstatus = vmstatus
				// stop target VM
				vmstopstate := snap.StopVM(ctx, vmi.VMid)
				loggers.Info.Printf("MakeSnapshot(): VM stop operation state = %s", vmstopstate)
				//------
				if vmstopstate == "STOPPED" {
					// define register data
					registerstatusunit.VMstatus = vmstopstate
					// define time
					t := time.Now()
					// get disk size
					diskinfo := disk.Client.GetDiskInfo(ctx, vmi.VMhddid)
					// snapshot description with timestamp
					date := strings.ToLower(strings.ReplaceAll(t.Format("2006-01-02T15:04:05"), ":", "-"))
					snapdesc := "autosnap" + "-size-" + diskinfo.Size + "-date-" + date
					// snapshot name with timestamp
					snapname := vmi.VMhddid + "-size-" + diskinfo.Size + "-date-" + date
					loggers.Info.Printf("MakeSnapshot(): Start creating snapshot for VM=%s, Disk=%s with snapshot id=%s", vmi.VMid, vmi.VMhddid, snapname)
					// REST API Call to create snapshot
					snapcreateflag := snap.Create(ctx, vmi.VMhddid, snapname, snapdesc)
					// if REST API Call was successfull, then check snapshot status after timeout
					if snapcreateflag == "200 OK" {
						// start cycle to check stapshot status
						awaiting := 0
						timeinterval := 10
						for i := 0; i < 12; i++ {
							awaiting += timeinterval
							loggers.Info.Printf("MakeSnapshot(): Check snapshot status - start timeout. Time=%d", awaiting)
							time.Sleep(10 * time.Minute)
							snapstatus, snapid := snap.GetSnapStatusByName(ctx, snapname)
							loggers.Info.Printf("MakeSnapshot(): Check snapshot id [%s] status = %s. End timeout. Time=%d", snapid, snapstatus, awaiting)
							if snapstatus == "READY" {
								loggers.Info.Printf("MakeSnapshot(): Snapshot status: %s, Time=%d", snapstatus, awaiting)
								// define register data
								registerstatusunit.SnapshotID = snapid
								registerstatusunit.Status = snapstatus
								// if snapshot is ready then start VM
								vmstartstate := snap.StartVM(ctx, vmi.VMid)
								// if VM has started then nothing to do
								if vmstartstate == "RUNNING" {
									loggers.Info.Printf("MakeSnapshot(): VM with ID=%s has started successfully", vmi.VMid)
									// define register data
									registerstatusunit.VMstatus = vmstartstate
								}
								break
							} else {
								loggers.Info.Printf("MakeSnapshot(): Snapshot status: %s, Time=%d", snapstatus, awaiting)
								// define register data
								registerstatusunit.SnapshotID = snapid
								registerstatusunit.Status = snapstatus
							}
						}
					} else {
						loggers.Error.Printf("MakeSnapshot(): Error in Create Snapshot REST API Call: %s", snapcreateflag)
						vmstartstate := snap.StartVM(ctx, vmi.VMid)
						// if VM has started then nothing to do
						if vmstartstate == "RUNNING" {
							loggers.Error.Printf("MakeSnapshot(): VM with ID=%s has started successfully after snapshot fail", vmi.VMid)
						}
						// define register data
						registerstatusunit.SnapshotID = "ERROR:CALL_API"
						registerstatusunit.Status = "ERROR:CALL_API"
					}
				}
			} else {
				loggers.Error.Printf("MakeSnapshot(): VM with VMid=%s is not in RUNNING state", vmi.VMid)
				loggers.Error.Printf("MakeSnapshot(): SEND EMAIL NOTIFICATION HERE")
				// define register data
				registerstatusunit.SnapshotID = "ERROR:VM_NOT_RUNNING"
				registerstatusunit.Status = "ERROR:VM_NOT_RUNNING"
			}
			// add registerstatusunit to register
			StatusRegister = append(StatusRegister, registerstatusunit)
			messagetobot := "Snapshot task result: \n"
			messagetobot += "VMname: " + registerstatusunit.VMname + "\n"
			messagetobot += "VM status: " + registerstatusunit.VMstatus + "\n"
			messagetobot += "DiskID: " + registerstatusunit.DiskID + "\n"
			messagetobot += "SnapshotID: " + registerstatusunit.SnapshotID + "\n"
			messagetobot += "Snapshot status: " + registerstatusunit.Status + "\n"
			telegrambot.Tgbot.SendMessage(messagetobot)
		}(vm, index)
	}
}

// StopVM - function for create snapshot
func (snap Snapshot) StopVM(ctx context.Context, vmid string) string {
	loggers.Info.Printf("Snapshot StopVM() starts")
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()
	// Call instance stop REST function
	snap.instance.Stop(ctx, vmid)
	// Check status of VM after sleep timer
	loggers.Info.Printf("Snapshot StopVM() Start sleep timer")
	time.Sleep(1 * time.Minute)
	loggers.Info.Printf("Snapshot StopVM() Check VM running status")
	vmstatus := snap.instance.Get(ctx, vmid)
	loggers.Info.Printf("Snapshot StopVM() VM status after shutdown = %s", vmstatus)
	if vmstatus == "STOPPED" {
		loggers.Info.Printf("Snapshot StopVM() VM with VMid=%s has stopped in sleep timer", vmid)
		return "STOPPED"
	}
	// ----
	loggers.Error.Printf("Snapshot StopVM() VM with VMid=%s hasn't stopped in sleep timer", vmid)
	return "DO_NOT_STOPPED"
}

// StartVM - function for create snapshot
func (snap Snapshot) StartVM(ctx context.Context, vmid string) string {
	loggers.Info.Printf("Snapshot StartVM() starts")
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()
	// Call instance stop REST function
	snap.instance.Start(ctx, vmid)
	loggers.Info.Printf("Snapshot StartVM() call f(): snap.instance.Start(ctx, vmid)")
	// Check status of VM after sleep timer
	loggers.Info.Printf("Snapshot StartVM() Start sleep timer")
	time.Sleep(3 * time.Minute)
	loggers.Info.Printf("Snapshot StartVM() Check VM running status")
	vmstatus := snap.instance.Get(ctx, vmid)
	loggers.Info.Printf("Snapshot StartVM() VM status after shutdown = %s", vmstatus)
	if vmstatus == "RUNNING" {
		loggers.Info.Printf("Snapshot StartVM() VM with VMid=%s has started in sleep timer", vmid)
		return "RUNNING"
	}
	// ----
	loggers.Error.Printf("Snapshot StartVM() VM with VMid=%s hasn't started in sleep timer", vmid)
	telegrambot.Tgbot.SendMessage(fmt.Sprintf("VM [%s] do not runnning. Cancel snapshot creation", vmid))
	return "ERROR:VM_DO_NOT_RUNNING"
}

// CleanUpOldSnapshots - function for listing of partucular snapshot
func (snap Snapshot) CleanUpOldSnapshots(ctx context.Context) {
	loggers.Info.Printf("Snapshot CleanUpOldSnapshots() starts")
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()
	// ---------
	//snaplist := snap.ListSnapshots(ctx)
	// ---------
	for _, vm := range snap.vms {
		go func(vmi config.VirtualMachine) {
			loggers.Info.Printf("Snapshot CleanUpOldSnapshots() CleanUp Snapshots for VM=%s", vmi.VMid)
		}(vm)
	}
}

// PrintStatusRegister - to get status information about snapshots
func (snap Snapshot) PrintStatusRegister() {
	data := [][]string{
		[]string{"A", "The Good", "500"},
		[]string{"B", "The Very very Bad Man", "288"},
		[]string{"C", "The Ugly", "120"},
		[]string{"D", "The Gopher", "800"},
	}
	// buffer to write table
	buf := new(bytes.Buffer)
	table := tablewriter.NewWriter(buf)
	table.SetHeader([]string{"VMname", "DiskID", "SnapshotID", "SnapshotDate", "Status"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
	loggers.Info.Println(buf.String())
}
