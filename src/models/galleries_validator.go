package models

type galleryValidator struct {
	GalleryDB
}

func newGalleryValidator(gg *galleryGorm) *galleryValidator {
	return &galleryValidator{
		GalleryDB: gg,
	}
}
