package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

func main() {

	godotenv.Load()

	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	baseURL := os.Getenv("DEEPSEEK_BASE_URL")

	config := openai.DefaultConfig(apiKey)
	config.BaseURL = baseURL

	client := openai.NewClientWithConfig(config)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "google/gemma-3-27b-it:free",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Hello",
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("Ошибка при вызове API: %v", err)
	}
	if len(resp.Choices) == 0 {
		log.Fatal("API вернул пустой список ответов")
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
