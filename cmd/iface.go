package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Statistic Value 'object' (something from /sys/class/net/<IFACE>/statistics/something
type Iface struct {
	name            string
	getStatFunc     getStat
	// Only I/O bytes are supported for now
	rxBytes         stat
	txBytes         stat
}

var ifaces []*Iface

type stat struct {
	name  string
	value uint64
}

type StatMessage struct {
	RxBytes uint64 `json:"rx_bytes"`
	TxBytes uint64 `json:"tx_bytes"`
}

type getStat func(stat string) uint64

func NewIface(name string) *Iface {

	sv := Iface{
		name: name,
		rxBytes: stat{
			name: "rx_bytes",
			value: 0,
		},
		txBytes: stat{
			name: "tx_bytes",
			value: 0,
		},
	}

	// On Linux we can read straight out of /sys/class/net but on Mac we
	// have to use netstat exec hackery. Windows...no thanks
	if runtime.GOOS == "linux"{
		sv.getStatFunc = sv.readFromFile
	} else if runtime.GOOS == "darwin" {
		sv.getStatFunc = sv.readFromNetstat
	} else {
		log.Fatalf("Unsupported OS type: %s\n", runtime.GOOS)
	}

	//// Minimum interval of 5 secs
	//if pubIntervalSecs < 5 {
	//	sv.pubIntervalSecs = 5
	//} else {
	//	sv.pubIntervalSecs = pubIntervalSecs
	//}

	ifaces = append(ifaces, &sv)

	return &sv
}

func GetIfaces() []*Iface {
	return ifaces
}

func (sv *Iface) Name() string {
	return sv.name
}

func (sv *Iface) Start() {

	// Populate the first readings
	sv.readAll()

	// Start the go routine to read
	go sv.ReadForever()
}

func (sv *Iface) ReadForever() {
	for {
		sv.readAll()
		time.Sleep(time.Second)
	}
}

//func (sv *Iface) PublishForever() {
//	for {
//		sv.publish()
//		time.Sleep(time.Duration(sv.pubIntervalSecs) * time.Second)
//	}
//}

func (sv *Iface) GetIfaceStatsMessage() StatMessage {
	messageMap := &StatMessage{
		RxBytes: sv.rxBytes.value,
		TxBytes: sv.txBytes.value,
	}
	return *messageMap
}

func (sv *Iface) readAll() {
	sv.rxBytes.update(sv.getStatFunc(sv.rxBytes.name))
	sv.txBytes.update(sv.getStatFunc(sv.txBytes.name))
}

func (st *stat) update(newValue uint64) {
	// A restart, interface reload, counter zeroed or just wrapped around
	if newValue < st.value {
		// Add the whole reading
		st.value += newValue
	} else {
		st.value += newValue - st.value
	}
}

func (sv *Iface) readFromFile(stat string) uint64 {
	filePath := fmt.Sprintf("/sys/class/net/%s/statistics/%s", sv.name, stat)

	result, errStr := isFile(filePath)
	if result == false {
		log.Printf("Invalid filePath for %s\n", sv.name)
		panic(errStr)
	}

	// Open and read
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	return processStringToUint64(data)
}

func (sv *Iface) readFromNetstat(stat string) uint64 {

	var nthAwk string

	// Disgusting
	if stat == "rx_bytes" {
		nthAwk = "$7"
	} else if stat == "tx_bytes" {
		nthAwk = "$10"
	} else{
		panic("unsupported stat")
	}

	cmd := fmt.Sprintf("netstat -I %s -nbf inet | tail -n 1 | awk '{print %s}'", sv.name, nthAwk)
	out, err := exec.Command("bash","-c",cmd).Output()
	if err != nil {
		fmt.Printf("Failed to execute command: %s", cmd)
	}

	return processStringToUint64(out)
}

func processStringToUint64(input []byte) uint64 {
	value := strings.TrimSuffix(string(input), "\n")
	final, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	return final
}

func isFile(filePath string) (bool, string) {
	info, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		errStr := fmt.Sprintf("Does not exist: %s\n\n", filePath)
		return false, errStr
	}

	if info.IsDir() {
		errStr := fmt.Sprintf("Not a valid file: %s\n\n", filePath)
		return false, errStr
	}

	return true, ""
}
