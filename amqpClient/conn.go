package amqpClient

import (
	"fmt"
	"log"
	"github.com/streadway/amqp"
)

type ConnectionDetails struct {
	hostname string
	port     int
	username string
	password string
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

var connections []*amqp.Connection

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

func Connect(connDetails *ConnectionDetails) *amqp.Connection {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		connDetails.username, connDetails.password,
		connDetails.hostname, connDetails.port)

	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
	connections = append(connections, conn)
	return conn
}

