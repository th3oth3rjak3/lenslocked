package lenslocked

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContactPage(t *testing.T){
	req, err := http.NewRequest("GET", "/contact", nil)
	if err != nil {
		t.Fatalf("An unexpected error occurred during testing: %s", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Contact)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned the wrong status code. Got: %d, Want %d", status, http.StatusOK)
	}
	expected := "To get in touch, please send an email to <a href='mailto:support@lenslocked.com'>Support</a>"
	actual := rr.Body.String()
	if expected != actual {
		t.Errorf("Got the wrong response body. Got: %s, Want: %s", actual, expected)
	}
}