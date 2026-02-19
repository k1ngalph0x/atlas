package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/k1ngalph0x/atlas/services/ingestion-service/config"
	"github.com/segmentio/kafka-go"
)

var Writer *kafka.Writer

func InitKafka(config *config.Config) error {
	conn, err := kafka.Dial("tcp", config.KAFKA.Brokers[0])
	if err != nil {
		return fmt.Errorf("failed to connect to kafka: %w", err)
	}
	defer conn.Close()


	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             "atlas-events",
		NumPartitions:     3,
		ReplicationFactor: 1,
	})
	
	if err != nil {
		log.Printf("Topic creation warning: %v", err)
	}

	Writer = &kafka.Writer{
		Addr:     kafka.TCP(config.KAFKA.Brokers...),
		Topic:    "atlas-events",
		Balancer: &kafka.LeastBytes{},
	}
	return nil
}


func Publish(projectID string, message []byte) error {
	err := Writer.WriteMessages(context.Background(),
		kafka.Message{
			Key: []byte(projectID),
			Value: message,
		},
	)

	if err != nil{
		return err
	}

	return nil
}