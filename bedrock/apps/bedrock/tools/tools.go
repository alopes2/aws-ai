package tools

func GetWeather(city string) string {
	return "It is sunny in " + city
}

func GetWeatherToolSchema() map[string]map[string]any {
	return map[string]map[string]any{
		"json": {
			"type": "object",
			"properties": map[string]any{
				"city": map[string]string{
					"type":        "string",
					"description": "The city to get the current weather for. Example cities are Berlin, New York, Paris",
				},
			},
			"required": []string{
				"city",
			},
		},
	}
}
