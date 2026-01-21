package main

import (
	"context"
	"log"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/titpetric/platform"

	"github.com/titpetric/platform-app/blog"
	"github.com/titpetric/platform-app/user"
)

func main() {
	ctx := context.Background()
	if err := start(ctx); err != nil {
		log.Fatalf("exit error: %v", err)
	}
}

func start(ctx context.Context) error {
	opts := platform.NewOptions()
	svc := platform.New(opts)

	svc.Use(middleware.Logger)
	svc.Register(user.NewHandler())
	svc.Register(blog.NewModule("./data"))

	if err := svc.Start(ctx); err != nil {
		return err
	}

	svc.Wait()

	return nil
}
