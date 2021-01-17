package rabbitmq

import (
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"sync/atomic"
	"time"
)

var clients []*client

type client struct {
	conn  *connection
	channel *channel
	channels []*channel
}

type connection struct {
	*amqp.Connection
}

type channel struct {
	*amqp.Channel
	closed int32
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func dialConfig(url string, config amqp.Config) (*connection, error) {
	conn, err := amqp.DialConfig(url, config)
	if err != nil {
		return nil, err
	}

	connection := &connection{
		conn,
	}

	go func() {
		for {
			reason, ok := <-connection.NotifyClose(make(chan *amqp.Error))
			// exit this goroutine if closed by system interrupt/ctrl-c
			if !ok {
				log.Debug("AMQP connection recovery cancelled")
				log.Info("Closing AMQP connection")
				break
			}
			log.Debugf("Connection closed: %v", reason)
			log.Debug("Recovering connection...")
			// reconnect if not closed by developer
			for {
				// wait 1s for reconnect
				time.Sleep(time.Second)

				conn, err := amqp.Dial(url)
				if err == nil {
					connection.Connection = conn
					log.Debug("Connection recovery success")
					break
				}
				//log.Printf("Reconnect failed: %v", err)
			}
		}
	}()

	return connection, nil
}

func (c *connection) channel() (*channel, error) {
	ch, err := c.Connection.Channel()
	if err != nil {
		return nil, err
	}

	channel := &channel{
		Channel: ch,
	}

	go func() {
		for {
			reason, ok := <-channel.Channel.NotifyClose(make(chan *amqp.Error))
			// exit this goroutine if closed by system interrupt/ctrl-c
			if !ok || channel.isClosed() {
				log.Debug("AMQP channel recovery cancelled")
				_ = channel.close() // close again, ensure closed flag set when connection closed
				break
			}
			log.Debugf("Channel closed: %v", reason)
			log.Debug("Recovering channel...")
			// reconnect if not closed by developer
			for {
				// wait 1s for connection reconnect
				time.Sleep(time.Second)

				ch, err := c.Connection.Channel()
				if err == nil {
					log.Debug("Channel recovery success")
					channel.Channel = ch
					break
				}
				//log.Printf("Channel recovery failed: %v", err)
			}
		}

	}()

	return channel, nil
}

func (ch *channel) close() error {
	if ch.isClosed() {
		return amqp.ErrClosed
	}

	atomic.StoreInt32(&ch.closed, 1)

	return ch.Channel.Close()
}

func (ch *channel) isClosed() bool {
	return atomic.LoadInt32(&ch.closed) == 1
}
