package models

import (
	"github.com/jinzhu/gorm"
)

// Returns the first value from a gorm.DB instance which matches the provided object.
// If the object is found, error will be nil.
// If the object is not found, the error will be set to ErrNotFound.
// If some other error occurs, first will return an error with
// more information about what went wrong. This may not be an error
// generated by the models package.
//
// As a general rule, any error but ErrNotFound should probably result in
// an HTTP 500 error.
func First(db *gorm.DB, data interface{}) error {
	return db.First(data).Error
}
