package main

import (
	"chanv/link-monitor/cmd"
	"chanv/link-monitor/rabbitmq"
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
	connDetails := rabbitmq.ConnectionDetails(
		conf.Broker.Host, conf.Broker.Port,
		conf.Broker.Username, conf.Broker.Password)

	msgClient := rabbitmq.Client(connDetails)

	rabbitmq.Configure(
		msgClient,
		conf.Broker.ExchangeName, conf.Broker.ExchangeType,
		cmd.ConstructRoutingKey(conf.Broker.RoutingKey, conf.KitId),
		)

	var links []*cmd.Iface

	for _, link := range conf.Links {
		iface := cmd.NewIface(link)
		links = append(links, iface)
		go iface.Start()
	}

	go rabbitmq.StartPublishing(msgClient, cmd.ConstructRoutingKey(conf.Broker.RoutingKey, conf.KitId), conf.Broker.PublishInterval, links)

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	//time for cleanup before exit
	rabbitmq.CloseAll()
	fmt.Println("Adios!")
}
