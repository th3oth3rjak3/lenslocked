package lenslocked

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotFoundPage(t *testing.T) {
	// Path /abcd1234 does not exist
	req, err := http.NewRequest("GET", "/abcd1234", nil)
	if err != nil {
		t.Fatalf("An unexpected error occurred during testing: %s", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(NotFound)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Handler returned the wrong status code. Got: %d, Want %d", status, http.StatusOK)
	}
	expected := "<h1>We could not find the page you were looking for :(</h1>" +
		"<p>Please email us if you keep being sent to an invalid page.</p>"
	actual := rr.Body.String()
	if expected != actual {
		t.Errorf("Got the wrong response body. Got: %s, Want: %s", actual, expected)
	}
}
