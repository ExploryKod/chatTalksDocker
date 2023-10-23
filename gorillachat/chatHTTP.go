package main

type UserItem struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RoomItem struct {
	ID          int                `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Clients     map[string]*Client `json:"-"`
}

type UserStoreInterface interface {
	AddUser(item UserItem) (int, error)
	GetUserByUsername(username string) (UserItem, error)
	GetUsers() ([]UserItem, error)
	AddRoom(item RoomItem) (int, error)
	GetRoomByName(name string) (RoomItem, error)
	GetRoomById(id int) (RoomItem, error)
	AddUserToRoom(roomID int, userID int) error
	GetUsersFromRoom(roomID int) ([]UserItem, error)
	GetOneUserFromRoom(roomID int, userID int) (UserItem, error)
	GetRooms() ([]RoomItem, error)
}
