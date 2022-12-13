package main

import (
	"fmt"

	"lenslocked/models"
)

const (
	host = "localhost"
	port = 5432
	user = "Jake"
	// Add password after user in the connection string if you have one.
	password = ""
	dbname   = "lenslocked_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.DestructiveReset()
	// usr := createUser(us)
	// fmt.Printf("%+v", getUserById(us, usr.ID))
}

func getUserById(us *models.UserService, id uint) *models.User {
	usr, err := us.ByID(id)
	if err != nil {
		panic(err)
	}
	return usr
}

func createUser(us *models.UserService) *models.User {
	usr := &models.User{
		Name:  "Fake User",
		Email: "fake.user@email.com",
	}
	if err := us.Create(usr); err != nil {
		panic(err)
	}
	return usr
}
