# weather_app

# please follow the steps to build the app in your local and run the below commands in your cli
git clone https://github.com/sheetal1972/weather_app.git ( clone the Git Repository )

## Installation

To install and run the Weather App, follow these steps:

1. **Clone the repository**:
    ```sh
    git clone <repository-url>
    cd weather-app
    ```

2. **Set up Go modules**:
    ```sh
    go mod tidy
    ```

3. **Build the application**:
    ```sh
    go build ./cmd/weather
    ```

4. **Run the application**:
    ```sh
    ./weather
    ```

now open a web browser http://localhost:8080 you are able see the weather Report page with box where user need to enter the coordinates 
with the same formate given in example

**Note:**  Use your own API_KEY in handler.go file line number 14 openWeatherMapAPIKey = "your_api_key"


![image](https://github.com/sheetal1972/weather_app/assets/160625825/8b4e5c4e-232d-4ea4-998b-e04f9fbd0a88)

Enter coordinates and click on **Get Weather** button

![image](https://github.com/sheetal1972/weather_app/assets/160625825/21d6bcd8-60db-4a3f-b492-87ece60df5d2)


Expected output:

![image](https://github.com/sheetal1972/weather_app/assets/160625825/5fc97dc1-7489-4c95-93c7-de0cd3b2db63)
