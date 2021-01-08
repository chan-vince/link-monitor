package cmd

import (
	"fmt"
	"runtime"
)

func CheckSys() bool {
	if runtime.GOOS == "darwin" {
		fmt.Println("Hello from Mac")
	}
	// todo make sure we're linux
}