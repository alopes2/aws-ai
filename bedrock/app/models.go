package main

type Request struct {
	Prompt string `json:"prompt"`
}

type Response struct {
	Text string `json:"text"`
}
