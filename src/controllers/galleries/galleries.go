package galleries

import (
	"fmt"
	"net/http"

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
	formData := &GalleryForm{}
	var vd views.Data
	if err := formData.Bind(r); err != nil {
		vd.SetAlert(err, true)
		gc.NewView.Render(w, vd)
		return
	}
	gallery := &models.Gallery{
		Title: formData.Title,
	}
	if err := gc.galleryService.Create(gallery); err != nil {
		vd.SetAlert(err, true)
		gc.NewView.Render(w, vd)
		return
	}
	fmt.Fprintln(w, *gallery)
}
