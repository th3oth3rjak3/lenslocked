package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func NewServices(connectionInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &Services{
		User:    NewUserService(db),
		Gallery: NewGalleryService(db),
		db:      db,
	}, nil
}

type Services struct {
	Gallery GalleryService
	User    UserService
	db      *gorm.DB
}

// Turns database log mode on or off. This is used primarily for testing purposes.
func (s *Services) LogMode(dbLogModeEnabled bool) {
	s.db.LogMode(dbLogModeEnabled)
}

// Closes the database connection. It can be deferred if desired.
func (s *Services) Close() error {
	return s.db.Close()
}

// Destructive Reset drops and automigrates all tables and rebuilds them
func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&User{}, &Gallery{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}

// Runs an automigration for all tables in the database.
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{}).Error
}
