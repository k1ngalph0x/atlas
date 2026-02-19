package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/k1ngalph0x/atlas/services/issue-service/api"
	"github.com/k1ngalph0x/atlas/services/issue-service/models"
	"github.com/segmentio/kafka-go"
)


func Consume(handler *api.IssueHandler){
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: handler.Config.KAFKA.Brokers,
		Topic: "atlas-events",
		GroupID: "issue-consumers",
	})

	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil{
			log.Println("Kafka read error:", err)
			continue
		}

		var event models.Event
		err = json.Unmarshal(msg.Value, &event)
		if err != nil{
			log.Println("Invalid event:", err)
			continue
		}

		handler.ProcessEvents(event)
		
	}
}

