// This package is used to create new views in the LensLocked application.
// It applies all of the default layout components by using the NewView function.
// Finally, it also defines the View type and its structure.
package views

import (
	"html/template"
)

func NewView(layout string, files ...string) *View {
	files = append(
		files,
		"views/layouts/bootstrap.html",
		"views/layouts/footer.html",
		"views/layouts/navbar.html",
	)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
		Layout: layout,
	}
}

type View struct {
	Template *template.Template
	Layout   string
}
