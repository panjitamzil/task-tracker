package main

import (
    "encoding/json"
    "os"
    "reflect"
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

func TestLoadTasks_EmptyFile(t *testing.T) {
    tmpfile, err := os.CreateTemp("", "tasks*.json")
    if err != nil {
        t.Fatal(err)
    }
    defer os.Remove(tmpfile.Name())

    tmpfile.Close()

    tasks, err := loadTasks(tmpfile.Name())
    if err == nil {
        t.Errorf("Expected error for empty file, got nil")
    }
    if tasks != nil {
        t.Errorf("Expected nil tasks, got %v", tasks)
    }
}

func TestLoadTasks_InvalidJSON(t *testing.T) {
    tmpfile, err := os.CreateTemp("", "tasks*.json")
    if err != nil {
        t.Fatal(err)
    }
    defer os.Remove(tmpfile.Name())

    tmpfile.Write([]byte("{invalid json"))
    tmpfile.Close()

    tasks, err := loadTasks(tmpfile.Name())
    if err == nil {
        t.Errorf("Expected error for invalid JSON, got nil")
    }
    if tasks != nil {
        t.Errorf("Expected nil tasks, got %v", tasks)
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

func TestSaveTasks_EmptyTasks(t *testing.T) {
    tmpfile, err := os.CreateTemp("", "tasks*.json")
    if err != nil {
        t.Fatal(err)
    }
    defer os.Remove(tmpfile.Name())

    tasks := []Task{}
    if err := saveTasks(tmpfile.Name(), tasks); err != nil {
        t.Errorf("Expected no error, got %v", err)
    }

    data, _ := os.ReadFile(tmpfile.Name())
    var saved []Task
    json.Unmarshal(data, &saved)
    if len(saved) != 0 {
        t.Errorf("Expected empty tasks, got %v", saved)
    }
}

func TestGetNextID_EmptyTasks(t *testing.T) {
    tasks := []Task{}
    id := getNextID(tasks)
    if id != 1 {
        t.Errorf("Expected ID 1 for empty tasks, got %d", id)
    }
}

func TestGetNextID_NonEmptyTasks(t *testing.T) {
    tasks := []Task{
        {ID: 1, Description: "Task 1"},
        {ID: 3, Description: "Task 3"},
        {ID: 2, Description: "Task 2"},
    }
    id := getNextID(tasks)
    if id != 4 {
        t.Errorf("Expected ID 4 for tasks with max ID 3, got %d", id)
    }
}

func TestAddTask(t *testing.T) {
    tasks := []Task{}
    updatedTasks, newID, err := addTask(tasks, "New task")
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    if newID != 1 {
        t.Errorf("Expected ID 1, got %d", newID)
    }
    if len(updatedTasks) != 1 || updatedTasks[0].Description != "New task" || updatedTasks[0].Status != StatusTodo {
        t.Errorf("Expected one task with description 'New task' and status 'todo', got %v", updatedTasks)
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
    if updatedTasks[0].UpdatedAt.Before(updatedTasks[0].CreatedAt) {
        t.Errorf("Expected UpdatedAt to be updated, got %v", updatedTasks[0].UpdatedAt)
    }
}

func TestUpdateTask_NotFound(t *testing.T) {
    tasks := []Task{{ID: 1, Description: "Old task", Status: StatusTodo}}
    _, err := updateTask(tasks, 2, "Updated task")
    if err == nil {
        t.Errorf("Expected error for nonexistent ID, got none")
    }
}

func TestUpdateTask_EmptyDescription(t *testing.T) {
    tasks := []Task{{ID: 1, Description: "Old task", Status: StatusTodo}}
    _, err := updateTask(tasks, 1, "")
    if err == nil {
        t.Errorf("Expected error for empty description, got none")
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

func TestDeleteTask_NotFound(t *testing.T) {
    tasks := []Task{{ID: 1, Description: "Task", Status: StatusTodo}}
    _, err := deleteTask(tasks, 2)
    if err == nil {
        t.Errorf("Expected error for nonexistent ID, got none")
    }
}

func TestDeleteTask_EmptyTasks(t *testing.T) {
    tasks := []Task{}
    _, err := deleteTask(tasks, 1)
    if err == nil {
        t.Errorf("Expected error for empty tasks, got none")
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
    if updatedTasks[0].UpdatedAt.Before(updatedTasks[0].CreatedAt) {
        t.Errorf("Expected UpdatedAt to be updated, got %v", updatedTasks[0].UpdatedAt)
    }
}

func TestMarkTask_NotFound(t *testing.T) {
    tasks := []Task{{ID: 1, Description: "Task", Status: StatusTodo}}
    _, err := markTask(tasks, 2, StatusDone)
    if err == nil {
        t.Errorf("Expected error for nonexistent ID, got none")
    }
}

func TestMarkTask_InvalidStatus(t *testing.T) {
    tasks := []Task{{ID: 1, Description: "Task", Status: StatusTodo}}
    _, err := markTask(tasks, 1, "invalid")
    if err == nil {
        t.Errorf("Expected error for invalid status, got none")
    }
}

func TestFilterTasks_AllTasks(t *testing.T) {
    tasks := []Task{
        {ID: 1, Description: "Task 1", Status: StatusTodo},
        {ID: 2, Description: "Task 2", Status: StatusDone},
    }
    filtered := filterTasks(tasks, "")
    if !reflect.DeepEqual(filtered, tasks) {
        t.Errorf("Expected all tasks, got %v", filtered)
    }
}

func TestFilterTasks_ByStatus(t *testing.T) {
    tasks := []Task{
        {ID: 1, Description: "Task 1", Status: StatusTodo},
        {ID: 2, Description: "Task 2", Status: StatusDone},
        {ID: 3, Description: "Task 3", Status: StatusTodo},
    }
    filtered := filterTasks(tasks, StatusTodo)
    if len(filtered) != 2 || filtered[0].ID != 1 || filtered[1].ID != 3 {
        t.Errorf("Expected two todo tasks (IDs 1, 3), got %v", filtered)
    }
}

func TestFilterTasks_NoMatch(t *testing.T) {
    tasks := []Task{
        {ID: 1, Description: "Task 1", Status: StatusTodo},
    }
    filtered := filterTasks(tasks, StatusDone)
    if len(filtered) != 0 {
        t.Errorf("Expected no tasks, got %v", filtered)
    }
}