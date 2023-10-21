package database

import (
	"chatHTTP"
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

//func (t *UserStore) GetUsers() ([]chatHTTP.UserItem, error) {
//	var users []chatHTTP.UserItem
//
//	rows, err := t.Query("SELECT id, username, password FROM Users")
//	if err != nil {
//		return []chatHTTP.UserItem{}, err
//	}
//
//	defer rows.Close()
//
//	for rows.Next() {
//		var user chatHTTP.UserItem
//		if err = rows.Scan(&user.ID, &user.Username, &user.Password); err != nil {
//			return []chatHTTP.UserItem{}, err
//		}
//		users = append(users, user)
//	}
//
//	if err = rows.Err(); err != nil {
//		return []chatHTTP.UserItem{}, err
//	}
//
//	return users, nil
//}

func (t *UserStore) GetUserByUsername(username string) (chatHTTP.UserItem, error) {
	var user chatHTTP.UserItem

	err := t.QueryRow("SELECT id, username, password FROM Users WHERE username = ?", username).
		Scan(&user.ID, &user.Username, &user.Password)

	if err == sql.ErrNoRows {
		// User not found
		return chatHTTP.UserItem{}, nil
	} else if err != nil {
		// Handle other database errors
		return chatHTTP.UserItem{}, err
	}

	return user, nil
}

func (t *UserStore) AddUser(item chatHTTP.UserItem) (int, error) {
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
