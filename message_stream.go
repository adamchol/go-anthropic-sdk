package anthropic

import "encoding/json"

type StreamEvent struct {
	Type  string `json:"type"`
	Index int    `json:"index"`

	// For message_start type
	Message MessageResponse `json:"message,omitempty"`

	// For content_block_start type
	ContentBlock ContentBlock `json:"content_block,omitempty"`

	// For content_block_delta type
	Delta ContentBlock `json:"delta,omitempty"`
}

type DeltaContent struct {
	Type string `json:"type"`

	Text string `json:"text,omitempty"`

	PartialJSON string `json:"partial_json,omitempty"`
}

func (se StreamEvent) MarshalJSON() ([]byte, error) {
	type alias StreamEvent

	temp := struct {
		alias
		Message      *MessageResponse `json:"message,omitempty"`
		ContentBlock *ContentBlock    `json:"content_block,omitempty"`
	}{
		alias: alias(se),
	}

	if se.Message.ID != "" {
		temp.Message = &se.Message
	}

	if se.ContentBlock.Type != "" {
		temp.ContentBlock = &se.ContentBlock
	}

	return json.Marshal(temp)
}
