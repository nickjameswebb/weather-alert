package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type WeatherResponse struct {
	Weather []Weather `json:"weather"`
	Main    Main      `json:"main"`
}

type Main struct {
	Temp    float32 `json:"temp"`
	TempMin float32 `json:"temp_min"`
	TempMax float32 `json:"temp_max"`
}

type Weather struct {
	Description string `json:"description"`
}

// func newKafkaWriter(kafkaURL, topic string) *kafka.Writer {
// 	return &kafka.Writer{
// 		Addr:     kafka.TCP(kafkaURL),
// 		Topic:    topic,
// 		Balancer: &kafka.LeastBytes{},
// 	}
// }

func getWeatherByZip(zip string, openWeatherMapAPIKey string) WeatherResponse {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?zip=%s&appid=%s",
		url.QueryEscape(zip),
		url.QueryEscape(openWeatherMapAPIKey),
	)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// log.Println(string(body))

	var weatherResponse WeatherResponse
	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		log.Fatalln(err)
	}

	return weatherResponse
}

func main() {
	openWeatherMapAPIKey := os.Getenv("OPEN_WEATHER_MAP_API_KEY")
	zipCode := os.Getenv("ZIP_CODE")
	// kafkaURL := os.Getenv("KAFKA_URL")
	// topic := os.Getenv("KAFKA_TOPIC")
	// writer := newKafkaWriter(kafkaURL, topic)
	// defer writer.Close()

	log.Println("==> starting producer...")

	_ = getWeatherByZip(zipCode, openWeatherMapAPIKey)
}
