package main

import (
	"chanv/link-monitor/cmd"
	"chanv/link-monitor/mq"
	"chanv/link-monitor/statByteVal"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// global conf stores the config file from viper
var conf *cmd.Config

func init() {
	conf = cmd.GetConf()

	// Show config
	s, _ := json.MarshalIndent(conf, "", "\t")
	fmt.Printf("Using config:\n%s", string(s))
}

func main() {
	//arguments := os.Args
	//if len(arguments) == 1 {
	//	fmt.Println("Please provide an argument!")
	//	os.Exit(1)
	//}
	iface := conf.Links[0]
	rx_bytes := statByteVal.New(iface)

	connDetails := mq.NewConnectionDetails("192.168.10.3", 5672, "test", "test")
	msgClient := mq.Connect(connDetails)
	mq.Configure(msgClient.Channel)

	rx_bytes.RegisterMsgClient(msgClient)
	go rx_bytes.Start()

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	//time for cleanup before exit
	mq.CloseAll()
	fmt.Println("Adios!")
}
