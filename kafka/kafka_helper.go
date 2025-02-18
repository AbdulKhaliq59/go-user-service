package kafka

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaHelper struct {
	producer *kafka.Producer
	consumer *kafka.Consumer
}

func NewKafkaHelper() (*KafkaHelper, error) {
	brokerURL := os.Getenv("KAFKA_BROKER_URL")
	if brokerURL == "" {
		return nil, fmt.Errorf("KAFKA_BROKER_URL environment variable is not set")
	}

	// Create Producer
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokerURL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	// Create Consumer
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokerURL,
		"group.id":          os.Getenv("KAFKA_GROUP_ID"),
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		producer.Close()
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	return &KafkaHelper{
		producer: producer,
		consumer: consumer,
	}, nil
}

func (k *KafkaHelper) Send(data interface{}, transactionName string, topic string) error {
	// Convert data to JSON
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	// Create message with headers
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          dataBytes,
		Headers: []kafka.Header{
			{Key: "appName", Value: []byte(os.Getenv("APPLICATION_NAME"))},
		},
	}

	// Send the message
	if err := k.producer.Produce(message, nil); err != nil {
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}

	// Wait for delivery report
	e := <-k.producer.Events()
	if m, ok := e.(*kafka.Message); ok {
		if m.TopicPartition.Error != nil {
			return fmt.Errorf("delivery failed: %w", m.TopicPartition.Error)
		}
	}

	return nil
}

func (k *KafkaHelper) Close() {
	if k.producer != nil {
		k.producer.Close()
	}
	if k.consumer != nil {
		k.consumer.Close()
	}
}
