package msg

import (
	"fmt"
	"log"
	"github.com/streadway/amqp"
)

var connections []*connection

type connectionDetails struct {
	hostname string
	port     int
	username string
	password string
}

func NewConnectionDetails(hostname string, port int, username string, password string) *connectionDetails {
	connDetails := connectionDetails{
		hostname: hostname,
		port:     port,
		username: username,
		password: password,
	}
	// Todo some more checking

	return &connDetails
}

type connection struct {
	conn *amqp.Connection
	channel *amqp.Channel
	channels []*amqp.Channel
}

func NewConnection(connDetails *connectionDetails) *connection {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		connDetails.username, connDetails.password,
		connDetails.hostname, connDetails.port)

	properties := amqp.Table{"user-id": "link-monitor"}

	config := amqp.Config{
		nil, "/", 1, 0, 10, nil,
		properties, "en_US", nil}

	conn, err := amqp.DialConfig(url, config)
	failOnError(err, "Failed to connect to RabbitMQ")

	channel, err := conn.Channel()
	failOnError(err, "Failed to get channel")

	connection := connection{
		conn: conn,
		channel: channel,
		channels: []*amqp.Channel{channel},
	}

	connections = append(connections, &connection)

	return &connection
}

func CloseAll() int {
	for _, conn := range connections {
		err := conn.channel.Close()
		failOnError(err, "Failed to close channel")
		err = conn.conn.Close()
		failOnError(err, "Failed to close connection")
	}
	return 0
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func ConfigureBroker(connection *connection, exchangeName, exchangeType, routingKey string) {
	// The exchange is hardcoded to the amq.topic exchange

	queueName := "testq"

	// Declare exchange
	err := connection.channel.ExchangeDeclare(exchangeName, exchangeType, true, false, false, false, nil)
	failOnError(err, "Failed to declare exchange")

	// Declare queueName
	_, err = connection.channel.QueueDeclare(queueName, true, true, false, false, nil)
	failOnError(err, "Failed to declare queueName")

	// Bind queue to exchange
	err = connection.channel.QueueBind(queueName, routingKey, exchangeName, false, nil)
	failOnError(err, "Failed to bind queue to exchange")
}

func (connection *connection) Publish(routingKey string, message string) bool {
	fmt.Printf("publish key: %s\n", routingKey)
	fmt.Printf("publish msg: %s\n", message)
	exchangeName := "amq.topic"

	publishing := amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(message),
	}

	err := connection.channel.Publish(exchangeName, routingKey, false, false, publishing)
	failOnError(err, "Failed to publish a message")

	return true
}