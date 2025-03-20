package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

type response struct {
	Body       string `json:"body"`
	StatusCode int    `json:"statusCode"`
}

type message struct {
	Message string `json:"message"`
}

type handler struct{}

func (h *handler) handleRequest(ctx context.Context, event json.RawMessage) (*response, error) {
	log.Printf("Received event %s", event)
	body, _ := json.Marshal(message{Message: "Hello from Lambda!"})

	message := &response{
		Body:       string(body),
		StatusCode: 200,
	}

	return message, nil
}

func main() {
	h := handler{}
	lambda.Start(h.handleRequest)
}
