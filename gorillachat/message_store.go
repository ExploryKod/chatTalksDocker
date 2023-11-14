package main

import (
	"database/sql"
	"errors"
)

func (t *UserStore) AddMessage(item MessageItem) (int, error) {
	res, err := t.DB.Exec("INSERT INTO messages (content, user_id, room_id, username) VALUES (?, ?, ?, ?)", item.Content, item.UserID, item.RoomID, item.Username)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (t *UserStore) GetMessagesFromRoom(id int) (MessageItem, error) {
	var messages MessageItem

	err := t.QueryRow("SELECT id, room_id, user_id, username, content FROM messages WHERE room_id = ?", id).
		Scan(&messages.ID, &messages.RoomID, &messages.UserID, &messages.Username, &messages.Content)

	if errors.Is(err, sql.ErrNoRows) {
		// User not found
		return MessageItem{}, nil
	} else if err != nil {
		// Handle other database errors
		return MessageItem{}, err
	}

	return messages, nil
}

func (t *UserStore) DeleteMessagesById(id int) error {
	_, err := t.DB.Exec("DELETE FROM messages WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func (t *UserStore) DeleteMessagesByRoomId(roomId int) error {
	_, err := t.DB.Exec("DELETE FROM messages WHERE room_id = ?", roomId)
	if err != nil {
		return err
	}

	return nil
}

func (t *UserStore) DeleteMessagesByUserId(userId int) error {
	_, err := t.DB.Exec("DELETE FROM messages WHERE user_id = ?", userId)
	if err != nil {
		return err
	}

	return nil
}
