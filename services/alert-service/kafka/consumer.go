package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/k1ngalph0x/atlas/services/alert-service/api"
	"github.com/k1ngalph0x/atlas/services/alert-service/models"
	"github.com/segmentio/kafka-go"
)

func Consume(handler *api.AlertHandler){
		reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: handler.Config.KAFKA.Brokers,
		// Topic: "atlas-events",
		// GroupID: "issue-consumers",
		Topic: "issue-updates",  
		GroupID: "alert-consumers",  
	})
	defer reader.Close()

	for 
	{
		msg, err := reader.ReadMessage(context.Background())
		if err != nil{
			log.Println("Kafka read error:", err)
			continue
		}	

		var event models.IssueUpdateEvent
				err = json.Unmarshal(msg.Value, &event)
		if err != nil{
			log.Println("Invalid event:", err)
			continue
		}

		handler.ProcessAlert(event)
	}
}