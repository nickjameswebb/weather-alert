package main

import (
	"context"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

func newKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func main() {
	openWeatherMapAPIKey := os.Getenv("OPEN_WEATHER_MAP_API_KEY")
	zipCode := os.Getenv("ZIP_CODE")
	kafkaURL := os.Getenv("KAFKA_URL")
	writer := newKafkaWriter(kafkaURL, zipCode)
	defer writer.Close()

	log.Println("==> starting producer...")

	weatherResponse, err := NewWeather(zipCode, openWeatherMapAPIKey)
	if err != nil {
		log.Fatalln("==> error getting weather: ", err)
	}

	log.Printf("==> writing to topic %s...\n", zipCode)

	err = writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:   []byte(zipCode),
			Value: []byte(weatherResponse.ToString()),
		},
	)
	if err != nil {
		log.Fatalf("==> error writing to topic %s: %s\n", zipCode, err)
	}

	log.Printf("==> writing to topic %s succeeded...\n", zipCode)
}
