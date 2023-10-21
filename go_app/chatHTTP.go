package chatHTTP

type UserItem struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserStoreInterface interface {
	AddUser(item UserItem) (int, error)
	GetUserByUsername(username string) (UserItem, error)
	//GetUsers(username string) (UserItem, error)
}
