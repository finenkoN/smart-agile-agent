package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

func main() {

	godotenv.Load()

	apiKey := os.Getenv("GROQ_API_KEY")
	baseURL := os.Getenv("GROQ_BASE_URL")

	config := openai.DefaultConfig(apiKey)
	config.BaseURL = baseURL

	client := openai.NewClientWithConfig(config)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "openai/gpt-oss-120b",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Hello",
				},
			},
		},
	)
	fmt.Println(resp.Choices[0].Message.Content)
}
