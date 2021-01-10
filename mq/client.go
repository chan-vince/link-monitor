package mq

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

	channel, err := conn.Channel()
	failOnError(err, "Failed to get channel")

	msgClient := NewMsgClient(channel)

	return msgClient
}

func Configure(channel *amqp.Channel) {
	// The exchange is hardcoded to the amq.topic exchange

	exchangeName := "amq.topic"
	exchangeType := "topic"
	queueName := "testq"
	routingKey := "routingKey"

	// Declare exchange
	err := channel.ExchangeDeclare(exchangeName, exchangeType, true, false, false, false, nil)
	failOnError(err, "Failed to declare exchange")

	// Declare queueName
	_, err = channel.QueueDeclare(queueName, true, true, false, false, nil)
	failOnError(err, "Failed to declare queueName")

	// Bind queue to exchange
	err = channel.QueueBind(queueName, routingKey, exchangeName, false, nil)
	failOnError(err, "Failed to bind queue to exchange")
}

func (msgClient *MsgClient) Publish(routingKey string, message string) bool {
	fmt.Printf("publish key: %s\n", routingKey)
	fmt.Printf("publish msg: %s\n", message)
	exchangeName := "amq.topic"

	publishing := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	}

	err := msgClient.Channel.Publish(exchangeName, routingKey, false, false, publishing)
	failOnError(err, "Failed to publish a message")

	return true
}