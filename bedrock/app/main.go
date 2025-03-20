package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type response struct {
	Message string `json:"message"`
}

type handler struct{}

func (h *handler) handleRequest(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Printf("Received event %s")
	body, _ := json.Marshal(response{Message: "Hello from Lambda!"})

	message := &events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}

	return message, nil
}

func main() {
	h := handler{}
	lambda.Start(h.handleRequest)
}
