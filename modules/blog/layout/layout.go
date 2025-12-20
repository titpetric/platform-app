package layout

import (
	"context"
	"io"
	"io/fs"

	"github.com/titpetric/vuego"
)

// Renderer handles page and layout rendering with vuego templates
type Renderer struct {
	root fs.FS
	tpl  vuego.Template

	data map[string]any
}

// NewRenderer creates a new Renderer with the given filesystem and shared data
func NewRenderer(root fs.FS, data map[string]any) *Renderer {
	return &Renderer{
		root: root,
		data: data,
		tpl:  vuego.NewFS(root, vuego.WithLessProcessor(), vuego.WithFuncs(Funcs)),
	}
}

// Render loads a template, and if the template contains "layout" in the metadata, it will
// load another template from layouts/%s.vuego; Layouts can be chained so one layout can
// again trigger another layout, like `blog.vuego -> layouts/post.vuego -> layouts/base.vuego`.
func (r *Renderer) Render(ctx context.Context, w io.Writer, filename string, data map[string]any) error {
	return r.tpl.Load(filename).Fill(data).Layout(ctx, w)
}
