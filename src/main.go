package main

import (
	"log"
	"net/http"

	"lenslocked/controllers"
	notfound "lenslocked/routing/notfound"
	"lenslocked/views"

	"github.com/gorilla/mux"
)

// TODO: Take these out of the global namespace.
var (
	homeView    *views.View
	contactView *views.View
)

// Renders the contact page.
func Contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(contactView.Render(w, nil))
}

// Renders the homepage.
func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(homeView.Render(w, nil))
}

// This function is used to define things that must work or else panic.
func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// Create controllers and views
	homeView = views.NewView("bootstrap", "views/home.html")
	contactView = views.NewView("bootstrap", "views/contact.html")
	usersC := controllers.NewUsers()

	// Create a router
	r := mux.NewRouter()
	r.HandleFunc("/", Home)
	r.HandleFunc("/contact", Contact)
	r.HandleFunc("/signup", usersC.New)
	r.NotFoundHandler = http.HandlerFunc(notfound.NotFound)

	// Start server
	addr := "localhost:3000"
	log.Println("Listening on:", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
