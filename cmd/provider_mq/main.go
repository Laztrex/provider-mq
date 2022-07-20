package main

import (
	"github.com/rs/zerolog/log"
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

		rmqProducer.Queue = conf.ReplyTo

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
