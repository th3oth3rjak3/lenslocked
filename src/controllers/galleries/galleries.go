package galleries

import (
	"fmt"
	"net/http"

	"lenslocked/context"
	"lenslocked/models"
	"lenslocked/views"
)

// The Galleries controller object.
type GalleriesController struct {
	NewView        *views.View
	galleryService models.GalleryService
}

// Instantiates a new Galleries controller.
// This will panic if templates are not parsed correctly.
// Only used during initial startup.
func NewGalleriesController(gs models.GalleryService) *GalleriesController {
	return &GalleriesController{
		NewView:        views.NewView("bootstrap", "galleries/new"),
		galleryService: gs,
	}
}

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
	gallery := &models.Gallery{
		Title:  formData.Title,
		UserID: usr.ID,
	}
	if err := gc.galleryService.Create(gallery); err != nil {
		vd.SetAlert(err, true)
		gc.NewView.Render(w, vd)
		return
	}
	fmt.Fprintln(w, *gallery)
}
