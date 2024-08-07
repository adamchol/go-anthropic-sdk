package anthropic

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageWithTextOmitEmpty(t *testing.T) {
	json1, err := json.Marshal(
		MessageRequest{
			Model: "mock",
			Messages: []InputMessage{
				{
					Role:    "user",
					Content: "content",
				},
			},
		},
	)
	assert.NoError(t, err)

	const expectedJSON1 = `{"model":"mock","messages":[{"role":"user","content":"content"}],"max_tokens":0}`

	assert.Equal(t, expectedJSON1, string(json1))
}

func TestMessageWithImageOmitEmpty(t *testing.T) {
	json2, err := json.Marshal(
		MessageRequest{
			Model: "mock",
			Messages: []InputMessage{
				{
					Role: "user",
					ContentBlocks: []ContentBlock{
						{
							Type: "image",
							Source: ImageSource{
								Type:      ImageSourceType,
								MediaType: ImagePNGMediaType,
								Data:      "data",
							},
						},
					},
				},
			},
			MaxTokens: 2000,
		},
	)
	assert.NoError(t, err)

	const expectedJSON2 = `{"model":"mock","messages":[{"role":"user","content":[{"type":"image","source":{"type":"base64","media_type":"image/png","data":"data"}}]}],"max_tokens":2000}`

	assert.Equal(t, expectedJSON2, string(json2))
}

func TestMessageWithToolUseOmitEmpty(t *testing.T) {
	json3, err := json.Marshal(
		MessageRequest{
			Model: "mock",
			Messages: []InputMessage{
				{
					Role: "user",
					ContentBlocks: []ContentBlock{
						{
							Type: "tool_use",
							Id:   "tool_id",
							Name: "tool_name",
							Input: map[string]interface{}{
								"key": "value",
							},
						},
					},
				},
			},
			MaxTokens: 2000,
		},
	)
	assert.NoError(t, err)

	const expectedJSON3 = `{"model":"mock","messages":[{"role":"user","content":[{"type":"tool_use","id":"tool_id","name":"tool_name","input":{"key":"value"}}]}],"max_tokens":2000}`

	assert.Equal(t, expectedJSON3, string(json3))
}

func TestMessageWithToolResultOmitEmpty(t *testing.T) {
	json4, err := json.Marshal(
		MessageRequest{
			Model: "mock",
			Messages: []InputMessage{
				{
					Role: "user",
					ContentBlocks: []ContentBlock{
						{
							Type:      "tool_result",
							ToolUseId: "tool_use_id",
							IsError:   true,
							ToolResultContent: ToolResultContent{
								Type: "text",
								Text: "tool result content text",
							},
						},
					},
				},
			},
		},
	)
	assert.NoError(t, err)

	const expectedJSON4 = `{"model":"mock","messages":[{"role":"user","content":[{"type":"tool_result","tool_use_id":"tool_use_id","is_error":true,"content":{"type":"text","text":"tool result content text"}}]}],"max_tokens":0}`

	assert.Equal(t, expectedJSON4, string(json4))

	json5, err := json.Marshal(
		MessageRequest{
			Model: "mock",
			Messages: []InputMessage{
				{
					Role: "assistant",
					ContentBlocks: []ContentBlock{
						{
							Type:      "tool_result",
							ToolUseId: "tool_use_id",
							IsError:   false,
							ToolResultContent: ToolResultContent{
								Type: "image",
								Source: ImageSource{
									Type:      ImageSourceType,
									MediaType: ImagePNGMediaType,
									Data:      "data",
								},
							},
						},
					},
				},
			},
		},
	)
	assert.NoError(t, err)

	const expectedJSON5 = `{"model":"mock","messages":[{"role":"assistant","content":[{"type":"tool_result","tool_use_id":"tool_use_id","content":{"type":"image","source":{"type":"base64","media_type":"image/png","data":"data"}}}]}],"max_tokens":0}`

	assert.Equal(t, expectedJSON5, string(json5))
}

func TestMessageContentDuplicateError(t *testing.T) {
	_, err := json.Marshal(
		MessageRequest{
			Model: "mock",
			Messages: []InputMessage{
				{
					Role:    "user",
					Content: "content",
					ContentBlocks: []ContentBlock{
						{
							Type: "text",
						},
					},
				},
			},
		},
	)
	assert.Error(t, err)
}

// TODO: Tests for CreateRequest method with a mock server
