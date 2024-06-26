package domain

type WeatherData struct {
	City        string
	Temperature float64
	Weather     string
	Condition   string
	Latitude    float64
	Longitude   float64
}

func GetCondition(temperature float64, weatherDescription string) string {
	if temperature >= 30.0 {
		return "hot"
	} else if temperature <= 10.0 {
		return "cold"
	}

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
