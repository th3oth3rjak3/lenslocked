package errorsModel

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	// ErrEmailNotFound is returned when a user email cannot be found in the database.
	ErrUserNotFound modelError = "no account exists with this email address"

	// ErrPasswordIncorrect is returned when the user enters an incorrect password
	ErrPasswordIncorrect modelError = "password is incorrect"

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

	// ErrNameRequired is returned when a user fails to provide a name on signup
	ErrNameRequired modelError = "name is required"

	// ErrGalleryNotFound is returned when a gallery cannot be found in the database.
	ErrGalleryNotFound modelError = "gallery does not exist"

	// ErrTitleRequired is returned when a gallery does not contain a title
	ErrTitleRequired modelError = "gallery title is required"

	// ErrIdInvalid is returned when an invalid ID is provided to a method like Delete.
	ErrIdInvalid privateError = "id provided was invalid"

	// ErrEnvironmentUnset is returned when there are missing environment variables
	ErrEnvironmentUnset privateError = "missing required environment variables"

	// ErrRememberTokenTooShort is returned when a remember token is generated
	// with fewer than 64 bytes.
	ErrRememberTokenTooShort privateError = "remember token generated with too few bytes"

	// ErrRememberTokenHashRequired is returned when a remember token hash is not generated.
	ErrRememberHashRequired privateError = "remember token hash is required"

	// ErrUserIdRequired is returned when a gallery is missing a UserID for
	// the user who owns the gallery
	ErrUserIdRequired privateError = "user id is required for each gallery"
)

// modelError is used for errors that are meant to be public to the user.
type modelError string

func (m modelError) Error() string {
	return "models: " + string(m)
}

func (m modelError) Public() string {
	split := strings.Split(string(m), " ")
	split[0] = cases.Title(language.AmericanEnglish).String(split[0])
	return strings.Join(split, " ") + "."
}

// privateError is used for errors that are more internal to the program
// and wouldn't make a lot of sense to a user.
type privateError string

func (e privateError) Error() string {
	return string(e)
}
