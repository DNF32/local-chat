package connected_user

type UserRepository interface {
	CreateUserProfile(username string) (User, error)
	GetUserByID(id int64) (User, error)
	GetUserByUsername(username string) (User, error)
	UpdateUserProfile(user User) error
	DeleteUser(id int64) error
}
