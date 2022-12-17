package galleriesModel

import (
	"github.com/jinzhu/gorm"
)

// A Gallery contains image resources that are viewed by our visitors.
type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not null;index"`
	Title  string `gorm:"not null"`
}

// GalleryDB is used to interact with the galleries database.
//
// For all single queries:
// If the gallery is found, error will be nil
// If the galler is not found, the error will be set to ErrGalleryNotFound
type GalleryDB interface {
	ByID(id uint) (*Gallery, error)
	Create(gallery *Gallery) error
}

// GalleryService is a set of methods to manipulate and work with the Gallery model.
type GalleryService interface {
	GalleryDB
}

// NewGalleryService initializes a GalleryService instance.
func NewGalleryService(db *gorm.DB) GalleryService {
	gg := &galleryGorm{db}
	gv := newGalleryValidator(gg)
	return &galleryService{
		GalleryDB: gv,
	}
}

// galleryService implements the GalleryService interface.
type galleryService struct {
	GalleryDB
}
