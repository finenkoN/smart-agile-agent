package tasks

import (
	"fmt"
	"strings"
	"sync"
)

type Task struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Assignee string `json:"assignee"`
	Status   string `json:"status"` // "todo", "in_progress", "done"
}

var (
	taskDB    = make(map[int]Task)
	currentID = 1
	mu        sync.Mutex
)

func CreateTask(title, assignee string) string {
	mu.Lock()
	defer mu.Unlock()

	task := Task{
		ID:       currentID,
		Title:    title,
		Assignee: assignee,
		Status:   "todo",
	}
	taskDB[task.ID] = task
	currentID++

	return fmt.Sprintf("Задача #%d '%s' создана и назначена на %s.", task.ID, task.Title, task.Assignee)
}

func GetTasks(assignee string) string {
	mu.Lock()
	defer mu.Unlock()

	var userTasks []string
	for _, task := range taskDB {
		if task.Assignee == assignee {
			userTasks = append(userTasks, fmt.Sprintf("#%d: %s (%s)", task.ID, task.Title, task.Status))
		}
	}
	if len(userTasks) == 0 {
		return fmt.Sprintf("На %s нет назначенных задач.", assignee)
	}
	return fmt.Sprintf("Задачи для %s:\n%s", assignee, strings.Join(userTasks, "\n"))
}
