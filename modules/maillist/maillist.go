package service

import (
	"context"
	"embed"
	"fmt"
	"net/http"

	"github.com/titpetric/platform"

	"github.com/titpetric/platform-app/modules/maillist/service"
)

//go:embed schema
var schema embed.FS

type MailList struct {
	platform.UnimplementedModule
}

func NewMailList() *MailList {
	return &MailList{}
}

func (*MailList) Name() string {
	return "maillist"
}

func (m *MailList) Start(ctx context.Context) error {
	return service.Migrate(ctx, schema)
}

func (m *MailList) Mount(r platform.Router) error {
	r.Route("/maillist", func(r platform.Router) {
		r.Get("/", m.Index)
		r.Post("/create", m.Create)
	})
	return nil
}

func (m *MailList) Index(w http.ResponseWriter, r *http.Request) {
	perms := service.NewPermissions(r)

	fmt.Fprintf(w, "Permissions: %#v", perms)
}

func (m *MailList) Create(w http.ResponseWriter, r *http.Request) {
}
