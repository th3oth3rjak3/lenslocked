package models

import (
	"os"
	"strings"

	"lenslocked/hash"
	"lenslocked/rand"

	"golang.org/x/crypto/bcrypt"
)

// userValidator is a chained type that performs validation and
// normalization of data before being passed to the final UserDB implementation
type userValidator struct {
	UserDB
	hmac hash.HMAC
}

// userValidationFunction is a function signature given to all user
// validation functions so that it is easier to iterate over all the
// user validation functions and call them in a loop.
type userValidationFunction func(*User) error

// Creates a new instance of the userValidator
func newUserValidator(connectionInfo string) (*userValidator, error) {
	key := os.Getenv("HASH_KEY")
	if key == "" {
		return nil, ErrEnvironmentUnset
	}
	gorm, err := newUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}
	hmac := hash.NewHMAC(key)
	return &userValidator{
		UserDB: gorm,
		hmac:   hmac,
	}, nil
}

// Query methods

// ByID will ensure the ID is valid and then call the ByID method
// on the subsequent UserDB layer.
func (uv *userValidator) ByID(id uint) (*User, error) {
	// TODO: validate ByID data
	return uv.UserDB.ByID(id)
}

// ByEmail will first convert the email to lowercase and then call
// ByEmail on the subsequent UserDB layer.
func (uv *userValidator) ByEmail(email string) (*User, error) {
	// TODO: validate by email
	// normalize email input
	email = strings.ToLower(email)
	return uv.UserDB.ByEmail(email)
}

// ByRemember will hash the token and then call ByRemember on the
// subsequent UserDB layer.
func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := &User{
		Remember: token,
	}
	if err := uv.runUserValidationFunctions(
		user,
		uv.hmacRemember,
	); err != nil {
		return nil, err
	}
	return uv.UserDB.ByRemember(user.RememberHash)
}

// Data Alteration Methods

// Create ensures that the password is not empty, meets the complexity
// requirements, and then generates a hash. It also normalizes the email
// address by setting it to lowercase. It also creates a remember token
// and finally calls the subsequent UserDB layer's Create method.
func (uv *userValidator) Create(user *User) error {
	// TODO: check complexity requirements.
	// TODO: name not empty string
	// TODO: email not already taken.
	if user.Password == "" {
		return ErrInvalidPassword
	}

	// normalize email address
	user.Email = strings.ToLower(user.Email)

	// run normalization/validation
	if err := uv.runUserValidationFunctions(
		user,
		uv.bcryptPassword,
		uv.generateDefaultRemember,
		uv.hmacRemember,
	); err != nil {
		return err
	}

	return uv.UserDB.Create(user)
}

func (uv *userValidator) Update(user *User) error {
	// TODO: email validation
	// TODO: password validation
	// TODO: password complexity
	// TODO: check name is not empty string
	// TODO: check email is not already taken
	// TODO: validation
	if err := uv.runUserValidationFunctions(
		user,
		uv.hmacRemember,
	); err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

func (uv *userValidator) Delete(id uint) error {
	// TODO: validation
	if id < 1 {
		return ErrInvalidId
	}
	return uv.UserDB.Delete(id)
}

// runUserValidationFunctions is a function which takes a user object
// and a variadic parameter of validation functions which are each called
// on the user object. This function returns an error if any of the
// validation functions return an error.
func (uv *userValidator) runUserValidationFunctions(user *User, fns ...userValidationFunction) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

// WARNING: bcryptPassword does not validate complexity requirements for a user
// password. It will only hash passwords that are not an empty string.
func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	pwBytes := []byte(user.Password)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = "" // Clear the user's actual password
	return nil
}

// hmacRemember takes a User object with a remember token set,
// hashes the remember token, and sets the user.RememberHash value.
//
// WARNING: If the remember token is the empty string, it returns
// without performing a hash.
func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

// generateDefaultRemember generates a new remember token if one is not set.
func (uv *userValidator) generateDefaultRemember(user *User) error {
	if user.Remember != "" {
		return nil
	}
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = token
	return nil
}
