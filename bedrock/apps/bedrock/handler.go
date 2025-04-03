package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/alopes2/aws-ai/bedrock/tools"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/document"
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

	if err != nil {
		log.Fatal("Failed to call bedrock")
	}

	var messages []string
	for _, content := range response.Content {
		if textContent, ok := content.(*types.ContentBlockMemberText); ok {
			messages = append(messages, textContent.Value)
		}
	}

	responseBody, _ := json.Marshal(Response{
		Messages: messages,
	})

	apiResponse := &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
	}

	return apiResponse, nil
}

func (h *Handler) callBedrock(prompt string, ctx *context.Context, connectionID string) (*types.Message, error) {
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

	messages := []types.Message{
		{
			Content: []types.ContentBlock{
				&types.ContentBlockMemberText{
					Value: prompt,
				},
			},
			Role: types.ConversationRoleUser,
		},
	}

	msg, err := h.newFunction(ctx, &messages, connectionID)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (h *Handler) newFunction(ctx *context.Context, messages *[]types.Message, connectionID string) (*types.Message, error) {
	input := &bedrockruntime.ConverseStreamInput{
		ModelId:  aws.String(h.modelID),
		Messages: *messages,
		System: []types.SystemContentBlock{
			&types.SystemContentBlockMemberText{
				Value: `You are a technology expert named Tony.
				You answer technology related questions in a friendly and casual tone.
				You break down complex topics into easy-to-understand explanations.
				It's ok to not know the answer, but try your best to point to where the user might find more information about the topic.
				You are also allowed to get the current weather.`,
			},
		},
		ToolConfig: &types.ToolConfiguration{
			Tools: []types.Tool{
				&types.ToolMemberToolSpec{
					Value: types.ToolSpecification{
						InputSchema: &types.ToolInputSchemaMemberJson{
							Value: document.NewLazyDocument(tools.GetWeatherToolSchema()),
						},
						Name:        aws.String("GetWeather"),
						Description: aws.String("Get the current weather for a city or location. It returns the weather simplified with the city in the response. Examples: 'It is sunny in Berlin', 'It is raining in Curitiba"),
					},
				},
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

	result, stopReason := h.handleOutput(outputMessage, connectionID, ctx)

	if stopReason == types.StopReasonToolUse {
		*messages = append(*messages, *result)

		result, err = h.newFunction(ctx, messages, connectionID)
	}

	return result, nil
}

func (h *Handler) handleOutput(outputMessage <-chan types.ConverseStreamOutput, connectionID string, ctx *context.Context) (*types.Message, types.StopReason) {
	var result string
	var toolInput string
	var stopReason types.StopReason

	var msg types.Message
	var toolUse types.ToolUseBlock

	for event := range outputMessage {
		switch e := event.(type) {
		case *types.ConverseStreamOutputMemberMessageStart:
			log.Print("Message start")

			h.SendWebSocketMessageToConnection(ctx, "", BedrockEventMessageStart, connectionID)

			msg.Role = e.Value.Role

		case *types.ConverseStreamOutputMemberContentBlockStart:
			log.Print("Content block start")

			switch blockStart := e.Value.Start.(type) {
			case *types.ContentBlockStartMemberToolUse:
				toolUse.Name = blockStart.Value.Name
				toolUse.ToolUseId = blockStart.Value.ToolUseId
			default:
				h.SendWebSocketMessageToConnection(ctx, "", BedrockEventContentStart, connectionID)

				result = ""
			}

		case *types.ConverseStreamOutputMemberContentBlockDelta:
			switch delta := e.Value.Delta.(type) {
			case *types.ContentBlockDeltaMemberText:
				h.SendWebSocketMessageToConnection(ctx, delta.Value, BedrockEventContent, connectionID)
				result = result + delta.Value

			case *types.ContentBlockDeltaMemberToolUse:
				if delta.Value.Input != nil {
					toolInput = toolInput + *delta.Value.Input
				}
			}

		case *types.ConverseStreamOutputMemberContentBlockStop:
			log.Print("Content block stop")

			if toolUse.Input != nil {
				if jsonBytes, err := json.Marshal(toolInput); err != nil {
					toolInput = string(jsonBytes)
					toolUse.Input = document.NewLazyDocument(toolInput)
					msg.Content = append(msg.Content, &types.ContentBlockMemberToolUse{
						Value: toolUse,
					})
				}
				toolUse = types.ToolUseBlock{}
			} else {
				h.SendWebSocketMessageToConnection(ctx, "", BedrockEventContentStop, connectionID)

				msg.Content = append(msg.Content, &types.ContentBlockMemberText{
					Value: result,
				})
			}

		case *types.ConverseStreamOutputMemberMetadata:
			log.Printf("Metadata %+v", e.Value)
			h.SendWebSocketMessageToConnection(ctx, fmt.Sprintf("%+v", e.Value), BedrockEventMetadata, connectionID)

		case *types.ConverseStreamOutputMemberMessageStop:
			log.Printf("Message stop. Reason: %s", e.Value.StopReason)
			stopReason = e.Value.StopReason
			h.SendWebSocketMessageToConnection(ctx, string(e.Value.StopReason), BedrockEventMessageStop, connectionID)

		case *types.UnknownUnionMember:
			log.Printf("unknown tag: %s", e.Tag)

		default:
			log.Printf("Received unexpected event type: %T", e)
		}
	}

	return &msg, stopReason
}

func (h *Handler) SendWebSocketMessageToConnection(ctx *context.Context, textResponse string, event string, connectionID string) {
	data, _ := json.Marshal(WebSocketMessage{Event: event, Data: textResponse})

	websocketInput := &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(connectionID),
		Data:         []byte(data),
	}

	_, err := h.apiGatewayManagementClient.PostToConnection(*ctx, websocketInput)

	if err != nil {
		log.Printf("ERROR %+v", err)
	}
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
