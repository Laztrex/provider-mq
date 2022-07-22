package controllers

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"provider_mq/internal/utils"
)

type RMQSpec struct {
	Queue            string
	ConnectionString string
	Connection       *amqp.Connection
	Channel          *amqp.Channel
	Exchange         string
	RoutingKey       string
	BindingKey       string
	ReplyTo          string
	Err              chan error
	Args             RMQArgs
}

type RMQArgs struct {
	Topic amqp.Table
	Queue amqp.Table
}

func (conn RMQSpec) OnError(err error, msg string) {
	if err != nil {
		log.Err(err).Msgf("Error message on '%s' queue. Error message: %s", conn.Queue, msg)
	}
}

func (conn *RMQSpec) Connect() error {
	var err error

	tlsConf := utils.GetTlsConf()

	conn.Connection, err = amqp.DialTLS(conn.ConnectionString, tlsConf)
	if err != nil {
		return fmt.Errorf("error in creating rabbitmq connection with %s : %s", conn.ConnectionString, err.Error())
	}

	go func() {
		<-conn.Connection.NotifyClose(make(chan *amqp.Error)) //Listen to NotifyClose
		conn.Err <- errors.New("connection Closed")
	}()

	conn.Channel, err = conn.Connection.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel")
	}

	return nil
}

func (conn *RMQSpec) ExchangeDeclare() error {
	err := conn.Channel.ExchangeDeclare(
		conn.Exchange, // name
		"topic",       // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return fmt.Errorf("error in Exchange ConsumeDeclare: %s", err)
	}

	return nil
}

func (conn *RMQSpec) QueueDeclare() error {
	_, err := conn.Channel.QueueDeclare(
		conn.Queue,      // name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		conn.Args.Queue, // arguments
	)
	if err != nil {
		return fmt.Errorf("error in declaring the queue %s", err)
	}

	//err = conn.Channel.Qos(1, 0, false)
	//if err != nil {
	//	return fmt.Errorf("error set qos %s", err)
	//}

	return nil
}

func (conn *RMQSpec) QueueBind() error {
	err := conn.Channel.QueueBind(
		conn.Queue,      // name of the queue
		conn.BindingKey, // binding key
		conn.Exchange,   // source exchange
		false,           // noWait
		nil,             // arguments
	)
	if err != nil {
		return fmt.Errorf("queue Bind error: %s", err)
	}

	return nil
}

//Reconnect reconnects the connection
func (conn *RMQSpec) Reconnect() error {

	if err := conn.Connect(); err != nil {
		return err
	}
	log.Printf("INFO: channel reconnection: %s", conn.ConnectionString)
	//if err := connProducer.BindQueue(); err != nil {
	//	return err
	//}

	return nil
}
