package main

import (
	"github.com/rs/zerolog/log"
	"sync"

	"provider_mq/internal/controllers"
	"provider_mq/internal/service"
	"provider_mq/internal/utils"
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

	err := rmqConsumer.GetConnect()
	if err != nil {
		rmqConsumer.OnError(err, "Failed to connect to RabbitMQ")
		panic(err)
	}

	err = rmqProducer.GetConnect()
	if err != nil {
		rmqProducer.OnError(err, "Failed to connect to RabbitMQ")
		panic(err)
	}

	wg.Add(len(configs) + 1)
	for _, conf := range configs {

		rmqConsumer.GetChannel()
		rmqProducer.GetChannel()

		rmqConsumer.Queue = conf.QueueName
		rmqConsumer.Exchange = conf.Topic
		rmqConsumer.BindingKey = conf.BindingKey
		rmqConsumer.Args.Queue = conf.ArgsQueue

		if conf.ReplyTo != "" {
			rmqProducer.Queue = conf.ReplyTo
		}

		if conf.DLE == true {
			rmqConsumer.SetDLE()
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

	go service.Predict()

	wg.Wait()
}
