package statByteVal

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

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

func (sv *netIface) readFromFile(iface string, stat string) uint64 {
	filePath := fmt.Sprintf("/sys/class/net/%s/statistics/%s", iface, stat)
	fmt.Println(filePath)
	result, errStr := isFile(filePath)
	if result == false {
		log.Printf("Invalid filePath for %s\n", iface)
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

func (sv *netIface) readFromNetstat(iface string, stat string) uint64 {

	var nthAwk string

	if stat == "rx_bytes" {
		nthAwk = "$7"
	} else if stat == "tx_bytes" {
		nthAwk = "$10"
	} else{
		panic("unsupported stat")
	}

	cmd := fmt.Sprintf("netstat -I %s -nbf inet | tail -n 1 | awk '{print %s}'", iface, nthAwk)
	out, err := exec.Command("bash","-c",cmd).Output()
	if err != nil {
		fmt.Printf("Failed to execute command: %s", cmd)
	}
	fmt.Println(processStringToUint64(out))

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