package anthropic

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageOmitEmpty(t *testing.T) {
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

	assert.Equal(t, string(json1), expectedJSON1)

	json2, err := json.Marshal(
		MessageRequest{
			Model: "mock",
			Messages: []InputMessage{
				{
					Role: "user",
					ContentBlocks: []ContentBlock{
						{
							Type: "image",
							Source: &ImageSource{
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

	assert.Equal(t, string(json2), expectedJSON2)

	_, err = json.Marshal(
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
