package main

import (
	"log"
	"net/http"
	"lenslocked/controllers"
	"github.com/go-chi/chi/v5"
)

func main() {
	// Create controllers and views
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers()

	// Create a router
	r := chi.NewRouter()
	r.Get("/", staticC.Home.ServeHTTP)
	r.Get("/contact", staticC.Contact.ServeHTTP)
	r.Get("/signup", usersC.New)
	r.Post("/signup", usersC.Create)

	// Start server
	addr := "localhost:3000"
	log.Println("Listening on:", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
