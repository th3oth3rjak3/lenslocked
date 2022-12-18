package galleriesController

import (
	"fmt"
	"net/http"
	"strconv"

	"lenslocked/context"
	"lenslocked/models/errorsModel"
	"lenslocked/models/galleriesModel"
	"lenslocked/views"

	"github.com/go-chi/chi/v5"
)

// The Galleries controller object.
type GalleriesController struct {
	NewView        *views.View
	ShowView       *views.View
	EditView       *views.View
	galleryService galleriesModel.GalleryService
}

// Instantiates a new Galleries controller.
// This will panic if templates are not parsed correctly.
// Only used during initial startup.
func NewGalleriesController(gs galleriesModel.GalleryService) *GalleriesController {
	return &GalleriesController{
		NewView:        views.NewView("bootstrap", "galleries/new"),
		ShowView:       views.NewView("bootstrap", "galleries/show"),
		EditView:       views.NewView("bootstrap", "galleries/edit"),
		galleryService: gs,
	}
}

// Used to show the "Create a Gallery" page to the user
//
// GET /galleries/new
func (gc *GalleriesController) New(w http.ResponseWriter, r *http.Request) {
	gc.NewView.Render(w, nil)
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
		gc.NewView.Render(w, vd)
		return
	}
	gallery := &galleriesModel.Gallery{
		Title:  formData.Title,
		UserID: usr.ID,
	}
	if err := gc.galleryService.Create(gallery); err != nil {
		vd.SetAlert(err, true)
		gc.NewView.Render(w, vd)
		return
	}
	url := fmt.Sprintf("%s/%d", r.URL.Path, gallery.ID)
	http.Redirect(w, r, url, http.StatusFound)
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
	data.Payload = gallery
	gc.ShowView.Render(w, data)
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
	gc.EditView.Render(w, data)
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
	vd.Payload = gallery
	if err != nil {
		vd.SetAlert(err, true)
		gc.EditView.Render(w, vd)
		return
	}
	formData := &GalleryForm{}
	if err := formData.Bind(r); err != nil {
		vd.SetAlert(err, true)
		gc.EditView.Render(w, vd)
		return
	}
	gallery.Title = formData.Title
	if err := gc.galleryService.Update(gallery); err != nil {
		vd.SetAlert(err, true)
		gc.EditView.Render(w, vd)
		return
	}
	vd.Alert = &views.Alert{
		Level:   views.AlertLevelSuccess,
		Message: "Gallery successfully updated!",
	}
	gc.EditView.Render(w, vd)
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
		gc.ShowView.Render(w, data)
		return nil, err
	}
	gallery, err := gc.galleryService.ByID(uint(id))
	if err != nil {
		data.SetAlert(errorsModel.ErrGalleryNotFound, true)
		gc.ShowView.Render(w, data)
		return nil, err
	}
	return gallery, nil
}
