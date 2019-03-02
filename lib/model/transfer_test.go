package model

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"testing"
	"time"

	"github.com/syncthing/syncthing/lib/config"
	"github.com/syncthing/syncthing/lib/db"
	"github.com/syncthing/syncthing/lib/fs"
	"github.com/syncthing/syncthing/lib/protocol"
)

func TestTransferFake(t *testing.T) {
	// l.SetDebug("model", true)
	mustDir := func() string {
		dir, err := ioutil.TempDir("", "")
		if err != nil {
			t.Fatal(err)
		}
		return dir
	}
	mustDb := func() *db.Lowlevel {
		db, err := db.Open(mustDir())
		if err != nil {
			t.Fatal(err)
		}
		return db
	}

	db1 := mustDb()
	db2 := mustDb()

	cfg1 := config.New(device1)
	cfg1.Devices = append(cfg1.Devices, config.NewDeviceConfiguration(device2, "device2"))

	workload := "fake/?files=1&seed=42&sizeavg=1073741824"
	workload = "fake/?files=10000&seed=42&sizeavg=1024"
	fcfg1 := config.NewFolderConfiguration(device1, "folder", "folder1", fs.FilesystemTypeFake, workload)
	fcfg1.Devices = append(fcfg1.Devices, config.FolderDeviceConfiguration{
		DeviceID: device2,
	})
	cfg1.Folders = append(cfg1.Folders, fcfg1)

	cfg2 := config.New(device2)
	cfg2.Devices = append(cfg2.Devices, config.NewDeviceConfiguration(device1, "device2"))
	fcfg2 := config.NewFolderConfiguration(device2, "folder", "folder2", fs.FilesystemTypeBasic, mustDir())
	fcfg2.Devices = append(fcfg2.Devices, config.FolderDeviceConfiguration{
		DeviceID: device1,
	})
	//fcfg2.Copiers = 10000
	cfg2.Folders = append(cfg2.Folders, fcfg2)

	cfg1w := config.Wrap("", cfg1)
	cfg2w := config.Wrap("", cfg2)

	model1 := NewModel(cfg1w, device1, "", "", db1, nil)
	go model1.Serve()
	for _, folder := range cfg1.Folders {
		model1.AddFolder(folder)
		model1.StartFolder(folder.ID)
	}

	model2 := NewModel(cfg2w, device2, "", "", db2, nil)
	go model2.Serve()
	for _, folder := range cfg2.Folders {
		model2.AddFolder(folder)
		model2.StartFolder(folder.ID)
	}

	device1pipe, device2pipe := net.Pipe()
	device1Conn := protocol.NewConnection(device1, device1pipe, device1pipe, model2, "", protocol.CompressNever)
	device1Pipe := &pipeConn{

		Connection: device1Conn,
		pipe:       device1pipe,
	}
	device2Conn := protocol.NewConnection(device2, device2pipe, device2pipe, model1, "", protocol.CompressNever)
	device2Pipe := &pipeConn{

		Connection: device2Conn,
		pipe:       device2pipe,
	}

	model1.AddConnection(device2Pipe, protocol.HelloResult{
		DeviceName:    "device2",
		ClientName:    "syncthing",
		ClientVersion: "v1.0.0",
	})

	model2.AddConnection(device1Pipe, protocol.HelloResult{
		DeviceName:    "device1",
		ClientName:    "syncthing",
		ClientVersion: "v1.0.0",
	})

	// Make sure stuff is sent out to the other device.
	model1.ScanFolders()

	interval := 5 * time.Second
	for {
		need := model2.NeedSize("folder")
		if need.Bytes == 0 && need.Files == 0 && need.Directories == 0 {
			break
		}
		log.Printf("%#v\n", need)
		time.Sleep(interval)
	}

	ifs := fcfg2.Filesystem().(*fs.InstrumentedFilesystem)
	durs := ifs.Durations()
	fmt.Println("Durations:", durs)
	counts := ifs.Counts()
	fmt.Println("Counts:", counts)
	fmt.Println("Per file:")
	tot := time.Duration(0)
	for k := range durs {
		tot += durs[k] / 10000
		fmt.Printf("%10s: %10s (%f)\n", k, durs[k]/10000, float64(counts[k])/10000)
	}
	fmt.Println("Total per file:", tot)
}

type pipeConn struct {
	protocol.Connection
	pipe net.Conn
}

func (p *pipeConn) Type() string {
	return "pipe"
}

func (p *pipeConn) Transport() string {
	return "pipe"
}

func (p *pipeConn) Crypto() string {
	return "none"
}

func (p *pipeConn) Priority() int {
	return 0
}

func (p *pipeConn) RemoteAddr() net.Addr {
	return p.pipe.RemoteAddr()
}

func (p *pipeConn) String() string {
	return p.pipe.RemoteAddr().String()
}
