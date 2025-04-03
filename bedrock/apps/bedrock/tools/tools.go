package tools

func GetWeather(location string) string {
	return "It is sunny in " + location
}

// func GetWeatherToolSchema() ToolJsonSchema {
func GetWeatherToolSchema() map[string]interface{} {
	// return ToolJsonSchema{
	// 	Type: "object",
	// 	Properties: ToolJsonSchemaProperties{
	// 		"city": ToolJsonSchemaProperty{
	// 			Type:        "string",
	// 			Description: "The city to get the current weather for. Example cities are Berlin, New York, Paris",
	// 		},
	// 	},
	// 	Required: []string{
	// 		"city",
	// 	},
	// }

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
