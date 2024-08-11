package anthropic

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

const (
	MessageRoleAssistant = "assistant"
	MessageRoleUser      = "user"
)

const messagesSuffix = "/messages"

var (
	ErrContentFieldsMisused             = errors.New("can't use both Content and ContentBlocks properties simultaneously")
	ErrChatCompletionStreamNotSupported = errors.New("streaming is not supported with this method, please use CreateChatCompletionStream") //nolint:lll
	ErrModelNotAvailable                = errors.New("this model is not available for Anthropic Messages API")
)

const (
	TextContentObjectType       = "text"
	ImageContentObjectType      = "image"
	ToolUseContentObjectType    = "tool_use"
	ToolResultContentObjectType = "tool_result"
)

// ContentBlock is used to provide the [InputMessage] with multiple input or input other than a simple string
type ContentBlock struct {
	Type string `json:"type"`

	// For Text type
	Text string `json:"text,omitempty"`

	// For Image type
	Source ImageSource `json:"source,omitempty"`

	// For Tool Use type
	ID    string                 `json:"id,omitempty"`
	Name  string                 `json:"name,omitempty"`
	Input map[string]interface{} `json:"input,omitempty"`

	// For Tool Result type
	ToolUseId         string            `json:"tool_use_id,omitempty"`
	IsError           bool              `json:"is_error,omitempty"`
	ToolResultContent ToolResultContent `json:"content,omitempty"`
}

func (cb ContentBlock) MarshalJSON() ([]byte, error) {
	type alias ContentBlock
	temp := struct {
		alias
		Source            *ImageSource       `json:"source,omitempty"`
		ToolResultContent *ToolResultContent `json:"content,omitempty"`
	}{
		alias: alias(cb),
	}

	if cb.Source != (ImageSource{}) {
		temp.Source = &cb.Source
	}

	if cb.ToolResultContent != (ToolResultContent{}) {
		temp.ToolResultContent = &cb.ToolResultContent
	}

	return json.Marshal(temp)
}

type ToolResultContent struct {
	Type string `json:"type"`

	// For Text type
	Text string `json:"text,omitempty"`

	// For Image type
	Source ImageSource `json:"source,omitempty"`
}

func (trs ToolResultContent) MarshalJSON() ([]byte, error) {
	type alias ToolResultContent
	temp := struct {
		alias
		Source *ImageSource `json:"source,omitempty"`
	}{
		alias: alias(trs),
	}

	if trs.Source != (ImageSource{}) {
		temp.Source = &trs.Source
	}

	return json.Marshal(temp)
}

const ImageSourceType = "base64"
const (
	ImageJPEGMediaType = "image/jpeg"
	ImagePNGMediaType  = "image/png"
	ImageGIFMediaType  = "image/gif"
	ImageWebPMediaType = "image/webp"
)

type ImageSource struct {
	Type      string `json:"type"`
	MediaType string `json:"media_type"`
	Data      string `json:"data"`
}

// InputMessage stores content of message request. When creating new message with [Client.CreateMessage], Role field is always equal to "user".
// Content field is used to pass just one string of content. ContentBlocks are used to pass multiple input content and/or content other than a string, like an image.
//
// Note that Content and ContentBlocks fields cannot be used simultaneously in one InputMessage.
type InputMessage struct {
	Role          string `json:"role"`
	Content       string `json:"content"`
	ContentBlocks []ContentBlock
}

func (m InputMessage) MarshalJSON() ([]byte, error) {
	if m.Content != "" && m.ContentBlocks != nil {
		return nil, ErrContentFieldsMisused
	}

	if len(m.ContentBlocks) > 0 {
		msg := struct {
			Role          string         `json:"role"`
			Content       string         `json:"-"`
			ContentBlocks []ContentBlock `json:"content"`
		}(m)
		return json.Marshal(msg)
	}

	msg := struct {
		Role          string         `json:"role"`
		Content       string         `json:"content"`
		ContentBlocks []ContentBlock `json:"-"`
	}(m)
	return json.Marshal(msg)
}

func (m *InputMessage) UnmarshalJSON(bs []byte) error {
	msg := InputMessage{}

	if err := json.Unmarshal(bs, &msg); err == nil {
		*m = msg
		return nil
	}

	objectMsg := struct {
		Role          string `json:"role"`
		Content       string
		ContentBlocks []ContentBlock `json:"content"`
	}{}

	if err := json.Unmarshal(bs, &objectMsg); err != nil {
		return err
	}

	*m = InputMessage(objectMsg)
	return nil
}

const (
	Claude35SonnetModel = "claude-3-5-sonnet-20240620"
	Claude3OpusModel    = "claude-3-opus-20240229"
	Claude3SonnetModel  = "claude-3-sonnet-20240229"
	Claude3HaikuModel   = "claude-3-haiku-20240307"
)

type MessageRequestMetadata struct {
	UserId string `json:"user_id,omitempty"`
}

// MessageRequest represents a request structure for Anthropic Messages API
type MessageRequest struct {
	Model     string         `json:"model"`
	Messages  []InputMessage `json:"messages"`
	MaxTokens int            `json:"max_tokens"`

	Temperature   int                     `json:"temperature,omitempty"`
	StopSequences []string                `json:"stop_sequences,omitempty"`
	Metadata      *MessageRequestMetadata `json:"metadata,omitempty"`
	Stream        bool                    `json:"stream,omitempty"`
	System        string                  `json:"system,omitempty"`
	TopK          int                     `json:"top_k,omitempty"`
	TopP          int                     `json:"top_p,omitempty"`

	Tools      []Tool      `json:"tools,omitempty"`
	ToolChoice *ToolChoice `json:"tool_choice,omitempty"`
}

type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	InputSchema map[string]interface{} `json:"input_schema"`
}

const ObjectToolInputSchemaType = "object"

type ToolInputSchema struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

const (
	AutoToolChoiceType = "auto"
	AnyToolChoiceType  = "any"
	ToolToolChoiceType = "tool"
)

type ToolChoice struct {
	Type string `json:"type"`
	Name string `json:"name,omitempty"`
}

type StopReason string

const (
	StopReasonEndTurn      StopReason = "end_turn"
	StopReasonMaxTokens    StopReason = "max_tokens"
	StopReasonStopSequence StopReason = "stop_sequence"
	StopReasonToolUser     StopReason = "tool_use"
)

type MessageResponse struct {
	ID           string         `json:"id"`
	Type         string         `json:"type"`
	Content      []ContentBlock `json:"content"`
	Role         string         `json:"role"`
	Model        string         `json:"model"`
	StopReason   StopReason     `json:"stop_reason"`
	StopSequence string         `json:"stop_sequence,omitempty"`
	Usage        Usage          `json:"usage"`
}

type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// CreateMessage - API call to Anthropic Messages API to create a message completion
func (c *Client) CreateMessage(ctx context.Context, request MessageRequest) (response MessageResponse, err error) {
	if request.Stream {
		err = ErrChatCompletionStreamNotSupported
		return
	}

	req, err := c.newRequest(context.Background(), http.MethodPost, c.fullURL(messagesSuffix), withBody(request))
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}
