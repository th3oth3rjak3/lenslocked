package main

import (
	"log"
	"net/http"
	"os"

	"lenslocked/controllers"
	"lenslocked/models"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Load configuration
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Could not load the environment configuration.")
	}

	// Create Services
	us, err := models.NewUserService(os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		log.Fatalf("Could not initialize the user service: %s", err)
	}
	defer us.Close()

	// Run migrations
	us.AutoMigrate()

	// Destructive Reset if AutoMigrate won't work.
	// us.DestructiveReset()

	// Create controllers and views
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(us)

	// Create a router
	r := chi.NewRouter()
	r.Get("/", staticC.Home.ServeHTTP)
	r.Get("/contact", staticC.Contact.ServeHTTP)
	r.Get("/signup", usersC.NewView.ServeHTTP)
	r.Post("/signup", usersC.Create)
	r.Get("/login", usersC.LoginView.ServeHTTP)
	r.Post("/login", usersC.Login)

	// Start server
	addr := "localhost:3000"
	log.Println("Listening on:", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
