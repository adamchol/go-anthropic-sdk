# Go Anthropic SDK
[![Go Reference](https://pkg.go.dev/badge/github.com/adamchol/go-anthropic-sdk.svg)](https://pkg.go.dev/github.com/adamchol/go-anthropic-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/adamchol/go-anthropic-sdk)](https://goreportcard.com/report/github.com/adamchol/go-anthropic-sdk)

This library provides unofficial Go clients for [Anthropic API](https://www.anthropic.com/api).

It is heavily inspired by the [Unofficial Go SDK for OpenAI API](https://github.com/sashabaranov/go-openai) by @sashabaranov

## Installation
```sh
go get github.com/adamchol/go-anthropic-sdk 
```

## Usage
### Basic Messages API example with Claude 3.5 Sonnet
```go
package main

import (
	"context"
	"fmt"

	anthropic "github.com/adamchol/go-anthropic-sdk"
)

func main() {
	client := anthropic.NewClient("your-token")

	resp, err := client.CreateMessage(context.Background(), anthropic.MessageRequest{
		Model: anthropic.Claude35SonnetModel,
		Messages: []anthropic.InputMessage{
			{
				Role:    anthropic.MessageRoleUser,
				Content: "Hello, how are you?",
			},
		},
		MaxTokens: 4096,
	})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp.Content[0].Text)
}
```

### Other examples

<details>
<summary>Claude 3.5 Sonnet with image and text input</summary>

```go
package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/adamchol/go-anthropic-sdk"
)

func main() {
	client := anthropic.NewClient("your-token")

	imageBytes, err := os.ReadFile("ant.jpg")
	if err != nil {
		log.Fatalf("Failed to read image file: %v", err)
	}

	imgData := base64.StdEncoding.EncodeToString(imageBytes) // Encoding the image into base64

	resp, err := client.CreateMessage(context.Background(), anthropic.MessageRequest{
		Model: anthropic.Claude35SonnetModel,
		Messages: []anthropic.InputMessage{
			{
				Role: "user",
				ContentBlocks: []anthropic.ContentBlock{ // Using ContentBlocks field instead of Content for multiple input
					{
						Type: "text",
						Text: "Is there a living organism on this image?",
					},
					{
						Type: "image",
						Source: anthropic.ImageSource{
							Type:      anthropic.ImageSourceType, // "base64"
							MediaType: anthropic.ImageJPEGMediaType, // "image/jpeg"
							Data:      imgData,
						},
					},
				},
			},
		},
		MaxTokens: 1000,
	})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp.Content[0].Text)
}
```

</details>

<details>
<summary>Streaming messages</summary>

```go
package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/adamchol/go-anthropic-sdk"
)

func main() {
	client := anthropic.NewClient("your-token")

	stream, err := client.CreateMessageStream(context.Background(), anthropic.MessageRequest{
		Model:     anthropic.Claude35SonnetModel,
		MaxTokens: 1000,
		Messages: []anthropic.InputMessage{
			{
				Role:    anthropic.MessageRoleUser,
				Content: "Hello, how are you doing today?",
			},
		},
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	defer stream.Close()

	for {
		delta, err := stream.Recv()

		if err == io.EOF {
			fmt.Println()
			fmt.Println("Done")
			return
		}

		if err != nil {
			fmt.Printf("Stream error: %s", err)
			return
		}

		fmt.Print(delta.Text)
	}
}
```

</details>
