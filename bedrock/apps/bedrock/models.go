package main

type Request struct {
	Prompt string `json:"prompt"`
}

type Response struct {
	Messages []string `json:"messages"`
}

type Message struct {
	Role    string
	Content []string
}

type WebSocketResponse struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}
