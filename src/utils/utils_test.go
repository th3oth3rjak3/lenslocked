package lenslocked

import "testing"

func TestHello(t *testing.T){
	expected := "Hello, world!"
	actual := Greeting()
	if actual != expected {
		t.Fail()
	}
}