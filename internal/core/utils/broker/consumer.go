package broker

import (
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Event func(*kafka.Message)

func KafkaConsumer(topic []string, group string, action Event) error {

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
		"group.id":          group,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return err
	}

	c.SubscribeTopics(topic, nil)

	// A signal handler or similar could be used to set this to false to break the loop.
	run := true

	for run {
		msg, err := c.ReadMessage(time.Second)
		if err == nil {
			action(msg)
			// fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			// fmt.Printf("Type of value %T\n", msg.Value)
		} else if !err.(kafka.Error).IsTimeout() {
			// The client will automatically try to recover from all errors.
			// Timeout is not considered an error because it is raised by
			// ReadMessage in absence of messages.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)

			return nil
		}
	}

	c.Close()

	return nil

}
