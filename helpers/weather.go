package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetWeather(city string, key string) (string, error) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, key)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get weather data")
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	main := data["main"].(map[string]interface{})
	weather := data["weather"].([]interface{})[0].(map[string]interface{})
	temp := main["temp"].(float64)
	description := weather["description"].(string)

	return fmt.Sprintf("The weather in %s: %0.1f°C, %s.", city, temp, description), nil
}