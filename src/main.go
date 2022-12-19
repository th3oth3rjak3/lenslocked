package main

import (
	"log"
	"net/http"
	"os"

	"lenslocked/controllers/galleriesController"
	"lenslocked/controllers/staticController"
	"lenslocked/controllers/usersController"
	mw "lenslocked/middleware"
	"lenslocked/models/servicesModel"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	services, err := servicesModel.NewServices(os.Getenv("DB_CONNECTION_STRING"))
	must(err, "Could not initialize services.")
	defer services.Close()

	// Run migrations
	services.AutoMigrate()

	// Destructive Reset if AutoMigrate won't work.
	// services.DestructiveReset()

	// Create controllers and views
	staticC := staticController.NewStatic()
	usersC := usersController.NewUsersController(services.User)
	galleriesC := galleriesController.NewGalleriesController(services.Gallery, services.Image)

	// Create a router
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	userMw := mw.User{
		UserService:    services.User,
		GalleryService: services.Gallery,
	}
	requireUser := mw.RequireUser{
		User: userMw,
	}

	r.Use(userMw.Invoke)
	fs := http.FileServer(http.Dir("./images/galleries/"))
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
		r.Route("/galleries", func(r chi.Router) {
			// Gallery Routes

			r.Use(requireUser.Invoke)
			r.Get("/", galleriesC.Index)
			r.Post("/", galleriesC.Create)
			r.Route("/new", func(r chi.Router) {
				r.Get("/", galleriesC.New)
			})
			r.Route("/{galleryId}", func(r chi.Router) {
				r.Get("/", galleriesC.Show)
				r.Get("/edit", galleriesC.Edit)
				r.Post("/update", galleriesC.Update)
				r.Post("/delete", galleriesC.Delete)
				r.Route("/images", func(r chi.Router) {
					r.Post("/", galleriesC.ImageUpload)
					r.Route("/{filename}", func(r chi.Router) {
						r.Post("/delete", galleriesC.ImageDelete)
					})
				})
			})
		})
		r.Route("/images/galleries/{galleryID}", func(r chi.Router) {
			r.Handle("/*", http.StripPrefix("/images/galleries/", requireUser.ImageSafety(fs)))
		})
	})

	// Start server
	addr := "localhost:3000"
	log.Println("Listening on:", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
