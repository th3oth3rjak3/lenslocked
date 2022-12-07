package lenslocked

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHomePage(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("An unexpected error occurred during testing: %s", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Home)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned the wrong status code. Got: %d, Want %d", status, http.StatusOK)
	}
	expected := "<h1>Welcome to my super awesome site!</h1>"
	actual := rr.Body.String()
	if expected != actual {
		t.Errorf("Got the wrong response body. Got: %s, Want: %s", actual, expected)
	}
}
