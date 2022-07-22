package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"

	//"provider_mq/internal/app"
	"provider_mq/internal/controllers"
	"provider_mq/internal/utils"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	connectionString := utils.GetEnvVar("RMQ_URL")
	configs := utils.GetQueueConf()

	rmqConsumer := controllers.RMQSpec{
		ConnectionString: connectionString,
		Err:              make(chan error),
	}

	rmqProducer := controllers.RMQSpec{
		ConnectionString: connectionString,
		Err:              make(chan error),
	}

	err := rmqConsumer.Connect()
	if err != nil {
		rmqConsumer.OnError(err, "Failed to connect to RabbitMQ")
		panic(err)
	}

	err = rmqProducer.Connect()
	if err != nil {
		rmqProducer.OnError(err, "Failed to connect to RabbitMQ")
		panic(err)
	}

	wg.Add(len(configs) + 1)
	for _, conf := range configs {

		rmqConsumer.Queue = conf.QueueName
		rmqConsumer.Exchange = conf.Topic
		rmqConsumer.BindingKey = conf.BindingKey
		rmqConsumer.Args.Queue = conf.ArgsQueue

		if conf.ReplyTo != "" {
			rmqProducer.Queue = conf.ReplyTo
		}

		if conf.DLE == true {
			err = rmqConsumer.Channel.ExchangeDeclare(
				"dead_letter_exchange", // name
				"fanout",               // type
				true,                   // durable
				false,                  // auto-deleted
				false,                  // internal
				false,                  // no-wait
				nil,
			)
			if err != nil {
				fmt.Printf("error in ExchangeDeclare: %s", err)
			}

			_, err = rmqConsumer.Channel.QueueDeclare(
				"dead_letter_queue", // name
				false,               // durable
				false,               // delete when unused
				false,               // exclusive
				false,               // no-wait
				amqp.Table{
					"x-message-ttl":          60000,
					"x-dead-letter-exchange": "ML.MQ",
				},
			)
			if err != nil {
				fmt.Printf("error in ConsumeDeclare: %s", err)
			}

			err = rmqConsumer.Channel.QueueBind(
				"dead_letter_queue",    // name of the queue
				"",                     // binding key
				"dead_letter_exchange", // source exchange
				false,                  // noWait
				nil,                    // arguments
			)
			if err != nil {
				fmt.Printf("error in BindQueue: %s", err)
			}
		}

		rmqConsumer.ConsumeDeclare()
		rmqProducer.ProduceDeclare()

		go func(consumer controllers.RMQSpec) {
			log.Info().Msgf("Run consuming queue '%v'", consumer.Queue)
			consumer.ConsumeMessages()
		}(rmqConsumer)

		go func(producer controllers.RMQSpec) {
			log.Info().Msgf("Run consuming queue '%v'", producer.Queue)
			producer.ProduceMessage()
		}(rmqProducer)

	}

	go controllers.Predict()

	wg.Wait()

}
