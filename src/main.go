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
	must(contactView.Render(w, nil))

}

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(homeView.Render(w, nil))
}

func initializeViews(){
	homeView = views.NewView("bootstrap", "views/home.html")
	contactView = views.NewView("bootstrap", "views/contact.html")
}

// This function is used to define things that must work or else panic.
func must(err error){
	if err != nil {
		panic(err)
	}
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
