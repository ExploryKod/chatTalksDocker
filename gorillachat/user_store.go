package main

import (
	"database/sql"
)

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db,
	}
}

type UserStore struct {
	*sql.DB
}

func (t *UserStore) GetUsers() ([]UserItem, error) {
	var users []UserItem

	rows, err := t.Query("SELECT id, username, password, admin FROM Users")
	if err != nil {
		return []UserItem{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var user UserItem
		if err = rows.Scan(&user.ID, &user.Username, &user.Password, &user.Admin); err != nil {
			return []UserItem{}, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return []UserItem{}, err
	}

	return users, nil
}

func (t *UserStore) GetUserByUsername(username string) (UserItem, error) {
	var user UserItem

	err := t.QueryRow("SELECT id, username, password FROM Users WHERE username = ?", username).
		Scan(&user.ID, &user.Username, &user.Password)

	if err == sql.ErrNoRows {
		// User not found
		return UserItem{}, nil
	} else if err != nil {
		// Handle other database errors
		return UserItem{}, err
	}

	return user, nil
}

func (t *UserStore) AddUser(item UserItem) (int, error) {
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

func (t *UserStore) DeleteUserById(id int) error {
	_, err := t.DB.Exec("DELETE FROM Users WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}
