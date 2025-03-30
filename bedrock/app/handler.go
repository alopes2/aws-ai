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
		Message: response,
	})

	apiResponse := &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
	}

	return apiResponse, nil
}

func (h *Handler) callBedrock(prompt string, ctx *context.Context) (*Message, error) {
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

	var combinedResult string

	msg := Message{}

	for event := range outputMessage {
		switch e := event.(type) {
		case *types.ConverseStreamOutputMemberContentBlockDelta:
			textResponse := e.Value.Delta.(*types.ContentBlockDeltaMemberText)
			combinedResult = combinedResult + textResponse.Value

		case *types.ConverseStreamOutputMemberMessageStart:
			log.Print("Message start")
			msg.Role = string(e.Value.Role)

		case *types.ConverseStreamOutputMemberMessageStop:
			log.Printf("Message stop. Reason: %s", e.Value.StopReason)

		case *types.ConverseStreamOutputMemberContentBlockStart:
			log.Print("Content block start")
			combinedResult = ""

		case *types.ConverseStreamOutputMemberContentBlockStop:
			log.Print("Content block stop")
			msg.Content = append(msg.Content, combinedResult)

		case *types.UnknownUnionMember:
			log.Printf("unknown tag: %s", e.Tag)

		default:
			log.Printf("Received unexpected event type: %T", e)
		}
	}

	return &msg, nil
}

func NewHandler(config aws.Config, modelID string) *Handler {
	return &Handler{
		bedrockClient: bedrockruntime.NewFromConfig(config),
		modelID:       modelID,
	}
}
