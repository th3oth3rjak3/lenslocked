package hash

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"hash"
)

// NewHMAC returns a new HMAC object.
func NewHMAC(key string) HMAC {
	h := hmac.New(sha512.New, []byte(key))
	return HMAC{
		hmac: h,
	}
}

// HMAC is a wrapper around a shared hash object to make it easier to use
// in our code.
type HMAC struct {
	hmac hash.Hash
}

// Hash takes in a new string to hash and returns the hashed value.
// The primary use case for this is hashing a remember token before
// saving into the database. It is also used for user authentication
// when the remember token is provided as part of logging in.
func (h HMAC) Hash(input string) string {
	h.hmac.Reset()
	h.hmac.Write([]byte(input))
	b := h.hmac.Sum(nil)
	return base64.URLEncoding.EncodeToString(b)
}
