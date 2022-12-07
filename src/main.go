package main

import (
	"log"
	"net/http"

	notfound "lenslocked/routing/notfound"
	"lenslocked/views"

	"github.com/gorilla/mux"
)

// TODO: Take these out of the global namespace.
var (
	homeView    *views.View
	contactView *views.View
)

func Contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := contactView.Template.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := homeView.Template.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}

func initializeViews(){
	homeView = views.NewView("views/home.html")
	contactView = views.NewView("views/contact.html")
}

func main() {
	initializeViews()
	r := mux.NewRouter()
	r.HandleFunc("/", Home)
	r.HandleFunc("/contact", Contact)
	r.NotFoundHandler = http.HandlerFunc(notfound.NotFound)
	addr := "localhost:3000"
	log.Println("Listening on:", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
