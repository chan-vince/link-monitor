package statByteVal

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Statistic Value 'object' (something from /sys/class/net/<IFACE>/statistics/something
type statByteVal struct {
	name            string
	filePath        string
	currentByteVal  uint64
	pubIntervalSecs int
}

func New(name string, filePath string) *statByteVal {

	result, errStr := isFile(filePath)
	if result == false {
		log.Fatalf("Invalid filePath for %s\n", name)
		panic(errStr)
	}

	sv := statByteVal{name: name, filePath: filePath}
	sv.currentByteVal = 0
	sv.pubIntervalSecs = 0

	return &sv
}

func (sv *statByteVal) ReadForever() uint64 {
	for {
		var value uint64
		value = ReadFromFile(sv.filePath)
		HandleReadValue(value)
		time.Sleep(time.Second)
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