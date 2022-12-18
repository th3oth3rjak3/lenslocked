package galleriesModel

import (
	"lenslocked/models/errorsModel"
)

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

// Create ensures that the gallery contains a userID fo the owner of
// the gallery and a title for the gallery.
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

// Update ensures that the gallery has a UserID for the owner of the gallery
// and a title.
func (gv *galleryValidator) Update(gallery *Gallery) error {
	// run normalization/validation
	if err := gv.runGalleryValidationFunctions(
		gallery,
		gv.userIdRequirer,
		gv.titleRequirer,
	); err != nil {
		return err
	}

	return gv.GalleryDB.Update(gallery)
}

// Delete validates a user id and then calls the underlying UserDB Delete method.
func (gv *galleryValidator) Delete(id uint) error {
	var gallery Gallery
	gallery.ID = id
	if err := gv.runGalleryValidationFunctions(
		&gallery,
		gv.idGreaterThan(0),
	); err != nil {
		return err
	}
	return gv.GalleryDB.Delete(id)
}

// idGreaterThan checks to see if the gallery has an ID greater than n.
func (gv *galleryValidator) idGreaterThan(n uint) galleryValidatorFunction {
	return func(gallery *Gallery) error {
		if gallery.ID <= n {
			return errorsModel.ErrIdInvalid
		}
		return nil
	}
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
		return errorsModel.ErrTitleRequired
	}
	return nil
}

// userIdRequirer is a function that requires a userId to not be 0 or nil
func (gv *galleryValidator) userIdRequirer(gallery *Gallery) error {
	if gallery.UserID <= 0 {
		return errorsModel.ErrUserIdRequired
	}
	return nil
}
