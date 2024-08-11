package anthropic

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type streamReader struct {
	reader   *bufio.Reader
	response *http.Response
}

// Recv is the same as RecvAll() but receives only events with the type "content_block_delta", which carry the content of the response,
// and returns them as [MessageStreamDelta]
func (stream *streamReader) Recv() (MessageStreamDelta, error) {
	for {
		streamEvent, err := stream.RecvAll()
		if err != nil {
			return *new(MessageStreamDelta), err
		}

		switch streamEvent.Type {
		case ContentBlockDeltaStreamEventType:
			return streamEvent.Delta, nil
		case ErrStreamEventType:
			return *new(MessageStreamDelta), fmt.Errorf("API error of type \"%s\": %s", streamEvent.Error.Type, streamEvent.Error.Message)
		case MessageStopStreamEventType:
			return *new(MessageStreamDelta), io.EOF
		default:
			continue
		}
	}
}

// RecvAll receives all types of events from Anthropic Messages API and returns them as [MessageStreamEvent]
// If you want to process all events, check the type of event first to know what fields are available.
func (stream *streamReader) RecvAll() (response MessageStreamEvent, err error) {
	return stream.processLines()
}

func (stream *streamReader) processLines() (MessageStreamEvent, error) {
	dataPrefix := []byte("data: ")

	for {
		line, readErr := stream.reader.ReadBytes('\n')
		if readErr != nil {
			return *new(MessageStreamEvent), readErr
		}

		if string(line) == "" {
			continue
		}

		if bytes.HasPrefix(line, dataPrefix) {
			lineData := bytes.TrimPrefix(line, dataPrefix)

			var resp MessageStreamEvent
			err := json.Unmarshal(lineData, &resp)
			if err != nil {
				return *new(MessageStreamEvent), err
			}

			return resp, nil
		}
	}
}

func (stream *streamReader) Close() error {
	return stream.response.Body.Close()
}
