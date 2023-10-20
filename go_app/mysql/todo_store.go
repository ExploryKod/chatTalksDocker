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
