package main

import (
	"fmt"
	"log"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1>")
}

func main() {
	http.HandleFunc("/", handlerFunc)
	addr := "localhost:3000"
	log.Println("Listening on", addr)
	log.Println("ugh 123 ")
	log.Fatal(http.ListenAndServe(addr, nil))
}
