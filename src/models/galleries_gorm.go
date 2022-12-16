package models

import "github.com/jinzhu/gorm"

type galleryGorm struct {
	db *gorm.DB
}

var _ GalleryDB = &galleryGorm{}

// Creates a new gallery and backfills data like ID, CreatedAt, and UpdatedAt fields.
//
// This doesn't check for errors, just returns any errors during processing.
func (gg *galleryGorm) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}
