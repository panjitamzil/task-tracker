package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Task represents a task in the tracker
type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Task status constants
const (
	StatusTodo       = "todo"
	StatusInProgress = "in-progress"
	StatusDone       = "done"
)

// loadTasks reads tasks from a JSON file
func loadTasks(filename string) ([]Task, error) {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var tasks []Task
	if err := json.NewDecoder(file).Decode(&tasks); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}
	return tasks, nil
}

// saveTasks writes tasks to a JSON file
func saveTasks(filename string, tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return os.WriteFile(filename, data, 0644)
}

// getNextID determines the next available ID
func getNextID(tasks []Task) int {
	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	return maxID + 1
}

// addTask adds a new task to the list
func addTask(tasks []Task, description string) ([]Task, int, error) {
	if description == "" {
		return nil, 0, fmt.Errorf("description cannot be empty")
	}
	newID := getNextID(tasks)
	now := time.Now()
	newTask := Task{
		ID:          newID,
		Description: description,
		Status:      StatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	return append(tasks, newTask), newID, nil
}

// updateTask updates a task's description
func updateTask(tasks []Task, id int, description string) ([]Task, error) {
	if description == "" {
		return nil, fmt.Errorf("description cannot be empty")
	}
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Description = description
			tasks[i].UpdatedAt = time.Now()
			return tasks, nil
		}
	}
	return nil, fmt.Errorf("task with ID %d not found", id)
}

// deleteTask removes a task by ID
func deleteTask(tasks []Task, id int) ([]Task, error) {
	for i, task := range tasks {
		if task.ID == id {
			return append(tasks[:i], tasks[i+1:]...), nil
		}
	}
	return nil, fmt.Errorf("task with ID %d not found", id)
}

// markTask changes a task's status
func markTask(tasks []Task, id int, status string) ([]Task, error) {
	if status != StatusTodo && status != StatusInProgress && status != StatusDone {
		return nil, fmt.Errorf("invalid status: %s", status)
	}
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Status = status
			tasks[i].UpdatedAt = time.Now()
			return tasks, nil
		}
	}
	return nil, fmt.Errorf("task with ID %d not found", id)
}

// filterTasks returns tasks filtered by status
func filterTasks(tasks []Task, status string) []Task {
	if status == "" {
		return tasks
	}
	var result []Task
	for _, task := range tasks {
		if task.Status == status {
			result = append(result, task)
		}
	}
	return result
}
