package view

import (
	"io/fs"

	"github.com/titpetric/vuego"
)

// Views is a type that provides type safe view helpers.
type Views struct {
	Loader *Loader
}

// AdminViews provides views for the admin panel.
type AdminViews Views

// NewViews creates a view object. All views are implemented here.
func NewViews(filesystem fs.FS) *Views {
	return &Views{
		Loader: NewLoader(vuego.NewFS(filesystem, vuego.WithFuncs(Funcs))),
	}
}

// Index renders the blog index/list page.
func (v *Views) Index(data *IndexData) vuego.Template {
	return v.Loader.Load("pages/index.vuego").Fill(data.Map())
}

// Blog renders the blog list page.
func (v *Views) Blog(data *IndexData) vuego.Template {
	return v.Loader.Load("pages/blog.vuego").Fill(data.Map())
}

// Post renders the post layout template.
func (v *Views) Post(data *PostData) vuego.Template {
	return v.Loader.Load("layouts/post.vuego").Fill(data.Map())
}

// NewAdminViews creates a view object. All views are implemented here.
func NewAdminViews(filesystem fs.FS) *AdminViews {
	return &AdminViews{
		Loader: NewLoader(vuego.NewFS(filesystem, vuego.WithFuncs(Funcs))),
	}
}

// Dashboard renders the admin dashboard page.
func (v *AdminViews) Dashboard(data *AdminDashboardData) vuego.Template {
	return v.Loader.Load("dashboard.vuego").Assign("data", data)
}

// List renders the admin article list page.
func (v *AdminViews) List(data *AdminListData) vuego.Template {
	return v.Loader.Load("list.vuego").Fill(data.Map())
}

// Edit renders the admin article edit page.
func (v *AdminViews) Edit(data *AdminEditData) vuego.Template {
	return v.Loader.Load("edit.vuego").Fill(data.Map())
}
