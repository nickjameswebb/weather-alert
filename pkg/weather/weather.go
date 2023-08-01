package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type WeatherResponse struct {
	Weather []Weather `json:"weather"`
	Main    Main      `json:"main"`
}

type Main struct {
	Temp float64 `json:"temp"`
}

type Weather struct {
	Description string `json:"description"`
}

func NewWeather(zip string, openWeatherMapAPIKey string) (*WeatherResponse, error) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?zip=%s&appid=%s",
		url.QueryEscape(zip),
		url.QueryEscape(openWeatherMapAPIKey),
	)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var weatherResponse WeatherResponse
	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		return nil, err
	}

	return &weatherResponse, nil
}

func (w *WeatherResponse) ToString() string {
	description := ""
	for i := 0; i < len(w.Weather); i++ {
		description += w.Weather[i].Description
		if i != len(w.Weather)-1 {
			description += ", "
		}
	}

	return fmt.Sprintf(
		"Conditions: %s | Temp: ~%.0f°F\n",
		description,
		KelvinToFahrenheit(w.Main.Temp),
		
	)
}

func KelvinToFahrenheit(kelvin float64) float64 {
	return (kelvin-273.15)*(9/5) + 32
}
