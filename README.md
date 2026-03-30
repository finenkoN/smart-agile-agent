# 🤖 Smart Agile AI-Agent (Go)

Микросервис на Go, который превращает обычные сообщения в Telegram в структурированные задачи в Agile-системе.

## 🚀 Основные возможности
- **AI-Управление:** Понимает команды вроде "Назначь на Никиту задачу поправить баг" через LLM Function Calling.
- **Безопасность:** Потокобезопасное хранилище (Mutex) и защита от дублей (Idempotency).
- **Интеграция:** Работает через n8n (Docker) и Telegram Bot API.

## 🛠 Стек технологий
- **Backend:** Go 1.22+ (net/http, encoding/json, sync)
- **AI:** OpenRouter / OpenAI API (Model: Gemini-2.0 / GPT-4o-mini)
- **Workflow:** n8n (Self-hosted in Docker)
- **Database:** In-memory Map (Thread-safe)

## 📐 Архитектура
1. **User** -> Telegram Bot
2. **n8n (Docker)** -> Опрашивает Telegram (Polling)
3. **Go Server** -> Принимает POST-запрос, валидирует UpdateID.
4. **AI API** -> Парсит текст в JSON-аргументы.
5. **Go Logic** -> Создает задачу в памяти и возвращает ответ.

## 📦 Запуск
1. Склонируйте репозиторий.
2. Создайте `.env` файл с `API_KEY` и `BASE_URL`.
3. Запустите Go сервер: `go run main.go`.
4. Разверните n8n в Docker и импортируйте workflow.

