# Task Tracker CLI

A command-line application to manage tasks, built with Go. Tasks are stored in a JSON file (`tasks.json`) and can be added, updated, deleted, marked as in-progress or done, and listed by status.

## Features
- Add a new task with a description
- Update a task's description
- Delete a task by ID
- Mark a task as in-progress or done
- List all tasks or filter by status (todo, in-progress, done)
- Persistent storage in a JSON file
- Error handling for invalid inputs

## Requirements
- Go 1.21 or later

## Installation
1. Clone or download this repository.
2. Navigate to the project directory:
```
cd task-tracker
```
4. Initialize Go modules if not already done
```
go mod init task-tracker
```
5. Build the application
```
go build -o task-tracker
```

## Usage

Run the application using the compiled binary (e.g., `./task-tracker`).

### Commands

| Command                           | Description                              | Example                                    |
|-----------------------------------|------------------------------------------|--------------------------------------------|
| `add <description>`               | Add a new task                           | `./task-tracker add "Buy groceries"`       |
| `update <id> <description>`       | Update a task's description              | `./task-tracker update 1 "Buy milk"`       |
| `delete <id>`                     | Delete a task                            | `./task-tracker delete 1`                  |
| `mark-in-progress <id>`           | Mark a task as in-progress               | `./task-tracker mark-in-progress 1`        |
| `mark-done <id>`                  | Mark a task as done                      | `./task-tracker mark-done 1`               |
| `list [todo\|in-progress\|done]`  | List tasks (optionally by status)        | `./task-tracker list done`                 |

**Note**: If no status is provided with `list`, all tasks are shown.

