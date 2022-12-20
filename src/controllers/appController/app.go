package appController

import (
	"log"
	"net/http"
	"os"

	"lenslocked/controllers/galleriesController"
	"lenslocked/controllers/staticController"
	"lenslocked/controllers/usersController"
	mw "lenslocked/middleware"
	"lenslocked/models/errorsModel"
	"lenslocked/models/servicesModel"
	"lenslocked/routers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

type App struct {
	Controllers *AppController
	Services    *servicesModel.Services
	AppRouter   *routers.AppRouter
	ImageServer http.Handler
	AssetServer http.Handler
}

type AppController struct {
	Static    *staticController.Static
	Users     *usersController.UsersController
	Galleries *galleriesController.GalleriesController
}

func NewApp() *App {
	// Load configuration
	errorsModel.Must(godotenv.Load(".env"), "Could not load the environment configuration.")

	// Create Services
	services, err := servicesModel.NewServices(os.Getenv("DB_CONNECTION_STRING"))
	errorsModel.Must(err, "Could not initialize services.")

	// Run migrations
	services.AutoMigrate()

	// Destructive Reset if AutoMigrate won't work.
	// services.DestructiveReset()

	appC := NewAppController(services)
	app := &App{
		Services:    services,
		Controllers: appC,
		ImageServer: http.FileServer(http.Dir("./images/galleries/")),
		AssetServer: http.FileServer(http.Dir("./assets/")),
	}
	app.AppRouter = app.NewAppRouter()
	app.AddRoutes(app.AppRouter)
	return app
}

func NewAppController(s *servicesModel.Services) *AppController {
	staticC := staticController.NewStatic()
	usersC := usersController.NewUsersController(s.User)
	galleriesC := galleriesController.NewGalleriesController(s.Gallery, s.Image)
	return &AppController{
		Static:    staticC,
		Users:     usersC,
		Galleries: galleriesC,
	}
}

func (app *App) NewAppRouter() *routers.AppRouter {
	userMw := mw.User{
		UserService:    app.Services.User,
		GalleryService: app.Services.Gallery,
	}
	requireUser := mw.RequireUser{
		User: userMw,
	}

	appRouter := routers.AppRouter{
		Router: chi.NewRouter(),
		Middleware: routers.AppMiddleware{
			UserMW:      &userMw,
			RequireUser: &requireUser,
		},
	}
	return &appRouter
}

func (app *App) AddRoute(router *routers.AppRouter, fn routers.RouteFunction) {
	fn(router)
}

func (app *App) Run() {
	defer app.Services.Close()
	addr := "localhost:3000"
	log.Println("Listening on:", addr)
	log.Fatal(http.ListenAndServe(addr, app.AppRouter.Router))
}

func (app *App) appMiddleware(ar *routers.AppRouter) {
	ar.Router = chi.NewRouter()
	r := ar.Router
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(ar.Middleware.UserMW.Invoke)
}

func (app *App) AddRoutes(ar *routers.AppRouter) {
	app.AddRoute(ar, app.appMiddleware)
	app.AddRoute(ar, app.defaultRoute)
}

func (app *App) defaultRoute(ar *routers.AppRouter) {
	r := ar.Router
	r.Route("/", func(r chi.Router) {
		r.Get("/", app.Controllers.Static.Home.ServeHTTP)
		app.AddRoute(ar, app.contactRoutes)
		app.AddRoute(ar, app.signupRoutes)
		app.AddRoute(ar, app.loginRoutes)
		app.AddRoute(ar, app.galleriesRoutes)
		app.AddRoute(ar, app.imagesRoutes)
		app.AddRoute(ar, app.assetsRoutes)
	})
}

func (app *App) contactRoutes(ar *routers.AppRouter) {
	r := ar.Router
	r.Route("/contact", func(r chi.Router) {
		r.Get("/", app.Controllers.Static.Contact.ServeHTTP)
	})
}

func (app *App) signupRoutes(ar *routers.AppRouter) {
	r := ar.Router
	r.Route("/signup", func(r chi.Router) {
		r.Get("/", app.Controllers.Users.Signup)
		r.Post("/", app.Controllers.Users.Create)
	})
}

func (app *App) loginRoutes(ar *routers.AppRouter) {
	r := ar.Router
	r.Route("/login", func(r chi.Router) {
		r.Get("/", app.Controllers.Users.LoginView.ServeHTTP)
		r.Post("/", app.Controllers.Users.Login)
	})
}

func (app *App) galleriesRoutes(ar *routers.AppRouter) {
	r := ar.Router
	requireUser := ar.Middleware.RequireUser
	r.Route("/galleries", func(r chi.Router) {
		// Gallery Routes

		r.Use(requireUser.Invoke)
		r.Get("/", app.Controllers.Galleries.Index)
		r.Post("/", app.Controllers.Galleries.Create)
		r.Route("/new", func(r chi.Router) {
			r.Get("/", app.Controllers.Galleries.New)
		})
		r.Route("/{galleryId}", func(r chi.Router) {
			r.Get("/", app.Controllers.Galleries.Show)
			r.Get("/edit", app.Controllers.Galleries.Edit)
			r.Post("/update", app.Controllers.Galleries.Update)
			r.Post("/delete", app.Controllers.Galleries.Delete)
			r.Route("/images", func(r chi.Router) {
				r.Post("/", app.Controllers.Galleries.ImageUpload)
				r.Route("/{filename}", func(r chi.Router) {
					r.Post("/delete", app.Controllers.Galleries.ImageDelete)
				})
			})
		})
	})
}

func (app *App) imagesRoutes(ar *routers.AppRouter) {
	r := ar.Router
	r.Route("/images/galleries/{galleryID}", func(r chi.Router) {
		r.Handle("/*", http.StripPrefix("/images/galleries/", ar.Middleware.RequireUser.ImageSafety(app.ImageServer)))
	})
}

func (app *App) assetsRoutes(ar *routers.AppRouter) {
	r := ar.Router
	r.Route("/assets", func(r chi.Router) {
		r.Handle("/*", http.StripPrefix("/assets/", app.AssetServer))
	})
}
