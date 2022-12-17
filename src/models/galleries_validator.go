package models

// galleryValidator is a chained type that performs validation and
// normalization of data before being passed to the final GalleryDB implementation
type galleryValidator struct {
	GalleryDB
}

// galleryValidatorFunction is a function signature given to all gallery
// validation functions so that it is easier to iterate over all the
// gallery validation functions and call them in a loop.
type galleryValidatorFunction func(*Gallery) error

// Creates a new instance of the galleryValidator
func newGalleryValidator(gg *galleryGorm) *galleryValidator {
	return &galleryValidator{
		GalleryDB: gg,
	}
}

// Create ensures that the password is not empty, meets the complexity
// requirements, and then generates a hash. It also normalizes the email
// address by setting it to lowercase. It also creates a remember token
// and finally calls the subsequent UserDB layer's Create method.
func (gv *galleryValidator) Create(gallery *Gallery) error {
	// run normalization/validation
	if err := gv.runGalleryValidationFunctions(
		gallery,
		gv.userIdRequirer,
		gv.titleRequirer,
	); err != nil {
		return err
	}

	return gv.GalleryDB.Create(gallery)
}

// runUserValidationFunctions is a function which takes a user object
// and a variadic parameter of validation functions which are each called
// on the user object. This function returns an error if any of the
// validation functions return an error.
func (gv *galleryValidator) runGalleryValidationFunctions(gallery *Gallery, fns ...galleryValidatorFunction) error {
	for _, fn := range fns {
		if err := fn(gallery); err != nil {
			return err
		}
	}
	return nil
}

// titleRequirer is a function that requires a gallery to have a title.
func (gv *galleryValidator) titleRequirer(gallery *Gallery) error {
	if gallery.Title == "" {
		return ErrTitleRequired
	}
	return nil
}

// userIdRequirer is a function that requires a userId to not be 0 or nil
func (gv *galleryValidator) userIdRequirer(gallery *Gallery) error {
	if gallery.UserID <= 0 {
		return ErrUserIdRequired
	}
	return nil
}
