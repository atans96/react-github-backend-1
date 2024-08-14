package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"server_go/src/server/server/routes/types"
)

type ChangeEventKafkaMessageTransformer struct {
	topic string
}

func (t *ChangeEventKafkaMessageTransformer) Transform(changeEvents chan *types.ChangeEvent) chan *Message {
	var messageChan = make(chan *Message, len(changeEvents))
	go func() {
		defer close(messageChan)
		for event := range changeEvents {
			documentID, err := event.DocumentID()
			if err != nil {
				fmt.Println("Mongo transformer: Unable to extract document id from event")
				continue
			}

			jsonBytes, err := event.Marshal()
			if err != nil {
				fmt.Println("Mongo transformer: Unable to unmarshal change event to json")
				continue
			}

			fmt.Println("Mongo transformer: Retrieve event")

			messageChan <- &Message{
				TopicPartition: kafka.TopicPartition{Topic: &t.topic, Partition: kafka.PartitionAny},
				Key:            []byte(documentID),
				Value:          jsonBytes,
			}
		}
	}()
	return messageChan
}

func NewChangeEventKafkaMessageTransformer(topic string) *ChangeEventKafkaMessageTransformer {
	return &ChangeEventKafkaMessageTransformer{
		topic: topic,
	}
}
