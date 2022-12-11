// This package is used to create new views in the LensLocked application.
// It applies all of the default layout components by using the NewView function.
// Finally, it also defines the View type and its structure.
package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// The NewView function creates a new View when provided a name for the layout definition and any new files for the view.
func NewView(layout string, files ...string) *View {
	processViewNames(files)
	files = append(files, layoutFiles()...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
		Layout:   layout,
	}
}

// The processViewNames function prepends common directory information
// to the front of the filename and appends the extensions on the end.
func processViewNames(files []string) {
	baseDir := "views/"
	extension := ".html"
	for i, file := range files {
		files[i] = baseDir + file + extension
	}
}

// ServeHttp is used to implement the http.Handler interface to render basic views.
func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := v.Render(w, nil); err != nil {
		panic(err)
	}
}

// Render is used to render the view with the predefined layout.
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	return v.Template.ExecuteTemplate(w, v.Layout, nil)
}

// An object to handle webpage Views.
type View struct {
	Template *template.Template
	Layout   string
}

// Uses glob to get all of the template files in the directory.
func layoutFiles() []string {
	layoutDir := "views/layouts/"
	templateExtension := ".html"
	files, err := filepath.Glob(layoutDir + "*" + templateExtension)
	if err != nil {
		panic(err)
	}
	return files
}
