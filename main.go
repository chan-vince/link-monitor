package main

import (
	"chanv/link-monitor/amqpClient"
	"chanv/link-monitor/statByteVal"
	"chanv/link-monitor/cmd"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide an argument!")
		os.Exit(1)
	}

	cmd.CheckSys()

	filePath := arguments[1]
	rx_bytes := statByteVal.New("rx_bytes", filePath)

	connDetails := amqpClient.NewConnectionDetails("192.168.10.3", 5672, "test", "test")
	msgClient := amqpClient.Connect(connDetails)
	amqpClient.Configure(msgClient.Channel)

	go rx_bytes.ReadForever(msgClient)

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	//time for cleanup before exit
	amqpClient.CloseAll()
	fmt.Println("Adios!")
}
