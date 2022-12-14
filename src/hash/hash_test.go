package hash

import (
	"encoding/base64"
	"testing"
)

// Tests the Hash function to ensure that a manually generated hash and
// the hash function return the same result given the same key and
// the same string value as inputs.
func TestHashFunction(t *testing.T) {
	key := "super-secret-key-that-noone-will-ever-guess"
	thingToHash := "this is some string that we want to hash"
	h := NewHMAC(key)
	h.hmac.Write([]byte(thingToHash))
	b := h.hmac.Sum(nil)
	expected := base64.URLEncoding.EncodeToString(b)
	actual := h.Hash(thingToHash)
	if expected != actual {
		t.Errorf("Hashes don't match. Have: %s, Want: %s", actual, expected)
	}
}
