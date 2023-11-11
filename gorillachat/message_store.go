package main

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
