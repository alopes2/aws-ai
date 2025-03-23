package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/alopes2/aws-ai/bedrock/requests"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type request struct {
	message string `json:message`
}

type handler struct {
	bedrockClient *bedrockruntime.Client
	modelID       string
}

func (h *handler) handleRequest(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Printf("Received event %+v", event)
	var requestBody request
	err := json.Unmarshal([]byte(event.Body), &requestBody)

	if err != nil {
		log.Fatal("Could not unmarshal request body")
	}

	promptTemplate := "User:%s \nAssistant:"
	inputText := fmt.Sprintf(promptTemplate, requestBody.message)
	modelRequest := requests.TitanTextRequest{
		InputText: inputText,
		TextGenerationConfig: requests.TitanTextGenerationConfig{
			Temperature:   0.5,
			TopP:          0.9,
			MaxTokenCount: 4089,
		},
	}

	modelRequestBody, err := json.Marshal(modelRequest)

	if err != nil {
		log.Fatal("Could not marshal model request body")
	}

	log.Printf("Sending input %d to model %s", inputText, h.modelID)

	input := &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(h.modelID),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        modelRequestBody,
	}

	response, err := h.bedrockClient.InvokeModel(ctx, input)

	if err != nil {
		log.Fatalf("Failed to invoke model with error %s", err.Error())
	}

	var responseBody requests.TitanTextResponse

	err = json.Unmarshal(response.Body, &responseBody)

	if err != nil {
		log.Fatalf("Failed to unmarshal response %s", err.Error())
	}

	log.Printf("Got results %+v", responseBody.Results)

	apiResponseBody, _ := json.Marshal(responseBody.Results[0])

	apiResponse := &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(apiResponseBody),
	}

	return apiResponse, nil
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	h := handler{
		bedrockClient: bedrockruntime.NewFromConfig(cfg),
		modelID:       os.Getenv("MODEL_ID"),
	}
	lambda.Start(h.handleRequest)
}
