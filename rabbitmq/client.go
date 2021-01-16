package rabbitmq

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
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

	var conn *connection
	var err error

	log.Info("Connecting to broker...")
	// If the broker isn't available when the application starts...we try with infinite persistence
	for {
		conn, err = dialConfig(url, config)
		if err != nil{
			time.Sleep(5 * time.Second)
		} else {
			log.Info("Connected!")
			break
		}
	}

	chann, err := conn.channel()
	failOnError(err, "Failed to get channel")

	client := client{
		conn:     conn,
		channel:  chann,
		channels: []*channel{chann},
	}

	clients = append(clients, &client)

	return &client
}

func (client *client) Configure(exchangeName, exchangeType string) {
	// The exchange is hardcoded to the amq.topic exchange

	// Declare exchange
	err := client.channel.ExchangeDeclare(exchangeName, exchangeType, true, false, false, false, nil)
	failOnError(err, "Failed to declare exchange")

}

func (client *client) Publish(routingKey string, message string) bool {
	exchangeName := "amq.topic"

	publishing := amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(message),
	}
	if ! client.channel.isClosed() {
		err := client.channel.Publish(exchangeName, routingKey, false, false, publishing)
		if err != nil {
			//log.Print("Failed to publish a message\n", err)
		} else {
			log.Debug("Published:")
			log.Debugf("\tkey: %s\n", routingKey)
			log.Debugf("\tmsg: %s\n", message)
		}
	}

	return true
}

func CloseAll() int {
	for _, client := range clients {
		err := client.channel.close()
		failOnError(err, "Failed to close channel")
		err = client.conn.Close()
		failOnError(err, "Failed to close connection")
	}
	return 0
}

