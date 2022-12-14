package rand

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
)

// The number of bytes to use as the length of each remember token.
const REMEMBER_TOKEN_BYTES = 64

var ErrOutOfRange = errors.New("strings: out of range exception")

// Bytes will help us generate n random bytes or will return an
// error if there was one. This uses the crypto/rand package so
// it is safe to use with things like remember tokens.
func Bytes(n int) ([]byte, error) {
	if n < 0 {
		return nil, ErrOutOfRange
	}
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Generates a random, base64 encoded string of size nBytes. This uses
// the crypto/rand package so it is safe to use with things like remember
// tokens.
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Generates a remember token of constant size.
func RememberToken() (string, error) {
	return String(REMEMBER_TOKEN_BYTES)
}
