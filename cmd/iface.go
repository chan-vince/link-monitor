package cmd

import (
	"encoding/json"
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
type iface struct {
	name            string
	getStatFunc     getStat
	rxBytes         uint64
	txBytes         uint64
	client          msgClient
	pubIntervalSecs uint
}

type message struct {
	Name   string
	Value  uint64
}

type msgClient interface {
	Publish(topic string, message string) bool
}

type getStat func(stat string) uint64

func NewIface(name string, pubIntervalSecs uint) *iface {

	sv := iface{name: name}

	// On Linux we can read straight out of /sys/class/net but on Mac we
	// have to use netstat exec hackery. Windows...no thanks
	if runtime.GOOS == "linux"{
		sv.getStatFunc = sv.readFromFile
	} else if runtime.GOOS == "darwin" {
		sv.getStatFunc = sv.readFromNetstat
	} else {
		log.Fatalf("Unsupported OS type: %s\n", runtime.GOOS)
	}

	// Only I/O bytes are supported for now
	sv.rxBytes = 0
	sv.txBytes = 0

	// The client only needs a Publish method, see msgClient interface
	sv.client = nil

	// Minimum interval of 5 secs
	if pubIntervalSecs < 5 {
		sv.pubIntervalSecs = 5
	} else {
		sv.pubIntervalSecs = pubIntervalSecs
	}

	return &sv
}

func (sv *iface) RegisterMsgClient(client msgClient) {
	sv.client = client
}

func (sv *iface) Start() {

	if sv.client == nil{
		panic("Message client not set - call .RegisterMsgClient() method first")
	}

	// Do the first reading
	reading :=sv.read()
	sv.process(reading)

	// Start all the go routines
	go sv.ReadForever()
	go sv.PublishForever()
}

func (sv *iface) ReadForever() {
	for {
		reading :=sv.read()
		sv.process(reading)
		time.Sleep(time.Second)
	}
}

func (sv *iface) PublishForever() {
	for {
		sv.publish()
		time.Sleep(time.Duration(sv.pubIntervalSecs) * time.Second)
	}
}

func (sv *iface) publish() {
	messageMap := &message{
		Name:   sv.name,
		Value: sv.rxBytes,
	}
	messageJson, _ := json.Marshal(messageMap)
	sv.client.Publish("routingKey", string(messageJson))
}

func (sv *iface) read() uint64 {
	return sv.getStatFunc("rx_bytes")
}

func (sv *iface) process(newReading uint64) {
	// A restart, interface reload, counter zeroed or just wrapped around
	if newReading < sv.rxBytes {
		// Add the whole reading
		sv.rxBytes += newReading
	} else {
		sv.rxBytes += newReading - sv.rxBytes
	}

	fmt.Printf("newReading: %d\n", newReading)
	fmt.Printf("totalBytes: %d\n", sv.rxBytes)
}

func (sv *iface) readFromFile(stat string) uint64 {
	filePath := fmt.Sprintf("/sys/class/net/%s/statistics/%s", sv.name, stat)
	fmt.Println(filePath)
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
	value := strings.TrimSuffix(string(data), "\n")
	final, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(final)

	return processStringToUint64(data)
}

func (sv *iface) readFromNetstat(stat string) uint64 {

	var nthAwk string

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
