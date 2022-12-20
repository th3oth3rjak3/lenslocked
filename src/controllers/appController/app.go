package appController

import (
	"fmt"
	"log"
	"net/http"

	"lenslocked/config"
	"lenslocked/controllers/galleriesController"
	"lenslocked/controllers/staticController"
	"lenslocked/controllers/usersController"
	mw "lenslocked/middleware"
	"lenslocked/models/errorsModel"
	"lenslocked/models/servicesModel"
	"lenslocked/routers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type App struct {
	Config      config.AppConfig
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
	cfg := config.DefaultConfig()
	dbCfg := config.DefaultPostGresConfig()

	// Create Services
	services, err := servicesModel.NewServices(
		servicesModel.WithGorm(dbCfg.Dialect(), dbCfg.ConnectionInfo()),
		servicesModel.WithUser(config.DefaultHashKeyConfig()),
		servicesModel.WithGallery(),
		servicesModel.WithImages(),
		servicesModel.WithLogMode(cfg.IsDev()),
	)
	errorsModel.Must(err, "Could not initialize services.")

	// Run migrations
	services.AutoMigrate()

	// Destructive Reset if AutoMigrate won't work.
	// services.DestructiveReset()

	appC := NewAppController(services)
	app := &App{
		Config:      cfg,
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
		Router: echo.New(),
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
	var addr string
	if app.Config.IsProd() {
		addr = fmt.Sprintf(":%d", app.Config.Port)
	} else {
		addr = fmt.Sprintf("localhost:%d", app.Config.Port)
	}
	log.Println("Listening on:", addr)
	log.Fatal(http.ListenAndServe(addr, app.AppRouter.Router))
}

func (app *App) appMiddleware(ar *routers.AppRouter) {
	r := ar.Router
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())
	r.Use(echo.WrapMiddleware(ar.Middleware.UserMW.Invoke))
	r.Pre(middleware.RemoveTrailingSlash())
	r.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "form:csrf",
	}))
}

func (app *App) AddRoutes(ar *routers.AppRouter) {
	app.AddRoute(ar, app.appMiddleware)
	app.AddRoute(ar, app.defaultRoute)
}

func (app *App) defaultRoute(ar *routers.AppRouter) {
	r := ar.Router
	base := r.Group("/")
	base.GET("", echo.WrapHandler(app.Controllers.Static.Home))
	app.AddRoute(ar, app.contactRoutes)
	app.AddRoute(ar, app.signupRoutes)
	app.AddRoute(ar, app.loginRoutes)
	app.AddRoute(ar, app.galleriesRoutes)
	app.AddRoute(ar, app.imagesRoutes)
	app.AddRoute(ar, app.assetsRoutes)
}

func (app *App) contactRoutes(ar *routers.AppRouter) {
	r := ar.Router
	contact := r.Group("/contact")
	contact.GET("", echo.WrapHandler(app.Controllers.Static.Contact))
}

func (app *App) signupRoutes(ar *routers.AppRouter) {
	r := ar.Router
	signup := r.Group("/signup")
	signup.GET("", echo.WrapHandler(http.HandlerFunc(app.Controllers.Users.Signup)))
	signup.POST("", echo.WrapHandler(http.HandlerFunc(app.Controllers.Users.Create)))
}

func (app *App) loginRoutes(ar *routers.AppRouter) {
	r := ar.Router
	login := r.Group("/login")
	login.GET("", echo.WrapHandler(app.Controllers.Users.LoginView))
	login.POST("", echo.WrapHandler(http.HandlerFunc(app.Controllers.Users.Login)))
}

func (app *App) galleriesRoutes(ar *routers.AppRouter) {
	r := ar.Router
	requireUser := ar.Middleware.RequireUser
	galleries := r.Group("/galleries", echo.WrapMiddleware(requireUser.Invoke))
	galleries.GET("", app.Controllers.Galleries.Index)
	galleries.POST("", app.Controllers.Galleries.Create)
	galleries.GET("/new", app.Controllers.Galleries.New)
	galleries.GET("/:galleryId", app.Controllers.Galleries.Show)
	galleries.GET("/:galleryId/edit", app.Controllers.Galleries.Edit)
	galleries.POST("/:galleryId/update", app.Controllers.Galleries.Update)
	galleries.POST("/:galleryId/delete", app.Controllers.Galleries.Delete)
	galleries.POST("/:galleryId/images", app.Controllers.Galleries.ImageUpload)
	galleries.POST("/:galleryId/images/:filename/delete", app.Controllers.Galleries.ImageDelete)
}

func (app *App) imagesRoutes(ar *routers.AppRouter) {
	r := ar.Router
	r.Static("/images", "images")
}

func (app *App) assetsRoutes(ar *routers.AppRouter) {
	r := ar.Router
	r.Static("/assets", "assets")
}
