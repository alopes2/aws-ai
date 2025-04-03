package tools

func GetWeather(city string) string {
	return "It is sunny in " + city
}

func GetWeatherToolSchema() ToolSchema {
	return ToolSchema{
		Json: ToolJsonSchema{
			Type: "object",
			Properties: ToolJsonSchemaProperties{
				"city": ToolJsonSchemaProperty{
					Type:        "string",
					Description: "The city to get the current weather for. Example cities are Berlin, New York, Paris",
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
	Type       string                   `json:"type"`
	Required   []string                 `json:"required"`
	Properties ToolJsonSchemaProperties `json:"properties"`
}

type ToolJsonSchemaProperties map[string]ToolJsonSchemaProperty

type ToolJsonSchemaProperty struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}
