package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

type Handler struct {
	bedrockClient              *bedrockruntime.Client
	modelID                    string
	apiGatewayManagementClient *apigatewaymanagementapi.Client
}

func (h *Handler) HandleRequest(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Printf("Received event %+v", event)

	var request Request
	err := json.Unmarshal([]byte(event.Body), &request)

	if err != nil {
		log.Fatal("Could not unmarshal request body")
	}

	if request.Prompt == "" {
		return &events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "",
		}, nil
	}

	log.Printf("Got user prompt message %s", request.Prompt)

	if err != nil {
		log.Fatal("Could not marshal model request body")
	}

	response, err := h.callBedrock(request.Prompt, &ctx, event.RequestContext.ConnectionID)

	responseBody, _ := json.Marshal(Response{
		Messages: response.Content,
	})

	apiResponse := &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
	}

	return apiResponse, nil
}

func (h *Handler) callBedrock(prompt string, ctx *context.Context, connectionID string) (*Message, error) {
	log.Printf("Sending input %s to model %s", prompt, h.modelID)

	// input := &bedrockruntime.ConverseInput{
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
	// 	System: []types.SystemContentBlock{
	// 		&types.SystemContentBlockMemberText{
	// 			Value: `You are a technology expert named Tony.
	// 			You answer technology related questions in a friendly and casual tone.
	// 			You break down complex topics into easy-to-understand explanations.
	// 			It's ok to not know the answer, but try your best to point to where the user might find more information about the topic.`,
	// 		},
	// 	},
	// }

	input := &bedrockruntime.ConverseStreamInput{
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
				Value: `You are a technology expert named Tony.
				You answer technology related questions in a friendly and casual tone.
				You break down complex topics into easy-to-understand explanations.
				It's ok to not know the answer, but try your best to point to where the user might find more information about the topic.`,
			},
		},
	}

	// response, err := h.bedrockClient.Converse(*ctx, input)
	response, err := h.bedrockClient.ConverseStream(*ctx, input)

	if err != nil {
		log.Printf("Failed to invoke model with error %s", err.Error())
		return nil, err
	}

	outputMessage := response.GetStream().Events()

	msg := h.handleOutput(outputMessage, connectionID, ctx)

	return msg, nil
}

func (h *Handler) handleOutput(outputMessage <-chan types.ConverseStreamOutput, connectionID string, ctx *context.Context) *Message {
	var combinedResult string

	msg := Message{}

	for event := range outputMessage {
		switch e := event.(type) {
		case *types.ConverseStreamOutputMemberContentBlockDelta:
			textResponse := e.Value.Delta.(*types.ContentBlockDeltaMemberText)

			h.SendWebSocketMessageToConnection(ctx, textResponse.Value, "content", connectionID)

			combinedResult = combinedResult + textResponse.Value

		case *types.ConverseStreamOutputMemberMessageStart:
			log.Print("Message start")

			h.SendWebSocketMessageToConnection(ctx, "", "message_start", connectionID)

			msg.Role = string(e.Value.Role)

		case *types.ConverseStreamOutputMemberMessageStop:
			log.Printf("Message stop. Reason: %s", e.Value.StopReason)

			h.SendWebSocketMessageToConnection(ctx, "", "message_stop", connectionID)

		case *types.ConverseStreamOutputMemberContentBlockStart:
			log.Print("Content block start")
			h.SendWebSocketMessageToConnection(ctx, "", "content_start", connectionID)

			combinedResult = ""

		case *types.ConverseStreamOutputMemberContentBlockStop:
			log.Print("Content block stop")
			h.SendWebSocketMessageToConnection(ctx, "", "content_stop", connectionID)

			msg.Content = append(msg.Content, combinedResult)

		case *types.ConverseStreamOutputMemberMetadata:
			log.Printf("Metadata %+v", e.Value)
			h.SendWebSocketMessageToConnection(ctx, fmt.Sprintf("%+v", e.Value), "metadata", connectionID)

		case *types.UnknownUnionMember:
			log.Printf("unknown tag: %s", e.Tag)

		default:
			log.Printf("Received unexpected event type: %T", e)
		}
	}

	return &msg
}

func (h *Handler) SendWebSocketMessageToConnection(ctx *context.Context, textResponse string, event string, connectionID string) {
	data, _ := json.Marshal(WebSocketResponse{Event: event, Data: textResponse})

	log.Printf("Sending to websocket %+v", data)

	websocketInput := &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(connectionID),
		Data:         []byte(data),
	}

	response, err := h.apiGatewayManagementClient.PostToConnection(*ctx, websocketInput)

	if err != nil {
		log.Printf("ERROR %+v", err)
	}

	log.Printf("Output from websocket %+v", response)
}

func NewHandler(config aws.Config) *Handler {
	modelID := os.Getenv("MODEL_ID")
	apiGatewayEndpoint := os.Getenv("API_GATEWAY_ENDPOINT")

	return &Handler{
		bedrockClient: bedrockruntime.NewFromConfig(config),
		modelID:       modelID,
		apiGatewayManagementClient: apigatewaymanagementapi.NewFromConfig(config, func(o *apigatewaymanagementapi.Options) {
			o.BaseEndpoint = &apiGatewayEndpoint
		}),
	}
}
