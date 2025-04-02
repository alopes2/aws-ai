package tools

func GetWeather(city string) string {
	return "It is sunny in " + city
}

func GetWeatherToolSchema() ToolSchema {
	return ToolSchema{
		Json: ToolJsonSchema{
			Type: "object",
			Properties: map[string]any{
				"city": map[string]string{
					"type":        "string",
					"description": "The city to get the current weather for. Example cities are Berlin, New York, Paris",
				},
			},
			Required: []string{
				"city",
			},
		},
	}
}

type ToolSchema struct {
	Json ToolJsonSchema `json:"json"`
}

type ToolJsonSchema struct {
	Type       string         `json:"type"`
	Required   []string       `json:"required"`
	Properties map[string]any `json:"properties"`
}
