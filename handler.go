package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const openWeatherMapAPIKey = "your_api_key"

type WeatherData struct {
	City        string
	Temperature float64
	Weather     string
	Condition   string
	Latitude    float64
	Longitude   float64
}

func getCondition(temperature float64, weatherDescription string) string {
	if temperature >= 30.0 {
		return "hot"
	} else if temperature <= 10.0 {
		return "cold"
	}

	// Additional conditions based on weather description
	switch weatherDescription {
	case "clear sky":
		return "clear sky"
	case "few clouds", "scattered clouds":
		return "partly cloudy"
	case "broken clouds", "overcast clouds":
		return "cloudy"
	case "light rain", "moderate rain", "heavy intensity rain", "very heavy rain":
		return "rainy"
	case "light snow", "moderate snow", "heavy snow", "sleet":
		return "snowy"
	default:
		return "unknown"
	}
}

func parseCoordinates(input string) (float64, float64, error) {
	// Use regular expression to extract latitude and longitude
	re := regexp.MustCompile(`(\d+\.\d+)\s*°\s*([NS])\s*,\s*(\d+\.\d+)\s*°\s*([EW])`)
	match := re.FindStringSubmatch(input)
	if len(match) != 5 {
		return 0, 0, fmt.Errorf("Invalid input format. Please provide latitude and longitude in the format '32.6518° N, 96.9083° W'")
	}

	// Parse latitude and longitude values
	lat, _ := strconv.ParseFloat(match[1], 64)
	lon, _ := strconv.ParseFloat(match[3], 64)

	// Adjust latitude based on hemisphere (North/South)
	if strings.ToUpper(match[2]) == "S" {
		lat = -lat
	}

	// Adjust longitude based on hemisphere (East/West)
	if strings.ToUpper(match[4]) == "W" {
		lon = -lon
	}

	return lat, lon, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Display the HTML form for input
		tmpl, err := template.New("index").Parse(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>Weather Report</title>
			</head>
			<body>
				<h1>Weather Report</h1>
				<form method="post" action="/">
					<label for="coordinates">Enter coordinates (e.g., 32.6518° N, 96.9083° W):</label>
					<input type="text" id="coordinates" name="coordinates" required>
					<button type="submit">Get Weather</button>
				</form>
			</body>
			</html>
		`)
		if err != nil {
			http.Error(w, "Error parsing HTML template", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Error executing HTML template", http.StatusInternalServerError)
			return
		}
	} else if r.Method == "POST" {
		// Retrieve coordinates from the form submission
		coordinatesStr := r.FormValue("coordinates")
		coordinatesStr = strings.TrimSpace(coordinatesStr)

		// Parse latitude and longitude from user input
		latitude, longitude, err := parseCoordinates(coordinatesStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Define the OpenWeatherMap API endpoint URL with user-provided latitude and longitude
		apiURL := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s", latitude, longitude, openWeatherMapAPIKey)
		// Make a GET request to the OpenWeatherMap API
		response, err := http.Get(apiURL)
		if err != nil {
			http.Error(w, "Error making GET request", http.StatusInternalServerError)
			return
		}
		defer response.Body.Close()

		// Check the status code
		if response.StatusCode != http.StatusOK {
			http.Error(w, fmt.Sprintf("Error: %s", response.Status), http.StatusInternalServerError)
			return
		}
		// Read the response body
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			http.Error(w, "Error reading response body", http.StatusInternalServerError)
			return
		}
		// Parse the JSON response
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			http.Error(w, "Error decoding JSON", http.StatusInternalServerError)
			return
		}
		// Create a WeatherData struct with relevant information
		weatherData := WeatherData{
			City:        result["name"].(string),
			Temperature: result["main"].(map[string]interface{})["temp"].(float64) - 273.15,
			Weather:     result["weather"].([]interface{})[0].(map[string]interface{})["description"].(string),
			Condition:   getCondition(result["main"].(map[string]interface{})["temp"].(float64)-273.15, result["weather"].([]interface{})[0].(map[string]interface{})["description"].(string)),
			Latitude:    latitude,
			Longitude:   longitude,
		}
		// Render the HTML template with the weather data
		tmpl, err := template.New("index").Parse(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>Weather Report</title>
			</head>
			<body>
				<h1>Weather Report for {{.City}}</h1>
				<p>Temperature: {{.Temperature}} °C</p>
				<p>Weather: {{.Weather}}</p>
				<p>Condition: {{.Condition}}</p>
				<p>City: {{.City}}</p>
				<p>Location: Latitude {{.Latitude}}, Longitude {{.Longitude}}</p>
			</body>
			</html>
		`)
		if err != nil {
			http.Error(w, "Error parsing HTML template", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, weatherData)
		if err != nil {
			http.Error(w, "Error executing HTML template", http.StatusInternalServerError)
			return
		}
	}
}
