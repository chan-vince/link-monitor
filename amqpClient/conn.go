package amqpClient

import (
	"fmt"
	"log"
	"github.com/streadway/amqp"
)

var connections []*amqp.Connection

type ConnectionDetails struct {
	hostname string
	port     int
	username string
	password string
}

type MsgClient struct {
	// Only support one channel per client for now
	Channel *amqp.Channel
}

func NewConnectionDetails(hostname string, port int, username string, password string) *ConnectionDetails {
	connDetails := ConnectionDetails{
		hostname: hostname,
		port:     port,
		username: username,
		password: password,
	}
	// Todo some more checking

	return &connDetails
}

func NewMsgClient(channel *amqp.Channel) *MsgClient {
	msgClient := MsgClient{Channel: channel}
	return &msgClient
}

func CloseAll() int {
	for _, conn := range connections {
		err := conn.Close()
		failOnError(err, "Failed to close connection")
	}
	return 0
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func Connect(connDetails *ConnectionDetails) *MsgClient {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		connDetails.username, connDetails.password,
		connDetails.hostname, connDetails.port)

	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
	connections = append(connections, conn)
	return conn
}

