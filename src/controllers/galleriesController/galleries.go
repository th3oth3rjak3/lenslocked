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
	galleryService galleriesModel.GalleryService
}

// Instantiates a new Galleries controller.
// This will panic if templates are not parsed correctly.
// Only used during initial startup.
func NewGalleriesController(gs galleriesModel.GalleryService) *GalleriesController {
	return &GalleriesController{
		NewView:        views.NewView("bootstrap", "galleries/new"),
		ShowView:       views.NewView("bootstrap", "galleries/show"),
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
	ctx := r.Context()
	usr := context.User(ctx)
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
	idStr := chi.URLParam(r, "galleryId")
	data := views.Data{}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		data.SetAlert(errorsModel.ErrGalleryNotFound, true)
		gc.ShowView.Render(w, data)
		return
	}
	gallery, err := gc.galleryService.ByID(uint(id))
	if err != nil {
		data.SetAlert(errorsModel.ErrGalleryNotFound, true)
		gc.ShowView.Render(w, data)
		return
	}
	data.Payload = gallery
	gc.ShowView.Render(w, data)
}
