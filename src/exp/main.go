package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"fmt"

	"lenslocked/rand"
)

func main() {
	token, err := rand.RememberToken()
	if err != nil {
		panic(err)
	}
	h := hmac.New(sha512.New, []byte("this-is-my-secret-key"))
	toHash := []byte(token)
	h.Write(toHash)
	b := h.Sum(nil)
	fmt.Println(token)
	fmt.Println(base64.URLEncoding.EncodeToString(b))
}
