package imagesModel

import (
	"fmt"
	"io"
	"os"
)

type ImageService interface {
	Create(galleryID uint, r io.ReadCloser, filename string) error
	// ByGalleryID(galleryID uint) []string
}

func NewImageService() ImageService {
	return &imageService{}
}

type imageService struct{}

func (is *imageService) Create(galleryID uint, r io.ReadCloser, filename string) error {
	defer r.Close()
	galleryPath, err := is.makeImagePath(galleryID)
	if err != nil {
		return err
	}
	dst, err := os.Create(galleryPath + filename)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, r)
	if err != nil {
		return err
	}
	return nil
}

// makeImagePath ensures that a gallery folder is created for the specific gallery
// to store photos and then returns the string representation of the path.
func (is *imageService) makeImagePath(galleryID uint) (string, error) {
	galleryPath := fmt.Sprintf("images/galleries/%v/", galleryID)
	err := os.MkdirAll(galleryPath, 0o755)
	if err != nil {
		return "", err
	}
	return galleryPath, nil
}
