package web

import (
	"io/fs"

	"github.com/titpetric/vuego"
)

// Renderer handles page and layout rendering with vuego templates.
type Renderer struct {
	viewFS fs.FS
	vuego  vuego.Template

	data map[string]any
}

// NewRenderer creates a new Renderer with the given filesystem and shared data.
func NewRenderer(viewFS fs.FS, data map[string]any) *Renderer {
	if data == nil {
		data = make(map[string]any)
	}

	return &Renderer{
		data:  data,
		vuego: vuego.NewFS(viewFS, vuego.WithLessProcessor(), vuego.WithFuncs(Funcs)).Fill(data),
	}
}

// Load loads a template, chaining any layouts declared in template metadata.
func (r *Renderer) Load(filename string, data any) vuego.Template {
	return r.vuego.Load(filename).Fill(data)
}

func (r *Renderer) Login(data LoginData) vuego.Template {
	return r.Load("login.vuego", data)
}

func (r *Renderer) Logout(data LogoutData) vuego.Template {
	return r.Load("logout.vuego", data)
}

func (r *Renderer) Register(data RegisterData) vuego.Template {
	return r.Load("register.vuego", data)
}
