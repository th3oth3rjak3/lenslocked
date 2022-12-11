package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "Jake"
	password = "" // Add this after user in the connection string if you have one.
	dbname   = "lenslocked_dev"
)

func main(){
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()
	rows, err := db.Query(`
		SELECT users.id AS user_id, users.name, users.email,
		orders.id AS order_id, orders.description, orders.amount
		FROM users
		INNER JOIN orders ON users.id = orders.user_id;
	`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var (
		user_id int
		name string
		email string
		order_id int
		description string
		amount int
	)
	for rows.Next() {
		err := rows.Scan(&user_id, &name, &email, &order_id, &description, &amount)
		if err != nil {
			panic(err)
		}
		fmt.Println(user_id, name, email, order_id, description, amount)
	}
}