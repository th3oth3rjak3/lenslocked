package servicesModel

import (
	"lenslocked/models/galleriesModel"
	"lenslocked/models/imagesModel"
	"lenslocked/models/usersModel"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type ServicesConfig func(*Services) error

func NewServices(cfgs ...ServicesConfig) (*Services, error) {
	var s Services
	for _, cfg := range cfgs {
		if err := cfg(&s); err != nil {
			return nil, err
		}
	}
	return &s, nil
}

func WithLogMode(logMode bool) ServicesConfig {
	return func(s *Services) error {
		s.db.LogMode(logMode)
		return nil
	}
}

func WithGorm(dialect, connectionInfo string) ServicesConfig {
	return func(s *Services) error {
		db, err := gorm.Open(dialect, connectionInfo)
		if err != nil {
			return err
		}
		s.db = db
		return nil
	}
}

func WithUser(hmacKey string) ServicesConfig {
	return func(s *Services) error {
		s.User = usersModel.NewUserService(s.db, hmacKey)
		return nil
	}
}

func WithGallery() ServicesConfig {
	return func(s *Services) error {
		s.Gallery = galleriesModel.NewGalleryService(s.db)
		return nil
	}
}

func WithImages() ServicesConfig {
	return func(s *Services) error {
		s.Image = imagesModel.NewImageService()
		return nil
	}
}

type Services struct {
	Gallery galleriesModel.GalleryService
	User    usersModel.UserService
	Image   imagesModel.ImageService
	db      *gorm.DB
}

// Closes the database connection. It can be deferred if desired.
func (s *Services) Close() error {
	return s.db.Close()
}

// Destructive Reset drops and automigrates all tables and rebuilds them
func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&usersModel.User{}, &galleriesModel.Gallery{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}

// Runs an automigration for all tables in the database.
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&usersModel.User{}, &galleriesModel.Gallery{}).Error
}
