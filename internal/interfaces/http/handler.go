package http

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"weather_app/internal/usecase"
)

func parseCoordinates(input string) (float64, float64, error) {
	re := regexp.MustCompile(`(\d+\.\d+)\s*°\s*([NS])\s*,\s*(\d+\.\d+)\s*°\s*([EW])`)
	match := re.FindStringSubmatch(input)
	if len(match) != 5 {
		return 0, 0, fmt.Errorf("Invalid input format. Please provide latitude and longitude in the format '32.6518° N, 96.9083° W'")
	}

	lat, _ := strconv.ParseFloat(match[1], 64)
	lon, _ := strconv.ParseFloat(match[3], 64)

	if strings.ToUpper(match[2]) == "S" {
		lat = -lat
	}

	if strings.ToUpper(match[4]) == "W" {
		lon = -lon
	}

	return lat, lon, nil
}

func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.New("index").Parse(`
            <!DOCTYPE html>
            <html>
            <head>
                <title>Weather Report</title>
            </head>
            <body>
			<h1>Welcome to Weather Reports</h1>
            <p>Please enter coordinates to get the weather information for that location.</p>
                <form method="POST">
                    <label for="coords">Enter coordinates (e.g., 32.6518° N, 96.9083° W):</label>
                    <input type="text" id="coords" name="coords">
                    <input type="submit" value="Get Weather">
                </form>
            </body>
            </html>
        `)
		if err != nil {
			http.Error(w, "Error parsing HTML template", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, nil)
		return
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusInternalServerError)
			return
		}

		coords := r.FormValue("coords")
		latitude, longitude, err := parseCoordinates(coords)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		weatherData, err := usecase.FetchWeatherData(latitude, longitude)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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
