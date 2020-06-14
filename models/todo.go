package models

import (
	"fmt"
)

// Todo model
type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// default db values
var (
	todos     []*Todo
	currentID = 1
)

// GetTodos get all todos
func GetTodos() []*Todo {
	return todos
}

// AddTodo add new todo
func AddTodo(todo Todo) (Todo, error) {
	// if todo.ID != 0 {

	// }

	currentID++
	todo.ID = currentID
	todos = append(todos, &todo)
	return todo, nil
}

// GetTodoByID retrieve single todo
func GetTodoByID(id int) (Todo, error) {
	for _, t := range todos {
		if t.ID == id {
			return *t, nil
		}
	}

	return Todo{}, fmt.Errorf("Todo with ID '%v' not found", id)
}

// UpdateTodo update a todo
func UpdateTodo(todo Todo) (Todo, error) {
	for i, item := range todos {
		if item.ID == todo.ID {
			todos[i] = &todo
			return todo, nil
		}
	}

	return Todo{}, fmt.Errorf("Todo with ID '%v' not found", todo.ID)
}

// DeleteTodo delete a todo
func DeleteTodo(id int) error {
	for i, item := range todos {
		if item.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("Todo with ID '%v' not found", id)
}
