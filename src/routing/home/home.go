package lenslocked

import (
	"fmt"
	"net/http"
	
)

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	message := "Welcome to my super awesome site!"
	fmt.Fprintf(w, "<h1>%s</h1>", message)
}
