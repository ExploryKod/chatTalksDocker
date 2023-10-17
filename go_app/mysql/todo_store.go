package database

import (
	"database/sql"
	"demoHTTP"
)

func NewTodoStore(db *sql.DB) *TodoStore {
	return &TodoStore{
		db,
	}
}

type TodoStore struct {
	*sql.DB
}

func (t *TodoStore) GetUsers() ([]demoHTTP.UserItem, error) {
	var users []demoHTTP.UserItem

	rows, err := t.Query("SELECT id, username, password FROM Users")
	if err != nil {
		return []demoHTTP.UserItem{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var user demoHTTP.UserItem
		if err = rows.Scan(&user.ID, &user.Username, &user.Password); err != nil {
			return []demoHTTP.UserItem{}, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return []demoHTTP.UserItem{}, err
	}

	return users, nil
}

func (t *TodoStore) GetUserByUsername(username string) (demoHTTP.UserItem, error) {
	var user demoHTTP.UserItem

	err := t.QueryRow("SELECT id, username, password FROM Users WHERE username = ?", username).
		Scan(&user.ID, &user.Username, &user.Password)

	if err == sql.ErrNoRows {
		// User not found
		return demoHTTP.UserItem{}, nil
	} else if err != nil {
		// Handle other database errors
		return demoHTTP.UserItem{}, err
	}

	return user, nil
}

func (t *TodoStore) GetTodos() ([]demoHTTP.TodoItem, error) {
	var todos []demoHTTP.TodoItem

	rows, err := t.Query("SELECT id, title, completed FROM Todos")
	if err != nil {
		return []demoHTTP.TodoItem{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var todo demoHTTP.TodoItem
		if err = rows.Scan(&todo.ID, &todo.Title, &todo.Completed); err != nil {
			return []demoHTTP.TodoItem{}, err
		}
		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		return []demoHTTP.TodoItem{}, err
	}

	return todos, nil
}

func (t *TodoStore) AddUser(item demoHTTP.UserItem) (int, error) {
	res, err := t.DB.Exec("INSERT INTO Users (username, password) VALUES (?, ?)", item.Username, item.Password)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (t *TodoStore) AddTodo(item demoHTTP.TodoItem) (int, error) {
	res, err := t.DB.Exec("INSERT INTO Todos (title, completed) VALUES (?, ?)", item.Title, item.Completed)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (t *TodoStore) DeleteTodo(id int) error {
	_, err := t.DB.Exec("DELETE FROM Todos WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func (t *TodoStore) ToggleTodo(id int) error {
	_, err := t.DB.Exec("UPDATE Todos SET `completed` = IF (`completed`, 0, 1) WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}
