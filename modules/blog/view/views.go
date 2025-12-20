package view

import (
	"io/fs"

	"github.com/titpetric/platform-app/modules/blog/layout"
)

type Views struct {
	*layout.Renderer
	data map[string]any
}

func NewViews(root fs.FS) (*Views, error) {
	data := map[string]any{}
	if err := fillTemplateData(&data); err != nil {
		return nil, err
	}

	return &Views{
		Renderer: layout.NewRenderer(root, data),
		data:     data,
	}, nil
}
