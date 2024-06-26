package usecase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"weather_app/internal/domain"
)

const openWeatherMapAPIKey = "your_api_key"

func FetchWeatherData(latitude, longitude float64) (*domain.WeatherData, error) {
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s", latitude, longitude, openWeatherMapAPIKey)
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request to OpenWeatherMap: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %s", response.Status)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	weatherData := &domain.WeatherData{
		City:        result["name"].(string),
		Temperature: result["main"].(map[string]interface{})["temp"].(float64) - 273.15,
		Weather:     result["weather"].([]interface{})[0].(map[string]interface{})["description"].(string),
		Condition:   domain.GetCondition(result["main"].(map[string]interface{})["temp"].(float64)-273.15, result["weather"].([]interface{})[0].(map[string]interface{})["description"].(string)),
		Latitude:    latitude,
		Longitude:   longitude,
	}

	return weatherData, nil
}
