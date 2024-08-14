package kafka

import (
	Kafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

// Message is used over a channel that is filled by kafkaWatch transformer
type Message struct {
	Headers        []Header
	TopicPartition Kafka.TopicPartition
	Key            []byte
	Value          []byte
}

// Header represents a message header
type Header struct {
	Key   string // Header name (utf-8 string)
	Value []byte // Header value (nil, empty, or binary)
}
