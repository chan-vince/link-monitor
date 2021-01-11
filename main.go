package main

import (
	"chanv/link-monitor/cmd"
	"chanv/link-monitor/rabbitmq"
	"chanv/link-monitor/syslog"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// global conf stores the config file from viper
var conf *cmd.Config

func init() {
	conf = cmd.GetConf()

	if err := syslog.InitLogger(conf.Logging.Level); err != nil {
		fmt.Printf("Error: %+v\n", err)
		os.Exit(-1)
	}

	// Show config
	s, _ := json.MarshalIndent(conf, "", "\t")
	fmt.Printf("Using config:\n%s\n", string(s))
}

func main() {
	log.Printf("Log level: %s\n", syslog.ZapcoreLevel)
	logger := syslog.GetLogger()
	defer logger.Sync()
	log.Println("Logging to syslog.")

	logger.Error("Using config:")

	connDetails := rabbitmq.ConnectionDetails(
		conf.Broker.Host, conf.Broker.Port,
		conf.Broker.Username, conf.Broker.Password)

	msgClient := rabbitmq.Client(connDetails)

	msgClient.Configure(conf.Broker.ExchangeName, conf.Broker.ExchangeType)

	for _, link := range conf.Links {
		iface := cmd.NewIface(link)
		go iface.Start()
	}

	// Allow time for the first reading to happen before publishing
	time.Sleep(time.Second)
	go rabbitmq.StartPublishing(msgClient, cmd.ConstructRoutingKey(conf.Broker.RoutingKey, conf.KitId), conf.Broker.PublishInterval, cmd.GetIfaces())

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel

	//time for cleanup before exit
	rabbitmq.CloseAll()

	time.Sleep(time.Millisecond * 250)
	fmt.Println("Adios!")
}
