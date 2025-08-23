package connected_user

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteUserRepo struct {
	db *sql.DB
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
