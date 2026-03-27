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

	apiKey := os.Getenv("GROQ_API_KEY")
	baseURL := os.Getenv("GROQ_BASE_URL")

	config := openai.DefaultConfig(apiKey)
	config.BaseURL = baseURL

	client := openai.NewClientWithConfig(config)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "mixtral-8x7b-32768",
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
