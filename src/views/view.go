// This package is used to create new views in the LensLocked application.
// It applies all of the default layout components by using the NewView function.
// Finally, it also defines the View type and its structure.
package views

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"lenslocked/context"
)

// The NewView function creates a new View when provided a name for the layout definition and any new files for the view.
func NewView(layout string, files ...string) *View {
	processViewNames(files)
	files = append(files, layoutFiles()...)
	t, err := template.New("").Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			return "", errors.New("CSRF NOT IMPLEMENTED")
		},
	}).ParseFiles(files...)
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
	v.Render(w, r, nil)
}

// Render is used to render the view with the predefined layout.
func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	var vd Data
	// Handle various data types:
	switch d := data.(type) {
	case Data:
		vd = d

	default:
		vd = Data{
			Payload: data,
		}
	}
	vd.User = context.User(r.Context())
	var buf bytes.Buffer
	csrfCookie, err := r.Cookie("_csrf")
	if err != nil {
		panic(err)
	}
	tpl := v.Template.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return template.HTML("<input type='hidden' name='csrf' value='" + csrfCookie.Value + "'/>")
		},
	})
	if err := tpl.ExecuteTemplate(&buf, v.Layout, vd); err != nil {
		log.Printf("Rendering Error Occurred: %s\n", err)
		http.Error(w, "Something went wrong. If the problem persists, please email support.", http.StatusInternalServerError)
		return
	}
	// This could return an error, but we don't have a good way of handling it.
	io.Copy(w, &buf)
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
