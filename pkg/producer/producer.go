package producer

import (
	"context"
	"log"
	"os"

	"github.com/nickjameswebb/weather-alert/pkg/weather"
	"github.com/segmentio/kafka-go"
)

func newKafkaWriter(kafkaURL string, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func StartProducer() {
	log.Println("==> starting producer...")

	openWeatherMapAPIKey := os.Getenv("OPEN_WEATHER_MAP_API_KEY")
	if openWeatherMapAPIKey == "" {
		log.Fatalln("==> missing required env var OPEN_WEATHER_MAP_API_KEY")
	}
	zipCode := os.Getenv("ZIP_CODE")
	if zipCode == "" {
		log.Fatalln("==> missing required env var ZIP_CODE")
	}
	kafkaURL := os.Getenv("KAFKA_URL")
	if kafkaURL == "" {
		log.Fatalln("==> missing required env var KAFKA_URL")
	}

	writer := newKafkaWriter(kafkaURL, zipCode)
	defer writer.Close()

	weatherResponse, err := weather.NewWeather(zipCode, openWeatherMapAPIKey)
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
