package kafka

import "fmt"

type EventHelper struct {
	kafkaHelper *KafkaHelper
}

func NewEventHelper(kafkaHelper *KafkaHelper) *EventHelper {
	return &EventHelper{
		kafkaHelper: kafkaHelper,
	}
}

func (e *EventHelper) SendEvent(message string, data interface{}, topic string) error {
	message = fmt.Sprintf("[Payment GTW] %s", message)
	err := e.kafkaHelper.Send(data, message, topic)
	if err != nil {
		return fmt.Errorf("failed to send event: %w", err)
	}
	return nil
}
