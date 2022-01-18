package consumer

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

func newKafkaReader(kafkaURL string, topic string, consumerGroup string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaURL},
		Topic:    topic,
		GroupID:  consumerGroup,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

func StartConsumer() {
	log.Println("==> starting consumer")

	zipCode := os.Getenv("ZIP_CODE")
	if zipCode == "" {
		log.Fatalln("==> missing required env var ZIP_CODE")
	}
	kafkaURL := os.Getenv("KAFKA_URL")
	if kafkaURL == "" {
		log.Fatalln("==> missing required env var KAFKA_URL")
	}
	consumerGroup := fmt.Sprintf("weather-consumers-%s", zipCode)

	log.Printf("==> building reader (topic=%s,consumerGroup=%s)\n", zipCode, consumerGroup)

	reader := newKafkaReader(kafkaURL, zipCode, consumerGroup)

	defer reader.Close()

	log.Println("==> reading messages")

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("==> reading from topic %s failed: %s\n", zipCode, err)
		}
		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
	}
}
