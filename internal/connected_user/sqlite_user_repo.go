package connected_user

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteUserRepo struct {
	db *sql.DB
}

func NewSQLiteUserRepo(connectionString string) (*SQLiteUserRepo, error) {
	if connectionString == "" {
		connectionString = DB_PATH
	}
	db, err := sql.Open("sqlite3", DB_PATH)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	repo := &SQLiteUserRepo{db: db}

	err = repo.Init()
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (repo *SQLiteUserRepo) Init() error {
	_, err := repo.db.Exec(`CREATE TABLE IF NOT EXISTS profiles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT not null unique);`)
	if err != nil {
		return err
	}
	return nil
}

func (repo *SQLiteUserRepo) CreateUserProfile(username string) (User, error) {
	result, err := repo.db.Exec(`insert into profiles (username) values (?)`, username)
	if err != nil {
		return User{}, err
	}

	// Get the last inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return User{}, err
	}

	// Return the created user
	return User{
		ID:       id,
		Username: username,
	}, nil
}

func (repo *SQLiteUserRepo) GetUserByID(id int64) (User, error) {
	var username string
	err := repo.db.QueryRow(`select username from profiles
		where id = ?`, id).Scan(&username)
	if err != nil {
		return User{}, err
	}
	return User{ID: int64(id), Username: username}, nil
}

func (repo *SQLiteUserRepo) GetUserByUsername(username string) (User, error) {
	var id int64
	err := repo.db.QueryRow(`select id from profiles
		where username = ?`, username).Scan(&id)

	if err != nil {
		return User{}, err
	}
	return User{ID: int64(id), Username: username}, nil
}

func (repo *SQLiteUserRepo) UpdateUserProfile(user User) error {
	// Assuming you only update the username
	_, err := repo.db.Exec(`UPDATE profiles 
		SET username = ? 
		WHERE id = ?`, user.Username, user.ID)
	return err
}

func (repo *SQLiteUserRepo) DeleteUser(id int64) error {
	_, err := repo.db.Exec(`DELETE FROM profiles 
		WHERE id = ?`, id)
	return err
}
