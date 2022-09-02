package rabbit

import (
	"crypto/tls"
	"os"

	"github.com/google/uuid"
	"github.com/streadway/amqp"

	"github.com/h14yhv/golang-lib/clock"
	"github.com/h14yhv/golang-lib/log"
)

type rabbitConnection struct {
	logger     log.Logger
	connection *amqp.Connection
	channel    *amqp.Channel
	config     Config
	tlsConfig  *tls.Config
	uuid       string
}

func NewService(conf Config, tlsConf *tls.Config) Service {
	logger, _ := log.New(Module, log.DebugLevel, true, os.Stdout)
	rb := &rabbitConnection{
		logger:    logger,
		config:    conf,
		tlsConfig: tlsConf,
	}
	if err := rb.newConnection(); err != nil {
		panic(err)
	}
	if err := rb.newChannel(); err != nil {
		panic(err)
	}
	// Monitor
	rb.monitor()
	// Success
	return rb
}

func (r *rabbitConnection) newConnection() error {
	var err error
	if r.config.Secure || r.tlsConfig != nil {
		if r.tlsConfig == nil {
			r.tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		r.connection, err = amqp.DialTLS(r.config.String(), r.tlsConfig)
	} else {
		r.connection, err = amqp.Dial(r.config.String())
	}
	if err != nil {
		return err
	}
	// Success
	return nil
}

func (r *rabbitConnection) newChannel() error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	r.uuid = id.String()
	r.channel, err = r.connection.Channel()
	if err != nil {
		return err
	}
	// Success
	return nil
}

func (r *rabbitConnection) monitor() {
	// Reconnect connection
	go func() {
		for {
			if r.connection != nil {
				reason, ok := <-r.connection.NotifyClose(make(chan *amqp.Error))
				if !ok {
					r.logger.Info("connection closed")
					break
				}
				if reason != nil {
					r.logger.Infof("connection closed, reason: %v", reason)
				}
			}
			// Reconnect
			for {
				clock.Sleep(ScheduleReconnect)
				if err := r.newConnection(); err == nil {
					r.logger.Info("recreate connection success!")
					break
				} else {
					r.logger.Errorf("recreate connection failed, reason: %v", err)
				}
			}
		}

	}()
	// Reconnect channel
	go func() {
		for {
			if r.channel != nil {
				reason, ok := <-r.channel.NotifyClose(make(chan *amqp.Error))
				if !ok {
					r.logger.Info("channel closed")
					break
				}
				if reason != nil {
					r.logger.Infof("channel closed, reason: %v", reason)
				}
			}
			// Reconnect
			for {
				clock.Sleep(ScheduleReconnect)
				if r.connection != nil && !r.connection.IsClosed() {
					if err := r.newChannel(); err == nil {
						r.logger.Info("recreate channel success!")
						break
					} else {
						r.logger.Errorf("recreate channel failed, reason: %v", err)
					}
				}
			}
		}
	}()
}

func (r *rabbitConnection) DeclareExchange(name, kind string, durable bool) error {
	// Success
	return r.channel.ExchangeDeclare(name, kind, durable, false, false, false, nil)
}

func (r *rabbitConnection) DeclareQueue(name string, durable bool, priority int, ttl clock.Duration) error {
	table := amqp.Table{}
	if priority > 0 {
		table["x-max-priority"] = priority
	}
	if ttl > 0 {
		table["message-ttl"] = ttl.Milliseconds()
	}
	if _, err := r.channel.QueueDeclare(name, durable, false, false, false, table); err != nil {
		return err
	}
	// Success
	return nil
}

func (r *rabbitConnection) BindQueue(queue, exchange string) error {
	// Success
	return r.channel.QueueBind(queue, queue, exchange, false, nil)
}

func (r *rabbitConnection) Qos(prefetchCount int) error {
	if err := r.channel.Qos(prefetchCount, 0, false); err != nil {
		return err
	}
	// Success
	return nil
}

func (r *rabbitConnection) Publish(exchange, queue string, message Message) error {
	if message.ContentType == "" {
		message.ContentType = MIMETextPlain
	}
	if message.Mode == 0 {
		message.Mode = Transient
	}
	err := r.channel.Publish(
		exchange,
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType:  message.ContentType,
			DeliveryMode: message.Mode,
			Priority:     message.Priority,
			Body:         message.Body,
		},
	)
	if err != nil {
		r.logger.Errorf("publish failed, reason: %v", err)
		clock.Sleep(SchedulePublish)
		return r.Publish(exchange, queue, message)
	}
	// Success
	return nil
}

func (r *rabbitConnection) Consume(queue string, auto bool, prefetchCount int, callback Consumer) error {
	if prefetchCount == 0 {
		prefetchCount = 1
	}
	if err := r.Qos(prefetchCount); err != nil {
		r.logger.Errorf("qos failed, reason: %v", err)
		clock.Sleep(ScheduleConsume)
		return r.Consume(queue, auto, prefetchCount, callback)
	}
	deliveries := make(chan amqp.Delivery)
	go func() {
		for {
			d, err := r.channel.Consume(queue, r.uuid, auto, false, false, false, nil)
			if err != nil {
				clock.Sleep(ScheduleConsume)
				continue
			}
			for msg := range d {
				deliveries <- msg
			}
		}
	}()
	// Process
	for msg := range deliveries {
		err := callback(msg.Body)
		if err == nil {
			if !auto {
				// Ack
				msg.Ack(false)
			}
		} else {
			if !auto {
				// Requeue
				msg.Nack(false, true)
			}
		}
	}
	// Success
	return nil
}
