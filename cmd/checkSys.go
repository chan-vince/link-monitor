package cmd

import (
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
)

func CheckSys() string {
	if runtime.GOOS == "darwin" {
		log.Debug("OS: Mac")
	} else if runtime.GOOS == "linux" {
		log.Debug("OS: Linux")
	} else if runtime.GOOS == "windows" {
		log.Error("Windows is not supported, exiting")
		os.Exit(1)
	} else {
		log.Error("Unknown OS, exiting")
		os.Exit(1)
	}
	return runtime.GOOS
}