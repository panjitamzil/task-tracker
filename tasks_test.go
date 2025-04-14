package main

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestLoadTasks_FileDoesNotExist(t *testing.T) {
	tasks, err := loadTasks("nonexistent.json")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("Expected empty tasks, got %v", tasks)
	}
}

func TestLoadTasks_ValidJSON(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "tasks*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	sampleTasks := []Task{
		{ID: 1, Description: "Test task", Status: StatusTodo, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	data, _ := json.Marshal(sampleTasks)
	tmpfile.Write(data)
	tmpfile.Close()

	tasks, err := loadTasks(tmpfile.Name())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(tasks) != 1 || tasks[0].Description != "Test task" {
		t.Errorf("Expected one task with description 'Test task', got %v", tasks)
	}
}

func TestSaveTasks(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "tasks*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	tasks := []Task{{ID: 1, Description: "Save test", Status: StatusTodo}}
	if err := saveTasks(tmpfile.Name(), tasks); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	data, _ := os.ReadFile(tmpfile.Name())
	var saved []Task
	json.Unmarshal(data, &saved)
	if len(saved) != 1 || saved[0].Description != "Save test" {
		t.Errorf("Expected saved task with description 'Save test', got %v", saved)
	}
}

func TestAddTask(t *testing.T) {
	tasks := []Task{}
	updatedTasks, newID, err := addTask(tasks, "New task")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if newID != 1 || len(updatedTasks) != 1 || updatedTasks[0].Status != StatusTodo {
		t.Errorf("Expected ID 1 and one todo task, got ID %d, tasks %v", newID, updatedTasks)
	}
}

func TestAddTask_EmptyDescription(t *testing.T) {
	tasks := []Task{}
	_, _, err := addTask(tasks, "")
	if err == nil {
		t.Errorf("Expected error for empty description, got none")
	}
}

func TestUpdateTask(t *testing.T) {
	tasks := []Task{{ID: 1, Description: "Old task", Status: StatusTodo}}
	updatedTasks, err := updateTask(tasks, 1, "Updated task")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if updatedTasks[0].Description != "Updated task" {
		t.Errorf("Expected description 'Updated task', got %s", updatedTasks[0].Description)
	}
}

func TestUpdateTask_NotFound(t *testing.T) {
	tasks := []Task{{ID: 1, Description: "Old task", Status: StatusTodo}}
	_, err := updateTask(tasks, 2, "Updated task")
	if err == nil {
		t.Errorf("Expected error for nonexistent ID, got none")
	}
}

func TestDeleteTask(t *testing.T) {
	tasks := []Task{{ID: 1, Description: "Task", Status: StatusTodo}}
	updatedTasks, err := deleteTask(tasks, 1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(updatedTasks) != 0 {
		t.Errorf("Expected empty tasks, got %v", updatedTasks)
	}
}

func TestMarkTask(t *testing.T) {
	tasks := []Task{{ID: 1, Description: "Task", Status: StatusTodo}}
	updatedTasks, err := markTask(tasks, 1, StatusDone)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if updatedTasks[0].Status != StatusDone {
		t.Errorf("Expected status 'done', got %s", updatedTasks[0].Status)
	}
}

func TestFilterTasks(t *testing.T) {
	tasks := []Task{
		{ID: 1, Description: "Task 1", Status: StatusTodo},
		{ID: 2, Description: "Task 2", Status: StatusDone},
	}
	filtered := filterTasks(tasks, StatusTodo)
	if len(filtered) != 1 || filtered[0].ID != 1 {
		t.Errorf("Expected one todo task, got %v", filtered)
	}
}
