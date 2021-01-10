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
	fmt.Printf("Using config:\n%s\n", string(s))
}

func main() {
	connDetails := msg.NewConnectionDetails(
		conf.Broker.Host, conf.Broker.Port,
		conf.Broker.Username, conf.Broker.Password)

	msgClient := msg.NewConnection(connDetails)

	msg.ConfigureBroker(
		msgClient,
		conf.Broker.ExchangeName, conf.Broker.ExchangeType,
		cmd.ConstructRoutingKey(conf.Broker.RoutingKey, conf.KitId),
		)

	// Maybe a wait for connected?

	var links []*cmd.Iface

	for _, link := range conf.Links {
		iface := cmd.NewIface(link)
		links = append(links, iface)
		go iface.Start()
	}

	go msg.StartPublishing(msgClient, cmd.ConstructRoutingKey(conf.Broker.RoutingKey, conf.KitId), conf.Broker.PublishInterval, links)

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	//time for cleanup before exit
	msg.CloseAll()
	fmt.Println("Adios!")
}
