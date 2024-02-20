package main

import "html/template"

var weatherTemplate *template.Template

func init() {
	weatherTemplate = template.Must(template.New("index").Parse(`
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
	`))
}
