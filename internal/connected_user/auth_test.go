package connected_user

import (
	"fmt"
	"testing"
)

func TestPasswordHandle(t *testing.T) {

	password := "my man"

	this := HashPassword(password)
	fmt.Println(this)
}

func TestDBConn(t *testing.T) {
	repo := NewSQLiteConnection(DB_PATH)
	repo.Init()

}
