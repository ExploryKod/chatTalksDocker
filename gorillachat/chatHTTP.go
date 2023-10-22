package main

type UserItem struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RoomItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserStoreInterface interface {
	AddUser(item UserItem) (int, error)
	GetUserByUsername(username string) (UserItem, error)
	//GetUsers(username string) (UserItem, error)
}
type RoomStoreInterface interface {
	AddRoom(item RoomItem) (int, error)
	GetRoomByName(name string) (RoomItem, error)
	GetRoomById(id int) (RoomItem, error)
}
