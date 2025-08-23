package connected_user

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

var ErrCreatingUserCredentials = errors.New("Failed to created user credentials")
var ErrGettingHashedPassword = errors.New("Failed to created user credentials")

var DB_PATH = "/Users/dnf/code/local-chat/data.db"

type SQLiteAuthRepo struct {
	db *sql.DB
}

var UserRepo *SQLiteAuthRepo

func NewSQLiteConnection(connectionString string) (*SQLiteAuthRepo, error) {
	db, err := sql.Open("sqlite3", DB_PATH)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	repo := &SQLiteAuthRepo{db: db}

	err = repo.Init()
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (repo *SQLiteAuthRepo) Init() error {
	_, err := repo.db.Exec(`CREATE TABLE IF NOT EXISTS credentials (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT not null unique,
	 	hashedPassword TEXT not null);`)
	if err != nil {
		return err
	}
	return nil
}

func (repo *SQLiteAuthRepo) CreateUserCredentials(username, hashedPassword string) error {
	_, err := repo.db.Exec(`insert into credential (username, hashedPassword)
values (?,?)`, username, hashedPassword)
	if err != nil {
		return ErrCreatingUserCredentials
	}
	return nil
}

func (repo *SQLiteAuthRepo) GetHashedPassword(username string) (string, error) {
	var hashedPassword string
	err := repo.db.QueryRow(`SELECT hashedPassword FROM credentials WHERE username = ?`, username).Scan(&hashedPassword)
	if err != nil {
		return "", ErrGettingHashedPassword
	}
	return hashedPassword, nil
}

func (repo *SQLiteAuthRepo) DeleteCredentials(username string) error {
	_, err := repo.db.Exec(`delete from credentials where username= ?`, username)
	if err != nil {
		return err
	}
	return nil
}
