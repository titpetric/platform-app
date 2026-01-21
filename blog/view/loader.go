package view

import (
	"github.com/titpetric/vuego"
)

// Loader handles page and layout rendering with vuego templates.
type Loader struct {
	tpl vuego.Template
}

// NewLoader creates a new Loader with the given filesystem and shared data.
func NewLoader(tpl vuego.Template) *Loader {
	return &Loader{tpl}
}

// Load loads the filename and returns a template ready for use.
func (r *Loader) Load(filename string) vuego.Template {
	return r.tpl.Load(filename)
}
