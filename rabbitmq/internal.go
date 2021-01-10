package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
	"sync/atomic"
	"time"
)

var clients []*client

type client struct {
	conn  *Connection
	channel *Channel
	channels []*Channel
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func dialConfig(url string, config amqp.Config) (*Connection, error) {
	conn, err := amqp.DialConfig(url, config)
	if err != nil {
		return nil, err
	}

	connection := &Connection{
		conn,
	}

	go func() {
		for {
			reason, ok := <-connection.NotifyClose(make(chan *amqp.Error))
			// exit this goroutine if closed by developer
			if !ok {
				log.Print("Auto reconnect cancelled")
				break
			}
			log.Print("connection closed, reason: %v", reason)

			// reconnect if not closed by developer
			for {
				// wait 1s for reconnect
				time.Sleep(time.Second)

				conn, err := amqp.Dial(url)
				if err == nil {
					connection.Connection = conn
					log.Print("reconnect success")
					break
				}

				log.Print("reconnect failed, err: %v", err)
			}
		}
	}()

	return connection, nil
}

// isClosed indicate closed by developer
func (ch *Channel) isClosed() bool {
	return atomic.LoadInt32(&ch.closed) == 1
}

func (c *Connection) channel() (*Channel, error) {
	ch, err := c.Connection.Channel()
	if err != nil {
		return nil, err
	}

	channel := &Channel{
		Channel: ch,
	}

	go func() {
		for {
			reason, ok := <-channel.Channel.NotifyClose(make(chan *amqp.Error))
			// exit this goroutine if closed by developer
			if !ok || channel.isClosed() {
				log.Print("auto channel cancelled")
				_ = channel.Close() // close again, ensure closed flag set when connection closed
				break
			}
			log.Print("channel closed, reason: %v", reason)

			// reconnect if not closed by developer
			for {
				// wait 1s for connection reconnect
				time.Sleep(time.Second)

				ch, err := c.Connection.Channel()
				if err == nil {
					log.Print("channel recreate success")
					channel.Channel = ch
					break
				}

				log.Print("channel recreate failed, err: %v", err)
			}
		}

	}()

	return channel, nil
}