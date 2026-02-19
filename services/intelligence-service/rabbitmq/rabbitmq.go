package rabbitmq

import (
	"fmt"
	"log"

	"github.com/k1ngalph0x/atlas/services/intelligence-service/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Connect(cfg *config.Config) (*amqp.Connection, *amqp.Channel, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.RABBITMQ.User,
		cfg.RABBITMQ.Password,
		cfg.RABBITMQ.Host,
		cfg.RABBITMQ.Port,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open channel: %w", err)
	}


	_, err = ch.QueueDeclare(
		"ai-analysis-jobs",
		true, 
		false, 
		false, 
		false, 
		nil,   
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	log.Println("Connected to RabbitMQ")
	return conn, ch, nil
}
