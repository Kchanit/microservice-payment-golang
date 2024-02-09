package broker

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func KafkaProducer(topic string, key string, content map[string]interface{}) error {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "119.59.99.166:9092"})
	if err != nil {
		log.Fatal(err)
		return err
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)

	// content := map[string]interface{}{
	// 	"message": "Hello World!",
	// }
	value, err := json.Marshal(content)

	if err != nil {
		log.Fatal(err)
		return err
	}

	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          value,
	}, nil)

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)

	return nil
}
