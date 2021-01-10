package cmd

import (
	"fmt"
	"os"
	"runtime"
)

func CheckSys() string {
	if runtime.GOOS == "darwin" {
		fmt.Println("Hello from Mac")
	} else if runtime.GOOS == "linux" {
		fmt.Println("Hello from Linux")
	} else if runtime.GOOS == "windows" {
		fmt.Println("Windows not supported, exiting")
		os.Exit(1)
	} else {
		fmt.Println("Unknown OS, exiting")
		os.Exit(1)
	}
	return runtime.GOOS
}