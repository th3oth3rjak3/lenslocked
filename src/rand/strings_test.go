package rand

import (
	"math"
	"testing"
)

// Tests to ensure that the requested number of bytes are returned.
func TestBytes(t *testing.T) {
	b, err := Bytes(10)
	if err != nil {
		t.Fatalf("Expected 10 bytes, got an error: %s", err)
	}
	if len(b) != 10 {
		t.Errorf("Expected 10 bytes, got %d bytes", len(b))
	}
}

// Tests the string functino to ensure it returns the correct length of
// base64 encoded string given the number of input bytes. It also ensures
// that two generated strings are different.
func TestString(t *testing.T) {
	byteLength := 10
	expectedLength := getStringLengthFromBytes(byteLength)
	s, err := String(byteLength)
	if err != nil {
		t.Fatalf("Expected a string of length %d but got an error: %s", expectedLength, err)
	}
	if len(s) != 16 {
		t.Errorf("Expected string of length %d but got length %d. String: %s", expectedLength, len(s), s)
	}
	// Expect the strings to be different
	s1, err := String(10)
	if err != nil {
		t.Fatalf("Expected a string of length %d but got an error: %s", expectedLength, err)
	}
	if s == s1 {
		t.Errorf("Expected different strings, but they were the same. String 1: %s, String 2: %s", s, s1)
	}

	byteLength = 17
	expectedLength = getStringLengthFromBytes(byteLength)
	s, err = String(byteLength)
	if err != nil {
		t.Fatalf("Expected a string of length %d but got an error: %s", expectedLength, err)
	}
	if len(s) != expectedLength {
		t.Errorf("Expected string length %d but got length %d, String: %s", expectedLength, len(s), s)
	}
}

// Tests to make sure a remember token of adequate length is generated.
func TestRememberToken(t *testing.T) {
	byteLength := 64
	expectedLength := getStringLengthFromBytes(byteLength)
	token, err := RememberToken()
	if err != nil {
		t.Fatalf("Expected a remember token of length %d, got an error: %s", expectedLength, err)
	}
	if len(token) != expectedLength {
		t.Errorf("Expected token of length %d, got length %d for token %s", expectedLength, len(token), token)
	}
}

// Ensures that Bytes generates an ErrOutOfRange when a negative value is
// provided.
func TestNegativeBytes(t *testing.T) {
	_, err := Bytes(-1)
	if err == nil {
		t.Errorf("Expected -1 to return an error")
	}
	if err != ErrOutOfRange {
		t.Errorf("Expected ErrOutOfRange, Got: %s", err)
	}
}

func TestHelperFunctions(t *testing.T) {
	expected := 88
	actual := getStringLengthFromBytes(64)
	if expected != actual {
		t.Errorf("Have: %d, Want: %d", actual, expected)
	}
}

// Returns the length of a base64 encoded string given the number of input
// bytes as nBytes.
func getStringLengthFromBytes(nBytes int) int {
	bits := nBytes * 8
	bitsPerChar := 6
	chars := (float64(bits) / float64(bitsPerChar))
	// In this case there was a decimal, so we need to round up to the next whole int.
	if math.Trunc(chars) < chars {
		chars = math.Trunc(chars) + 1
	}

	stringSize := int(chars)
	if stringSize%4 != 0 {
		stringSize += 4 - (stringSize % 4)
	}
	return stringSize
}
