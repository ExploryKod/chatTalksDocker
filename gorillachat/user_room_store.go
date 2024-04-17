package main

import (
	"database/sql"
	"log"
)

func (t *UserStore) GetUserDiscussions(id int) ([]DiscussionItem, error) {
	rows, err := t.Query("SELECT * FROM User_Room JOIN Users ON user_id = Users.id JOIN Rooms ON room_id = Rooms.id WHERE Users.id = ?", id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var discussions []DiscussionItem

	for rows.Next() {
		var discussion DiscussionItem
		err := rows.Scan(&discussion.ID, &discussion.UserID, &discussion.RoomID, &discussion.User.ID, &discussion.User.Username, &discussion.User.Password, &discussion.User.Admin, &discussion.User.Email, &discussion.Room.ID, &discussion.Room.Name, &discussion.Room.Description, &discussion.Room.Private)
		if err != nil {
			log.Println("ERREUR OCCURRED:", err)
			return nil, err
		}

		err = t.QueryRow("SELECT messages.content, messages.created_at, Users.username FROM messages JOIN Users ON messages.user_id = Users.id WHERE room_id = ? ORDER BY created_at DESC LIMIT 1", discussion.RoomID).Scan(&discussion.LastMessage.Content, &discussion.LastMessage.CreatedAt, &discussion.LastMessage.Username)
		if err != nil && err != sql.ErrNoRows {
			log.Println("ERROR OCCURRED:", err)
			return nil, err
		}

		discussions = append(discussions, discussion)
	}

	return discussions, nil
}

func (t *UserStore) AddUserToRoom(roomID int, userID int) error {
	_, err := t.DB.Exec("INSERT INTO User_Room (user_id, room_id) VALUES (?, ?)", userID, roomID)
	if err != nil {
		return err
	}
	return nil
}

func (t *UserStore) DeleteUserRoom(userid int, roomid int) error {
	_, err := t.DB.Exec("DELETE FROM User_Room WHERE user_id = ? AND room_id = ?", userid, roomid)
	if err != nil {
		return err
	}
	return nil
}
