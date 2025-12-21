package view

import (
	"io/fs"

	"github.com/titpetric/vuego"
)

// Views is a type that provides type safe view helpers.
type Views struct {
	Loader *Loader
}

// NewViews creates a view object. All views are implemented here.
func NewViews(filesystem fs.FS) *Views {
	return &Views{
		Loader: NewLoader(vuego.NewFS(filesystem, vuego.WithLessProcessor(), vuego.WithFuncs(Funcs))),
	}
}

// Index provides the index page view template.
func (v *Views) Index(data Data) vuego.Template {
	return v.Loader.Load("login.vuego").Fill(data)
}
