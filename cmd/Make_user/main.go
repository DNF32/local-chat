package main

import (
	"fmt"
	"local-chat/internal/connected_user"
	"strings"
)

func main() {
	fmt.Print("Enter username: ")
	var username string
	fmt.Scanln(&username)
	username = strings.TrimSpace(username)

	// Password (visible)
	fmt.Print("Enter password: ")
	var password string
	fmt.Scanln(&password)
	password = strings.TrimSpace(password)
	hashPassword := connected_user.HashPassword(password)

	authRepo, _ := connected_user.NewSQLiteAuthRepo("")
	userRepo, _ := connected_user.NewSQLiteUserRepo("")

	err := authRepo.CreateUserCredentials(username, hashPassword)
	if err!=nil{
		panic(err)
	}
	_, _ = userRepo.CreateUserProfile(username)
}
