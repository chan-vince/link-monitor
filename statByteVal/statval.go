package statByteVal

import (
	"fmt"
	"log"
	"os"
	"time"
	"encoding/json"
)

// Statistic Value 'object' (something from /sys/class/net/<IFACE>/statistics/something
type statByteVal struct {
	name            string
	filePath        string
	currentByteVal  uint64
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
	sv.currentByteVal = 0
	sv.pubIntervalSecs = 2

	return &sv
}

func (sv *statByteVal) ReadForever(client msgClient) uint64 {
	for {
		var value uint64
		value = ReadFromFile(sv.filePath)
		HandleReadValue(value)
		time.Sleep(time.Duration(sv.pubIntervalSecs) * time.Second)
		messageMap := &message{
			Name:   sv.name,
			Value: value,
		}
		messageJson, _ := json.Marshal(messageMap)
		client.Publish("routingKey", string(messageJson))
	}

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