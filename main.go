package main

import (
	"chanv/link-monitor/cmd"
	"chanv/link-monitor/msg"
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
	connDetails := msg.NewConnectionDetails(
		"192.168.10.3", 5672,
		"test", "test")
	msgClient := msg.Connect(connDetails)
	msg.Configure(msgClient.Channel)

	// Maybe a wait for connected?

	for i, link := range conf.Links {
		fmt.Println(i, link)
		iface := cmd.NewIface(link, conf.Broker.PublishInterval)
		iface.RegisterMsgClient(msgClient)
		go iface.Start()
	}

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	//time for cleanup before exit
	msg.CloseAll()
	fmt.Println("Adios!")
}
