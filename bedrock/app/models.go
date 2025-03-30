package main

type Request struct {
	Prompt string `json:"prompt"`
}

type Response struct {
	Message *Message `json:"message"`
}

type Message struct {
	Role    string   `json:"role"`
	Content []string `json:"content"`
}
