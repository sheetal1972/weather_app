package main

import (
	"fmt"
	stdhttp "net/http"
	"weather_app/internal/interfaces/http"
)

func main() {
	stdhttp.HandleFunc("/", http.WeatherHandler)
	fmt.Println("Server listening on :8080...")
	stdhttp.ListenAndServe(":8080", nil)
}
