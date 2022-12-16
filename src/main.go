package main

import (
	"log"
	"net/http"
	"os"

	"lenslocked/controllers/static"
	"lenslocked/controllers/users"
	"lenslocked/models"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func must(err error, failureMessage string) error {
	if err != nil {
		log.Fatalf(failureMessage+": %s", err)
		return err
	}
	return nil
}

func main() {
	// Load configuration
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Could not load the environment configuration.")
	}

	// Create Services
	services, err := models.NewServices(os.Getenv("DB_CONNECTION_STRING"))
	must(err, "Could not initialize services.")
	// TODO: fix the services.Close and Automigrate/DestructiveReset functions.
	// defer services.Close()

	// Run migrations
	// services.AutoMigrate()

	// Destructive Reset if AutoMigrate won't work.
	// us.DestructiveReset()

	// Create controllers and views
	staticC := static.NewStatic()
	usersC := users.NewUsersController(services.User)

	// Create a router
	r := chi.NewRouter()
	r.Get("/", staticC.Home.ServeHTTP)
	r.Get("/contact", staticC.Contact.ServeHTTP)
	r.Get("/signup", usersC.Signup)
	r.Post("/signup", usersC.Create)
	r.Get("/login", usersC.LoginView.ServeHTTP)
	r.Post("/login", usersC.Login)
	r.Get("/cookietest", usersC.CookieTest)

	// Start server
	addr := "localhost:3000"
	log.Println("Listening on:", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
