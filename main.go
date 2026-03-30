package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/finenkoN/smart-agile-agent/tasks"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

func main() {

	result1 := tasks.CreateTask("Исправить баг в авторизации", "Никита")
	fmt.Println(result1)

	result2 := tasks.CreateTask("Написать документацию", "Анна")
	fmt.Println(result2)

	result3 := tasks.CreateTask("Развернуть сервис в проде", "Никита")
	fmt.Println(result3)

	fmt.Println("\n--- Проверяем задачи для Никиты ---")
	nikitaTasks := tasks.GetTasks("Никита")
	fmt.Println(nikitaTasks)

	fmt.Println("\n--- Проверяем задачи для Петра ---")
	peterTasks := tasks.GetTasks("Петр")
	fmt.Println(peterTasks)

	godotenv.Load()

	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	baseURL := os.Getenv("DEEPSEEK_BASE_URL")

	config := openai.DefaultConfig(apiKey)
	config.BaseURL = baseURL

	client := openai.NewClientWithConfig(config)
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "Ты полезный ассистент, отвечаешь в стихах",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: "Как подговиться к собесодованию по бекеду и го",
		},
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    "openrouter/free",
			Messages: messages,
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
