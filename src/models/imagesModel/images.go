package imagesModel

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Image is not stored in the database
type Image struct {
	GalleryID uint
	Filename  string
}

func (i *Image) Path() string {
	url := url.URL{
		Path: fmt.Sprintf("/images/galleries/%v/%v", i.GalleryID, i.Filename),
	}
	return url.String()
}

type ImageService interface {
	Create(galleryID uint, r io.ReadCloser, filename string) error
	Delete(galleryID uint, filename string) error
	ByGalleryID(galleryID uint) ([]Image, error)
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

	for is.fileExists(galleryPath + filename) {

		fileSlice := strings.Split(filename, ".")
		if len(fileSlice) > 1 {
			// Add the extension to the end of the slice. We will overwrite the old extension with the _copy value
			fileSlice[len(fileSlice)-2] += "_copy"
			filename = strings.Join(fileSlice, ".")
		}
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

func (is *imageService) ByGalleryID(galleryID uint) ([]Image, error) {
	path := is.imagePath(galleryID)
	names, err := filepath.Glob(path + "*")
	if err != nil {
		return nil, err
	}
	extensions := []string{"jpg", "jpeg", "png"}
	names = is.filter(names, extensions)
	images := is.toImages(galleryID, names)
	return images, nil
}

func (is *imageService) Delete(galleryID uint, filename string) error {
	path := is.imagePath(galleryID) + filename
	return os.Remove(path)
}

func (is *imageService) toImages(galleryID uint, names []string) []Image {
	var images []Image
	path := is.imagePath(galleryID)
	for i := range names {
		names[i] = strings.Replace(names[i], path, "", 1)
		image := Image{Filename: names[i], GalleryID: galleryID}
		images = append(images, image)
	}
	return images
}

func (is *imageService) filter(names []string, extensions []string) []string {
	var imageStrings []string
	for _, fpath := range names {
		fpathSlice := strings.Split(fpath, ".")
		if contains(extensions, string(fpathSlice[len(fpathSlice)-1])) {
			imageStrings = append(imageStrings, fpath)
		}
	}
	return imageStrings
}

func contains(okExtensions []string, val string) bool {
	for _, ext := range okExtensions {
		if strings.EqualFold(ext, val) {
			return true
		}
	}
	return false
}

func (is *imageService) imagePath(galleryID uint) string {
	return fmt.Sprintf("images/galleries/%v/", galleryID)
}

// makeImagePath ensures that a gallery folder is created for the specific gallery
// to store photos and then returns the string representation of the path.
func (is *imageService) makeImagePath(galleryID uint) (string, error) {
	galleryPath := is.imagePath(galleryID)
	err := os.MkdirAll(galleryPath, 0o755)
	if err != nil {
		return "", err
	}
	return galleryPath, nil
}

func (is *imageService) fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
