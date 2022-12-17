package galleries

import (
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
