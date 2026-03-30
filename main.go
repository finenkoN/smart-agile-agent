package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/finenkoN/smart-agile-agent/tasks"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

var aiClient *openai.Client

// Структуры для приема и отправки JSON по HTTP
type ChatRequest struct {
	Message  string `json:"message"`
	UpdateID int    `json:"update_id"`
}

type ChatResponse struct {
	Reply string `json:"reply"`
}

var lastUpdateID int

func HandleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только POST запросы", http.StatusMethodNotAllowed)
		return
	}

	// Читаем входящий JSON
	var reqBody ChatRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}
	if reqBody.UpdateID <= lastUpdateID {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ChatResponse{Reply: "SKIP"}) // Говорим n8n, что это дубликат
		return
	}
	// Запоминаем новый ID
	lastUpdateID = reqBody.UpdateID

	fmt.Printf("\n[СЕРВЕР ПОЛУЧИЛ ЗАПРОС]: %s\n", reqBody.Message)

	var createTaskTool = openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "create_task",
			Description: "Создает новую задачу в системе и назначает ее на сотрудника.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"title":    map[string]any{"type": "string"},
					"assignee": map[string]any{"type": "string"},
				},
				"required": []string{"title", "assignee"},
			},
		},
	}

	aiReq := openai.ChatCompletionRequest{
		Model: "openrouter/free",
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: reqBody.Message},
		},
		Tools: []openai.Tool{createTaskTool},
	}

	// Делаем запрос к нейросети
	ctx := context.Background()
	resp, err := aiClient.CreateChatCompletion(ctx, aiReq)
	if err != nil {
		http.Error(w, "Ошибка AI", http.StatusInternalServerError)
		return
	}
	if len(resp.Choices) == 0 {
		log.Println("AI вернул пустой ответ (Choices)")
		json.NewEncoder(w).Encode(ChatResponse{Reply: "Извини, нейросеть промолчала."})
		return
	}
	msg := resp.Choices[0].Message
	finalReply := ""

	// Если AI вызвал функцию
	if len(msg.ToolCalls) > 0 {
		toolCall := msg.ToolCalls[0]
		if toolCall.Function.Name == "create_task" {
			var args struct {
				Title    string `json:"title"`
				Assignee string `json:"assignee"`
			}
			json.Unmarshal([]byte(toolCall.Function.Arguments), &args)

			// Вызываем Go-функцию
			finalReply = tasks.CreateTask(args.Title, args.Assignee)
		}
	} else {
		finalReply = msg.Content
	}

	// Формируем ответ клиенту
	w.Header().Set("Content-Type", "application/json")
	response := ChatResponse{Reply: finalReply}

	json.NewEncoder(w).Encode(response)
	fmt.Printf("[СЕРВЕР ОТВЕТИЛ]: %s\n", finalReply)
}

func main() {
	godotenv.Load()
	apiKey := os.Getenv("API_KEY")
	baseURL := os.Getenv("BASE_URL")

	config := openai.DefaultConfig(apiKey)
	config.BaseURL = baseURL
	aiClient = openai.NewClientWithConfig(config)

	http.HandleFunc("/chat", HandleChat)

	port := ":8080"
	fmt.Println("Сервер запущен на http://localhost" + port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
