package view

import "github.com/titpetric/vuego"

// Views is a type that provides type safe view helpers.
type Views struct {
	Loader *Loader
}

// NewViews creates a view object. All views are implemented here.
func NewViews(tpl vuego.Template) *Views {
	return &Views{
		Loader: NewLoader(tpl),
	}
}

// Index renders the blog index/list page
func (v *Views) Index(data *IndexData) vuego.Template {
	return v.Loader.Load("pages/index.vuego").Fill(data)
}

// Blog renders the blog list page
func (v *Views) Blog(data *IndexData) vuego.Template {
	return v.Loader.Load("pages/blog.vuego").Fill(data)
}

// Post renders the post layout template
func (v *Views) Post(data *PostData) vuego.Template {
	return v.Loader.Load("layouts/post.vuego").Fill(data)
}
