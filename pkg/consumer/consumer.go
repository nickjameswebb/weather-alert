package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/slack-go/slack"
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

func writeWeatherMessage(weatherMessage kafka.Message, slackWebhookURL string) {
	log.Printf("==> message at topic:%v partition:%v offset:%v	%s = %s\n",
		weatherMessage.Topic,
		weatherMessage.Partition,
		weatherMessage.Offset,
		string(weatherMessage.Key),
		string(weatherMessage.Value),
	)
	if slackWebhookURL != "" {
		log.Println("==> also sending message to slack webhook")
		attachment := slack.Attachment{
			Color: "good",
			Text:  string(weatherMessage.Value),
			Ts:    json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
		}
		msg := slack.WebhookMessage{
			Attachments: []slack.Attachment{attachment},
		}
		err := slack.PostWebhook(slackWebhookURL, &msg)
		if err != nil {
			log.Fatalln(err)
		}
	}
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
	slackWebhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	consumerGroup := fmt.Sprintf("weather-consumers-%s", zipCode)

	log.Printf("==> building reader (topic=%s,consumerGroup=%s)\n", zipCode, consumerGroup)

	reader := newKafkaReader(kafkaURL, zipCode, consumerGroup)

	defer reader.Close()

	log.Println("==> reading messages")

	for {
		weatherMessage, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("==> reading from topic %s failed: %s\n", zipCode, err)
		}
		writeWeatherMessage(weatherMessage, slackWebhookURL)
	}
}
