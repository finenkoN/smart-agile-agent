package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

type uaTransport struct{}

func (t *uaTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	return http.DefaultTransport.RoundTrip(req)
}

func main() {

	godotenv.Load()

	apiKey := os.Getenv("GROQ_API_KEY")
	baseURL := os.Getenv("GROQ_BASE_URL")

	config := openai.DefaultConfig(apiKey)
	config.BaseURL = baseURL

	config.HTTPClient = &http.Client{
		Transport: &http.Transport{
			// можно оставить пустым, но User-Agent зададим отдельно
		},
	}
	// Задаём User-Agent для всех запросов
	config.HTTPClient = &http.Client{
		Transport: &uaTransport{},
	}

	client := openai.NewClientWithConfig(config)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "llama-3.1-8b-instant",
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
