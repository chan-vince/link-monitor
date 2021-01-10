package statByteVal

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"
)

// Statistic Value 'object' (something from /sys/class/net/<IFACE>/statistics/something
type netIface struct {
	name            string
	getStatFunc     getStat
	rxBytes         uint64
	txBytes         uint64
	client          msgClient
	pubIntervalSecs int
}

type message struct {
	Name   string
	Value  uint64
}

type msgClient interface {
	Publish(topic string, message string) bool
}

type getStat func(iface string, stat string) uint64

func New(iface string) *netIface {

	sv := netIface{name: iface}

	if runtime.GOOS == "linux"{
		sv.getStatFunc = sv.readFromFile
	} else if runtime.GOOS == "darwin" {
		sv.getStatFunc = sv.readFromNetstat
	}

	sv.rxBytes = 0
	sv.txBytes = 0
	sv.client = nil
	sv.pubIntervalSecs = 1

	return &sv
}

func (sv *netIface) RegisterMsgClient(client msgClient) {
	sv.client = client
}

func (sv *netIface) Start() {

	if sv.client == nil{
		panic("Message client not set")
	}

	// Do the first reading
	reading :=sv.read()
	sv.process(reading)

	// Start all the go routines
	go sv.ReadForever()
	go sv.PublishForever()
}

func (sv *netIface) ReadForever() {
	for {
		reading :=sv.read()
		sv.process(reading)
		time.Sleep(time.Second)
	}
}

func (sv *netIface) PublishForever() {
	for {
		sv.publish()
		time.Sleep(time.Duration(sv.pubIntervalSecs) * time.Second)
	}
}

func (sv *netIface) publish() {
	messageMap := &message{
		Name:   sv.name,
		Value: sv.rxBytes,
	}
	messageJson, _ := json.Marshal(messageMap)
	sv.client.Publish("routingKey", string(messageJson))
}

func (sv *netIface) read() uint64 {
	return sv.getStatFunc(sv.name, "rx_bytes")
}

func (sv *netIface) process(newReading uint64) {
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
