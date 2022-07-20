package main

import (
	"github.com/rs/zerolog/log"
	//"provider_mq/internal/app"
	"provider-mq/internal/controllers"
	"provider-mq/internal/utils"
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

	err := rmqConsumer.Connect()
	if err != nil {
		rmqConsumer.OnError(err, "Failed to connect to RabbitMQ")
		panic(err)
	}

	wg.Add(len(configs) + 1)
	for _, conf := range configs {

		rmqConsumer.Queue = conf.QueueName
		err := rmqConsumer.QueueDeclare()
		if err != nil {
			rmqConsumer.OnError(err, "Failed to connect to RabbitMQ while consuming queue")
			panic(err)
		}

		go func(consumer controllers.RMQSpec) {
			log.Info().Msgf("Run consuming queue '%v'", consumer.Queue)
			consumer.ConsumeMessages()
		}(rmqConsumer)

	}

	go controllers.ProduceMessage()
	go controllers.Predict()

	wg.Wait()

	//go rmqConsumer.ConsumeMessages()

}
