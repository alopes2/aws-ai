package main

type Request struct {
	Prompt string `json:"prompt"`
}

type Response struct {
	Messages []string `json:"messages"`
}

type WebSocketMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

const (
	BedrockEventContent      string = "content"
	BedrockEventMessageStart string = "message_start"
	BedrockEventMessageStop  string = "message_stop"
	BedrockEventContentStart string = "content_start"
	BedrockEventContentStop  string = "content_stop"
	BedrockEventMetadata     string = "metadata"
)

type ToolUse struct {
	ToolUseID string
	Name      string
	Input     string
}
