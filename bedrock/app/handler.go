package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

type Handler struct {
	bedrockClient *bedrockruntime.Client
	modelID       string
}

func (h *Handler) HandleRequest(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Printf("Received event %+v", event)

	var request Request
	err := json.Unmarshal([]byte(event.Body), &request)

	if err != nil {
		log.Fatal("Could not unmarshal request body")
	}

	log.Printf("Got user prompt message %s", request.Prompt)

	if err != nil {
		log.Fatal("Could not marshal model request body")
	}

	response, err := h.callBedrock(request.Prompt, &ctx)

	responseBody, _ := json.Marshal(Response{
		Text: response,
	})

	apiResponse := &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
	}

	return apiResponse, nil
}

func (h *Handler) callBedrock(prompt string, ctx *context.Context) (string, error) {
	log.Printf("Sending input %s to model %s", prompt, h.modelID)

	input := &bedrockruntime.ConverseInput{
		ModelId: aws.String(h.modelID),
		Messages: []types.Message{
			{
				Content: []types.ContentBlock{
					&types.ContentBlockMemberText{
						Value: prompt,
					},
				},
				Role: types.ConversationRoleUser,
			},
		},
		System: []types.SystemContentBlock{
			&types.SystemContentBlockMemberText{
				Value: "You are a weather assistant. You return the weather conditions (for example: sunny, cloudy, windy), the humidity, and the temperature in Celsius.",
			},
		},
	}

	// input := &bedrockruntime.ConverseStreamInput{
	// 	ModelId: aws.String(h.modelID),
	// 	Messages: []types.Message{
	// 		{
	// 			Content: []types.ContentBlock{
	// 				&types.ContentBlockMemberText{
	// 					Value: prompt,
	// 				},
	// 			},
	// 			Role: types.ConversationRoleUser,
	// 		},
	// 	},
	// }

	response, err := h.bedrockClient.Converse(*ctx, input)
	// response, err := h.bedrockClient.ConverseStream(*ctx, input, func(o *bedrockruntime.Options) {})

	if err != nil {
		log.Printf("Failed to invoke model with error %s", err.Error())
		return "", err
	}

	outputMessage, _ := response.Output.(*types.ConverseOutputMemberMessage)

	text, _ := outputMessage.Value.Content[0].(*types.ContentBlockMemberText)

	return text.Value, nil
}

func NewHandler(config aws.Config, modelID string) *Handler {
	return &Handler{
		bedrockClient: bedrockruntime.NewFromConfig(config),
		modelID:       modelID,
	}
}
