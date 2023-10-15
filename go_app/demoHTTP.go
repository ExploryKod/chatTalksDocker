package demoHTTP

import "embed"

//go:embed templates/*
var EmbedTemplates embed.FS

type TodoItem struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type UserItem struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type TodoStoreInterface interface {
	GetTodos() ([]TodoItem, error)
	AddTodo(item TodoItem) (int, error)
	AddUser(item UserItem) (int, error)
	GetUserByUsername(username string) (UserItem, error)
	DeleteTodo(id int) error
	ToggleTodo(id int) error
}

// type UserStoreInterface interface {
// 	Getuser() ([]UserItem, error)

// }
