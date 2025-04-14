package main

import (
	"fmt"
	"os"
	"strconv"
)

func printTasks(tasks []Task) {
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}
	fmt.Printf("%-5s %-30s %-15s %-20s %-20s\n", "ID", "Description", "Status", "Created At", "Updated At")
	for _, task := range tasks {
		fmt.Printf("%-5d %-30s %-15s %-20s %-20s\n",
			task.ID,
			task.Description,
			task.Status,
			task.CreatedAt.Format("2006-01-02 15:04"),
			task.UpdatedAt.Format("2006-01-02 15:04"))
	}
}

func printUsage() {
	fmt.Println("Usage: task-cli <command> [arguments]")
	fmt.Println("Commands:")
	fmt.Println("  add <description>              Add a new task")
	fmt.Println("  update <id> <description>      Update a task's description")
	fmt.Println("  delete <id>                    Delete a task")
	fmt.Println("  mark-in-progress <id>          Mark a task as in-progress")
	fmt.Println("  mark-done <id>                 Mark a task as done")
	fmt.Println("  list [todo|in-progress|done]   List tasks (optionally by status)")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	filename := "tasks.json"
	tasks, err := loadTasks(filename)
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "add":
		if len(os.Args) != 3 {
			printUsage()
			os.Exit(1)
		}
		updatedTasks, newID, err := addTask(tasks, os.Args[2])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		if err := saveTasks(filename, updatedTasks); err != nil {
			fmt.Printf("Error saving tasks: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Task added successfully (ID: %d)\n", newID)

	case "update":
		if len(os.Args) != 4 {
			printUsage()
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid ID")
			os.Exit(1)
		}
		updatedTasks, err := updateTask(tasks, id, os.Args[3])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		if err := saveTasks(filename, updatedTasks); err != nil {
			fmt.Printf("Error saving tasks: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Task updated successfully")

	case "delete":
		if len(os.Args) != 3 {
			printUsage()
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid ID")
			os.Exit(1)
		}
		updatedTasks, err := deleteTask(tasks, id)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		if err := saveTasks(filename, updatedTasks); err != nil {
			fmt.Printf("Error saving tasks: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Task deleted successfully")

	case "mark-in-progress":
		if len(os.Args) != 3 {
			printUsage()
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid ID")
			os.Exit(1)
		}
		updatedTasks, err := markTask(tasks, id, StatusInProgress)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		if err := saveTasks(filename, updatedTasks); err != nil {
			fmt.Printf("Error saving tasks: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Task marked as in-progress")

	case "mark-done":
		if len(os.Args) != 3 {
			printUsage()
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid ID")
			os.Exit(1)
		}
		updatedTasks, err := markTask(tasks, id, StatusDone)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		if err := saveTasks(filename, updatedTasks); err != nil {
			fmt.Printf("Error saving tasks: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Task marked as done")

	case "list":
		status := ""
		if len(os.Args) == 3 {
			status = os.Args[2]
			if status != StatusTodo && status != StatusInProgress && status != StatusDone {
				fmt.Println("Invalid status. Use 'todo', 'in-progress', or 'done'.")
				os.Exit(1)
			}
		} else if len(os.Args) > 3 {
			printUsage()
			os.Exit(1)
		}
		filteredTasks := filterTasks(tasks, status)
		printTasks(filteredTasks)

	default:
		printUsage()
		os.Exit(1)
	}
}
