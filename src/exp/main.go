package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	host = "localhost"
	port = 5432
	user = "Jake"
	// Add password after user in the connection string if you have one.
	password = ""
	dbname   = "lenslocked_dev"
)

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()
	// LogMode typically not for use in production, but good for learning and dev.
	db.LogMode(true)
	db.AutoMigrate(&User{})
	var u User
	db.Where("id < ?", 10).
		Where("id < ?", 5).
		First(&u)
	fmt.Printf("%+v\n", u)
	// say something nice...
}
