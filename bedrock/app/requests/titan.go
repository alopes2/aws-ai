package requests

type TitanTextRequest struct {
	InputText            string                    `json:"inputText"`
	TextGenerationConfig TitanTextGenerationConfig `json:"textGenerationConfig"`
}

type TitanTextGenerationConfig struct {
	Temperature   float64  `json:"temperature"`
	TopP          float64  `json:"topP"`
	MaxTokenCount int      `json:"maxTokenCount"`
	StopSequences []string `json:"stopSequences,omitempty"`
}

type TitanTextResponse struct {
	InputTextTokenCount int           `json:"inputTextTokenCount"`
	Results             []TitanResult `json:"results"`
}

type TitanResult struct {
	TokenCount       int    `json:"tokenCount"`
	OutputText       string `json:"outputText"`
	CompletionReason string `json:"completionReason"`
}
