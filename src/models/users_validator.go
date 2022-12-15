package models

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"lenslocked/hash"
	"lenslocked/rand"

	"golang.org/x/crypto/bcrypt"
)

const MIN_PASSWORD_LENGTH = 8

var (
	// ErrEmailMissing is returned when an email is an empty string
	ErrEmailMissing = errors.New("models: email address is required")

	// ErrEmailInvalid is returned when a user provides an email that does
	// not match the allowed pattern.
	ErrEmailInvalid = errors.New("models: email address is not correctly formatted")

	// ErrEmailTaken is returned when a user tries to update or create a
	// user object with an email that is already owned by another user.
	ErrEmailTaken = errors.New("models: email address is already taken")

	// ErrPasswordTooShort is returned when a user provides a password that does not
	// meet the minimum length requirement.
	ErrPasswordTooShort = fmt.Errorf("models: password must be at least %d characters long", MIN_PASSWORD_LENGTH)

	// ErrPasswordRequired is returned when a user does not provide a password.
	ErrPasswordRequired = errors.New("models: password is required")
)

// userValidator is a chained type that performs validation and
// normalization of data before being passed to the final UserDB implementation
type userValidator struct {
	UserDB
	hmac       hash.HMAC
	emailRegex *regexp.Regexp
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
	regex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`)
	return &userValidator{
		UserDB:     gorm,
		hmac:       hmac,
		emailRegex: regex,
	}, nil
}

// Query methods

// ByID will ensure the ID is valid and then call the ByID method
// on the subsequent UserDB layer.
func (uv *userValidator) ByID(id uint) (*User, error) {
	var user User
	user.ID = id
	if err := uv.runUserValidationFunctions(
		&user,
		uv.idGreaterThan(0),
	); err != nil {
		return nil, err
	}
	return uv.UserDB.ByID(id)
}

// ByEmail will first convert the email to lowercase and then call
// ByEmail on the subsequent UserDB layer.
func (uv *userValidator) ByEmail(email string) (*User, error) {
	var user User
	user.Email = email
	if err := uv.runUserValidationFunctions(
		&user,
		uv.emailNormalizer,
	); err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(user.Email)
}

// ByRemember will hash the token and then call ByRemember on the
// subsequent UserDB layer.
func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := &User{
		Remember: token,
	}
	if err := uv.runUserValidationFunctions(
		user,
		uv.rememberTokenHasher,
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
	// run normalization/validation
	if err := uv.runUserValidationFunctions(
		user,
		uv.passwordRequirer,
		uv.passwordMinLengthChecker,
		uv.passwordCryptographer,
		uv.passwordHashRequirer,
		uv.rememberTokenGenerator,
		uv.rememberTokenHasher,
		uv.emailNormalizer,
		uv.emailRequirer,
		uv.emailPatternMatcher,
		uv.emailAvailabilityChecker,
	); err != nil {
		return err
	}

	return uv.UserDB.Create(user)
}

// Update validates a user and then calls the underlying UserDB Update method.
func (uv *userValidator) Update(user *User) error {
	if err := uv.runUserValidationFunctions(
		user,
		uv.passwordMinLengthChecker,
		uv.passwordCryptographer,
		uv.passwordHashRequirer,
		uv.rememberTokenHasher,
		uv.emailNormalizer,
		uv.emailRequirer,
		uv.emailPatternMatcher,
		uv.emailAvailabilityChecker,
	); err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

// Delete validates a user id and then calls the underlying UserDB Delete method.
func (uv *userValidator) Delete(id uint) error {
	var user User
	user.ID = id
	if err := uv.runUserValidationFunctions(
		&user,
		uv.idGreaterThan(0),
	); err != nil {
		return err
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

// WARNING: passwordCryptographer does not validate complexity requirements for a user
// password. It will only hash passwords that are not an empty string.
func (uv *userValidator) passwordCryptographer(user *User) error {
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

// rememberTokenHasher takes a User object with a remember token set,
// hashes the remember token, and sets the user.RememberHash value.
//
// WARNING: If the remember token is the empty string, it returns
// without performing a hash.
func (uv *userValidator) rememberTokenHasher(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

// rememberTokenGenerator generates a new remember token if one is not set.
func (uv *userValidator) rememberTokenGenerator(user *User) error {
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

// idGreaterThan checks to see if the user has an ID greater than n.
func (uv *userValidator) idGreaterThan(n uint) userValidationFunction {
	return func(user *User) error {
		if user.ID <= n {
			return ErrIdInvalid
		}
		return nil
	}
}

// emailNormalizer handles all of the normalization required for a user's
// email address. This includes forcing it to lowercase, and trimming
// off any whitespace characters.
func (uv *userValidator) emailNormalizer(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

// emailRequirer assumes that the email address has already been
// normalized. This means that emailRequirer would expect an email address
// that contained, for example, 5 empty spaces to be equal to the empty string.
func (uv *userValidator) emailRequirer(user *User) error {
	if user.Email == "" {
		return ErrEmailMissing
	}
	return nil
}

// emailPatternMatcher checks to make sure that an email address for a user
// matches a regular expression to ensure email addresses are not malformed.
func (uv *userValidator) emailPatternMatcher(user *User) error {
	if !uv.emailRegex.MatchString(user.Email) {
		return ErrEmailInvalid
	}
	return nil
}

// emailAvailabilityChecker checks to see if the email is unused in the database.
func (uv *userValidator) emailAvailabilityChecker(user *User) error {
	testUser, err := uv.ByEmail(user.Email)
	switch err {
	case ErrNotFound:
		return nil
	case nil:
		if testUser.ID == user.ID {
			return nil
		} else {
			return ErrEmailTaken
		}
	default:
		return err
	}
}

// passwordMinLengthChecker checks to see if a password meets the minimum length
// requirements.
//
// WARNING: this does not check for the empty string. That condition
// should be handled by another validation method such as passwordRequirer
func (uv *userValidator) passwordMinLengthChecker(user *User) error {
	if user.Password == "" || len(user.Password) >= MIN_PASSWORD_LENGTH {
		return nil
	}
	return ErrPasswordTooShort
}

// passwordRequirer checks to see if a password was provided or not.
func (uv *userValidator) passwordRequirer(user *User) error {
	if user.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

// passwordHashRequirer is a developer helper function that ensure a password
// hash is being generated before storing the user into the database.
func (uv *userValidator) passwordHashRequirer(user *User) error {
	if user.PasswordHash == "" {
		return ErrPasswordRequired
	}
	return nil
}
