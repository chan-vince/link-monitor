package main

import (
	"chanv/link-monitor/cmd"
	"chanv/link-monitor/rabbitmq"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// global conf stores the config file from viper
var conf *cmd.Config

func init() {
	conf = cmd.GetConf()

	log.SetOutput(os.Stdout)
	level, err := log.ParseLevel(conf.Logging.Level)
	if err != nil {
		level = log.ErrorLevel
	}
	log.SetLevel(level)

	// Show config
	s, _ := json.MarshalIndent(conf, "", "\t")
	log.Debugf("Using config:\n%s\n", string(s))
}

func main() {
	log.Printf("Log level: %s\n", log.GetLevel())

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
	log.Println("Adios!")
}
