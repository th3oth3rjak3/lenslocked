package main

import (
	"log"
	"net/http"
	"os"

	"lenslocked/controllers/galleries"
	"lenslocked/controllers/static"
	"lenslocked/controllers/users"
	"lenslocked/middleware"
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
	defer services.Close()

	// Run migrations
	services.AutoMigrate()

	// Destructive Reset if AutoMigrate won't work.
	// services.DestructiveReset()

	// Create controllers and views
	staticC := static.NewStatic()
	usersC := users.NewUsersController(services.User)
	galleriesC := galleries.NewGalleriesController(services.Gallery)

	// Create a router
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/", staticC.Home.ServeHTTP)
		r.Route("/contact", func(r chi.Router) {
			r.Get("/", staticC.Contact.ServeHTTP)
		})
		r.Route("/signup", func(r chi.Router) {
			r.Get("/", usersC.Signup)
			r.Post("/", usersC.Create)
		})
		r.Route("/login", func(r chi.Router) {
			r.Get("/", usersC.LoginView.ServeHTTP)
			r.Post("/", usersC.Login)
		})
		r.Route("/cookietest", func(r chi.Router) {
			r.Get("/", usersC.CookieTest)
		})
		r.Route("/galleries", func(r chi.Router) {
			// Gallery Routes
			requireUser := middleware.RequireUser{
				UserService: services.User,
			}
			r.Use(requireUser.Invoke)
			r.Get("/new", galleriesC.New)
			r.Post("/", galleriesC.Create)
		})
	})

	// Start server
	addr := "localhost:3000"
	log.Println("Listening on:", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
