package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func Client(connDetails *connectionDetails) *client {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		connDetails.username, connDetails.password,
		connDetails.hostname, connDetails.port)

	properties := amqp.Table{"user-id": "link-monitor"}

	config := amqp.Config{
		nil, "/", 1, 0, 10, nil,
		properties, "en_US", nil}

	var conn *Connection
	var err error

	log.Print("Connecting to broker...")
	// If the broker isn't available when the application starts...we try with infinite persistence
	for {
		conn, err = dialConfig(url, config)
		if err != nil{
			time.Sleep(5 * time.Second)
		} else {
			log.Print("Connected!")
			break
		}
	}

	channel, err := conn.channel()
	failOnError(err, "Failed to get channel")

	client := client{
		conn: conn,
		channel: channel,
		channels: []*Channel{channel},
	}

	clients = append(clients, &client)

	return &client
}

type Connection struct {
	*amqp.Connection
}

type Channel struct {
	*amqp.Channel
	closed int32
}

func CloseAll() int {
	for _, client := range clients {
		err := client.channel.Close()
		failOnError(err, "Failed to close channel")
		err = client.conn.Close()
		failOnError(err, "Failed to close connection")
	}
	return 0
}

func Configure(client *client, exchangeName, exchangeType, routingKey string) {
	// The exchange is hardcoded to the amq.topic exchange

	queueName := "testq"

	// Declare exchange
	err := client.channel.ExchangeDeclare(exchangeName, exchangeType, true, false, false, false, nil)
	failOnError(err, "Failed to declare exchange")

	// Declare queueName
	_, err = client.channel.QueueDeclare(queueName, true, true, false, false, nil)
	failOnError(err, "Failed to declare queueName")

	// Bind queue to exchange
	err = client.channel.QueueBind(queueName, routingKey, exchangeName, false, nil)
	failOnError(err, "Failed to bind queue to exchange")
}

func (client *client) Publish(routingKey string, message string) bool {
	fmt.Printf("publish key: %s\n", routingKey)
	fmt.Printf("publish rabbitmq: %s\n", message)
	exchangeName := "amq.topic"

	publishing := amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(message),
	}

	err := client.channel.Publish(exchangeName, routingKey, false, false, publishing)
	if err != nil {
		log.Print(err, "Failed to publish a message")
	}

	return true
}

