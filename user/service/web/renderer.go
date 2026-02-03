package web

import (
	"io/fs"

	"github.com/titpetric/vuego"
	"github.com/titpetric/vuego-cli/basecoat"

	"github.com/titpetric/platform-app/user/view"
)

// Renderer handles page and layout rendering with vuego templates
type Renderer struct {
	fs    fs.FS
	vuego vuego.Template

	data map[string]any
}

// NewRenderer creates a new Renderer with the given filesystem and shared data
func NewRenderer(data map[string]any) *Renderer {
	if data == nil {
		data = make(map[string]any)
	}

	ofs := vuego.NewOverlayFS(view.FS, basecoat.FS)

	return &Renderer{
		fs:    ofs,
		data:  data,
		vuego: vuego.NewFS(ofs, vuego.WithLessProcessor(), vuego.WithFuncs(Funcs)).Fill(data),
	}
}

// Render loads a template, and if the template contains "layout" in the metadata, it will
// load another template from layouts/%s.vuego; Layouts can be chained so one layout can
// again trigger another layout, like `blog.vuego -> layouts/post.vuego -> layouts/base.vuego`.
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
