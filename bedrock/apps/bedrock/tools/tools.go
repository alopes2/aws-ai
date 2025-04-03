package tools

func GetWeather(location string) string {
	return "It is sunny in " + location
}

const (
	ToolGetWeather = "GetWeather"
)

func GetWeatherToolSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"location": map[string]interface{}{
				"type":        "string",
				"description": "The city or location to get the weather for. Example locations are Berlin, New York, Paris",
			},
		},
		"required": []string{"location"},
	}
}
