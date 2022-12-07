package lenslocked

import (
	"fmt"
	"net/http"
)

func Contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	message := "To get in touch, please send an email to <a href='mailto:support@lenslocked.com'>Support</a>"
	fmt.Fprint(w, message)
}
