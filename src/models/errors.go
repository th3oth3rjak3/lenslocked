package models

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	// ErrEmailNotFound is returned when a user email cannot be found in the database.
	ErrUserNotFound modelError = "no account exists with this email address"

	// ErrIdInvalid is returned when an invalid ID is provided to a method like Delete.
	ErrIdInvalid modelError = "id provided was invalid"

	// ErrPasswordIncorrect is returned when the user enters an incorrect password
	ErrPasswordIncorrect modelError = "password is incorrect"

	// ErrEnvironmentUnset is returned when there are missing environment variables
	ErrEnvironmentUnset modelError = "missing required environment variables"

	// ErrEmailMissing is returned when an email is an empty string
	ErrEmailMissing modelError = "email address is required"

	// ErrEmailInvalid is returned when a user provides an email that does
	// not match the allowed pattern.
	ErrEmailInvalid modelError = "email address is not correctly formatted"

	// ErrEmailTaken is returned when a user tries to update or create a
	// user object with an email that is already owned by another user.
	ErrEmailTaken modelError = "email address is already taken"

	// ErrPasswordTooShort is returned when a user provides a password that does not
	// meet the minimum length requirement.
	ErrPasswordTooShort modelError = "password must be at least 8 characters long"

	// ErrPasswordRequired is returned when a user does not provide a password.
	ErrPasswordRequired modelError = "password is required"

	// ErrRememberTokenTooShort is returned when a remember token is generated
	// with fewer than 64 bytes.
	ErrRememberTokenTooShort modelError = "remember token generated with too few bytes"

	// ErrRememberTokenHashRequired is returned when a remember token hash is not generated.
	ErrRememberHashRequired modelError = "remember token hash is required"

	// ErrNameRequired is returned when a user fails to provide a name on signup
	ErrNameRequired modelError = "name is required"

	// ErrGalleryNotFound is returned when a gallery cannot be found in the database.
	ErrGalleryNotFound modelError = "gallery does not exist"
)

type modelError string

func (m modelError) Error() string {
	return "models: " + string(m)
}

func (m modelError) Public() string {
	split := strings.Split(string(m), " ")
	split[0] = cases.Title(language.AmericanEnglish).String(split[0])
	return strings.Join(split, " ") + "."
}
