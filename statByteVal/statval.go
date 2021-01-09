package statByteVal

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// Statistic Value 'object' (something from /sys/class/net/<IFACE>/statistics/something
type statByteVal struct {
	name            string
	filePath        string
	totalBytes      uint64
	pubIntervalSecs int
}

type message struct {
	Name   string
	Value  uint64
}

type msgClient interface {
	Publish(topic string, message string) bool
}

func New(name string, filePath string) *statByteVal {

	result, errStr := isFile(filePath)
	if result == false {
		log.Fatalf("Invalid filePath for %s\n", name)
		panic(errStr)
	}

	sv := statByteVal{name: name, filePath: filePath}
	sv.totalBytes = 0
	sv.pubIntervalSecs = 1

	return &sv
}

func (sv *statByteVal) ReadForever(client msgClient) uint64 {
	for {
		var newReading uint64
		newReading = ReadFromFile(sv.filePath)
		sv.processNewReading(newReading)
		time.Sleep(time.Duration(sv.pubIntervalSecs) * time.Second)
		messageMap := &message{
			Name:   sv.name,
			Value: sv.totalBytes,
		}
		messageJson, _ := json.Marshal(messageMap)
		client.Publish("routingKey", string(messageJson))
	}

}

func (sv *statByteVal) processNewReading(newReading uint64) {
	fmt.Printf("newReading: %d\n", newReading)
	fmt.Printf("totalBytes: %d\n", sv.totalBytes)

	// A restart, interface reload, counter zeroed or just wrapped around
	if newReading < sv.totalBytes {
		// Add the whole reading
		sv.totalBytes += newReading
	} else {
		sv.totalBytes += newReading - sv.totalBytes
	}
}