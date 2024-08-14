package kafka

import (
	kafkaconfluent "github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
)

type Container struct {
	kafkaProducer *kafkaconfluent.Producer
	kafkaClient   Client
}

func (container *Container) GetKafkaProducer() KafkaProducer {
	if container.kafkaProducer == nil {
		producer, err := kafkaconfluent.NewProducer(&kafkaconfluent.ConfigMap{
			"bootstrap.servers":       os.Getenv("KAFKA_HOST") + ":" + os.Getenv("KAFKA_PORT"),
			"go.produce.channel.size": 10000,
		})
		if err != nil {
			panic(err)
		}
		container.kafkaProducer = producer
	}

	return container.kafkaProducer
}

func (container *Container) GetKafkaClient() Client {
	if container.kafkaClient == nil {
		container.kafkaClient = container.getKafkaBaseClient()
	}

	return container.kafkaClient
}

func (container *Container) getKafkaBaseClient() Client {
	return NewClient(
		container.GetKafkaProducer(),
	)
}
