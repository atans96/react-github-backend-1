package service

import (
	"os"
	"server_go/src/server/server/routes/types"
	kafka2 "server_go/src/service/kafka"
)

func KafkaProducer(events chan *types.ChangeEvent) {
	kafkaContainer := &kafka2.Container{}
	kafkaMessageChan := kafka2.NewChangeEventKafkaMessageTransformer(os.Getenv("KAFKA_TOPIC_APOLLO")).Transform(events)
	kafkaContainer.GetKafkaClient().Produce(kafkaMessageChan)
}
