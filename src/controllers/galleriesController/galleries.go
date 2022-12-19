package galleriesController

import (
	"fmt"
	"net/http"
	"strconv"

	"lenslocked/context"
	"lenslocked/models/errorsModel"
	"lenslocked/models/galleriesModel"
	"lenslocked/models/imagesModel"
	"lenslocked/views"

	"github.com/go-chi/chi/v5"
)

const (
	// The bit shift is to convert to MB. In the case below, we have
	// 1 << 20 to represent 1 MB.
	MAX_MULTIPART_MEMORY = 1 << 20
)

// The Galleries controller object.
type GalleriesController struct {
	NewView        *views.View
	ShowView       *views.View
	EditView       *views.View
	IndexView      *views.View
	galleryService galleriesModel.GalleryService
	imageService   imagesModel.ImageService
}

// Instantiates a new Galleries controller.
// This will panic if templates are not parsed correctly.
// Only used during initial startup.
func NewGalleriesController(gs galleriesModel.GalleryService, is imagesModel.ImageService) *GalleriesController {
	return &GalleriesController{
		NewView:        views.NewView("bootstrap", "galleries/new"),
		ShowView:       views.NewView("bootstrap", "galleries/show"),
		EditView:       views.NewView("bootstrap", "galleries/edit"),
		IndexView:      views.NewView("bootstrap", "galleries/index"),
		galleryService: gs,
		imageService:   is,
	}
}

// Used to show the "Create a Gallery" page to the user
//
// GET /galleries/new
func (gc *GalleriesController) New(w http.ResponseWriter, r *http.Request) {
	gc.NewView.Render(w, r, nil)
}

// Used to process the new gallery request
//
// POST /galleries
func (gc *GalleriesController) Create(w http.ResponseWriter, r *http.Request) {
	usr := context.User(r.Context())
	if usr == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	formData := &GalleryForm{}
	var vd views.Data
	if err := formData.Bind(r); err != nil {
		vd.SetAlert(err, true)
		gc.NewView.Render(w, r, vd)
		return
	}
	gallery := &galleriesModel.Gallery{
		Title:  formData.Title,
		UserID: usr.ID,
	}
	if err := gc.galleryService.Create(gallery); err != nil {
		vd.SetAlert(err, true)
		gc.NewView.Render(w, r, vd)
		return
	}
	url := fmt.Sprintf("%s/%d/edit", r.URL.Path, gallery.ID)
	http.Redirect(w, r, url, http.StatusFound)
}

// Shows all galleries owned by a user
//
// GET /galleries
func (gc *GalleriesController) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	galleries, err := gc.galleryService.ByUserID(user.ID)
	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	var vd views.Data
	vd.Payload = galleries
	gc.IndexView.Render(w, r, vd)
}

// Get a specific gallery by the ID
//
// GET /galleries/:id
func (gc *GalleriesController) Show(w http.ResponseWriter, r *http.Request) {
	data := views.Data{}
	gallery, err := gc.galleryById(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	data.Payload = gallery
	gc.ShowView.Render(w, r, data)
}

// Edit a specific gallery by the ID
//
// GET /galleries/:id/edit
func (gc *GalleriesController) Edit(w http.ResponseWriter, r *http.Request) {
	data := views.Data{}
	gallery, err := gc.galleryById(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	data.Payload = gallery
	gc.EditView.Render(w, r, data)
}

// Used to process the updated gallery request
//
// POST /galleries/:id/update
func (gc *GalleriesController) Update(w http.ResponseWriter, r *http.Request) {
	usr := context.User(r.Context())
	if usr == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	var vd views.Data
	gallery, err := gc.galleryById(w, r)
	if err != nil {
		vd.SetAlert(err, true)
		gc.EditView.Render(w, r, vd)
		return
	}
	if gallery.UserID != usr.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	vd.Payload = gallery
	formData := &GalleryForm{}
	if err := formData.Bind(r); err != nil {
		vd.SetAlert(err, true)
		gc.EditView.Render(w, r, vd)
		return
	}
	gallery.Title = formData.Title
	if err := gc.galleryService.Update(gallery); err != nil {
		vd.SetAlert(err, true)
		gc.EditView.Render(w, r, vd)
		return
	}
	vd.Alert = &views.Alert{
		Level:   views.AlertLevelSuccess,
		Message: "Gallery successfully updated!",
	}
	gc.EditView.Render(w, r, vd)
}

// Used to delete a gallery by its given id
//
// POST /galleries/:id/delete
func (gc *GalleriesController) Delete(w http.ResponseWriter, r *http.Request) {
	usr := context.User(r.Context())
	if usr == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	var vd views.Data
	gallery, err := gc.galleryById(w, r)
	if err != nil {
		vd.SetAlert(err, true)
		gc.EditView.Render(w, r, vd)
		return
	}
	if gallery.UserID != usr.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	vd.Payload = gallery
	if err := gc.galleryService.Delete(gallery.ID); err != nil {
		vd.SetAlert(err, true)
		gc.EditView.Render(w, r, vd)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

// Used to process the updated gallery image uploads
//
// POST /galleries/:id/images
func (gc *GalleriesController) Upload(w http.ResponseWriter, r *http.Request) {
	usr := context.User(r.Context())
	if usr == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	var vd views.Data
	gallery, err := gc.galleryById(w, r)
	if err != nil {
		vd.SetAlert(err, true)
		gc.EditView.Render(w, r, vd)
		return
	}
	if gallery.UserID != usr.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	vd.Payload = gallery
	// TODO: parse a multipart form with multiple images.
	err = r.ParseMultipartForm(MAX_MULTIPART_MEMORY)
	if err != nil {
		vd.SetAlert(err, true)
		gc.EditView.Render(w, r, vd)
		return
	}

	files := r.MultipartForm.File["images"]
	for _, f := range files {
		file, err := f.Open()
		if err != nil {
			vd.SetAlert(err, true)
			gc.EditView.Render(w, r, vd)
			return
		}
		defer file.Close()
		err = gc.imageService.Create(gallery.ID, file, f.Filename)
		if err != nil {
			vd.SetAlert(err, true)
			gc.EditView.Render(w, r, vd)
			return
		}
	}
	fmt.Fprintln(w, "Files successfully uploaded.")
}

// galleryById gets a gallery by the id passed in the URL params if one exists.
// It then returns that gallery and an error if one occurs. This helper function
// is used for the Show and Edit methods.
func (gc *GalleriesController) galleryById(w http.ResponseWriter, r *http.Request) (*galleriesModel.Gallery, error) {
	idStr := chi.URLParam(r, "galleryId")
	data := views.Data{}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		data.SetAlert(errorsModel.ErrGalleryNotFound, true)
		gc.ShowView.Render(w, r, data)
		return nil, err
	}
	gallery, err := gc.galleryService.ByID(uint(id))
	if err != nil {
		data.SetAlert(errorsModel.ErrGalleryNotFound, true)
		gc.ShowView.Render(w, r, data)
		return nil, err
	}
	return gallery, nil
}
