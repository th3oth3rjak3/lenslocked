package usersModel

import (
	"lenslocked/models/errorsModel"

	"github.com/jinzhu/gorm"
)

// The User object is a GORM model that that represents the user's information.
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

// userGorm implements the UserDB interface
type userGorm struct {
	db *gorm.DB
}

// Creates a provided user and backfills data like the ID, CreatedAt, and UpdatedAt fields.
//
// This doesn't check for errors, just returns any errors during processing.
func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

// Returns the first value from a gorm.DB instance which matches the user object.
// If the user is found, error will be nil.
// If the user is not found, the error will be set to ErrNotFound.
// If some other error occurs, first will return an error with
// more information about what went wrong. This may not be an error
// generated by the models package.
//
// As a general rule, any error but ErrNotFound should probably result in
// an HTTP 500 error.
func first(db *gorm.DB, user *User) error {
	err := db.First(user).Error
	if err == gorm.ErrRecordNotFound {
		return errorsModel.ErrUserNotFound
	}
	return err
}

// ByID will look up a user with the provided ID.
// If the user is found, error will be nil.
// If the user is not found, the error will be set to ErrNotFound.
// If some other error occurs, ByID will return an error with
// more information about what went wrong. This may not be an error
// generated by the models package.
//
// As a general rule, any error but ErrNotFound should probably result in
// an HTTP 500 error.
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

// Returns a user by the provided email address.
// If the user is found, error will be nil.
// If the user is not found, the error will be set to ErrNotFound.
// If some other error occurs, ByEmail will return an error with
// more information about what went wrong. This may not be an error
// generated by the models package.
//
// As a general rule, any error but ErrNotFound should probably result in
// an HTTP 500 error.
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// ByRemember uses a remember token to look up a user in the database
// who has the matching token. This method expects the remember value
// to already be hashed.
// If the user is found, error will be nil.
// If the user is not found, the error will be set to ErrNotFound.
// If some other error occurs, ByRemember will return an error with
// more information about what went wrong. This may not be an error
// generated by the models package.
//
// As a general rule, any error but ErrNotFound should probably result in
// an HTTP 500 error.
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var usr User
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &usr)
	if err != nil {
		return nil, err
	}
	return &usr, nil
}

// Updates a user in the database. This update method requires a full user object
// because it overwrites the existing user object. This would be similar to an HTTP PUT,
// rather than an HTTP PATCH method.
func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

// Delete will delete the user with the provided ID.
func (ug *userGorm) Delete(id uint) error {
	usr := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&usr).Error
}
