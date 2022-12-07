package main

import (
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type User struct {
	Name string
	Script string
	Int int
	Kids []string
}

func main() {
	t, err := template.ParseFiles("./home.html")
	if err != nil {
		panic(err)
	}
	jake := User{"Jake Hathaway", "alert('hello')", 33, []string{"Noah", "Sage"}}
	err = t.Execute(os.Stdout, jake)
	if err != nil {
		panic(err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		t.Execute(w, jake)
	})
	http.ListenAndServe("localhost:3000", r)
}
