package anthropic

import (
	"context"
	"net/http"
)

type MessageStreamEventType string

const (
	MessageStartStreamEventType      MessageStreamEventType = "message_start"
	ContentBlockStartStreamEventType MessageStreamEventType = "content_block_start"
	PingStreamEventType              MessageStreamEventType = "ping"
	ContentBlockDeltaStreamEventType MessageStreamEventType = "content_block_delta"
	MessageDeltaStreamEventType      MessageStreamEventType = "message_delta"
	MessageStopStreamEventType       MessageStreamEventType = "message_stop"
	ContentBlockStopStreamEventType  MessageStreamEventType = "content_block_stop"
	ErrStreamEventType               MessageStreamEventType = "error"
)

type MessageStreamEvent struct {
	Type  MessageStreamEventType `json:"type"`
	Index int                    `json:"index,omitempty"`

	// For message_start type
	Message MessageResponse `json:"message,omitempty"`

	// For content_block_start type
	ContentBlock ContentBlock `json:"content_block,omitempty"`

	// For content_block_delta type
	Delta MessageStreamDelta `json:"delta,omitempty"`

	// For error type
	Error MessageStreamError `json:"error,omitempty"`
}

type MessageStreamDelta struct {
	Type string `json:"type"`

	Text string `json:"text,omitempty"`

	PartialJSON string `json:"partial_json,omitempty"`
}

type MessageStreamError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type MessageStream struct {
	*streamReader
}

// CreateMessageStream — API call to create a message w/ streaming
// support. It sets whether to stream back partial progress. If set, tokens will be
// sent as server-sent events as they become available, with the
// stream terminated by the last event of type "message_stop"
//
// See Recv() and RecvAll() methods of [MessageStream] for more details of how to
// receive data from stream.
func (c *Client) CreateMessageStream(ctx context.Context, request MessageRequest) (stream *MessageStream, err error) {
	request.Stream = true
	req, err := c.newRequest(context.Background(), http.MethodPost, c.fullURL(messagesSuffix), withBody(request))
	if err != nil {
		return
	}

	resp, err := c.sendStreamRequest(req)
	if err != nil {
		return
	}

	return &MessageStream{
		streamReader: resp,
	}, nil
}
