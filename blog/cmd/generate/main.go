package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/titpetric/platform-app/blog"
	"github.com/titpetric/platform-app/blog/storage"
)

func main() {
	outputDir := flag.String("output", "public", "Output directory for generated files")
	dataDir := flag.String("data", "data", "Data directory for markdown files")
	flag.Parse()

	ctx := context.Background()

	// Initialize platform (database only)
	if err := generate(ctx, *dataDir, *outputDir); err != nil {
		log.Fatalf("generation failed: %v", err)
	}
}

func generate(ctx context.Context, dataDir, outputDir string) error {
	start := time.Now()

	// Get database from platform
	db, err := storage.DB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}

	// Create module and load articles
	module := blog.NewModule(dataDir)

	// Create storage and schema
	repo, err := storage.NewStorage(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}

	// Set repository on module
	module.SetRepository(repo)

	// Scan markdown files
	count, err := module.ScanMarkdownFiles(ctx)
	if err != nil {
		return fmt.Errorf("failed to scan markdown files: %w", err)
	}
	fmt.Printf("Scanned %d markdown files\n", count)

	// Generate static files
	gen := blog.NewGenerator(module, outputDir)
	if err := gen.Generate(ctx); err != nil {
		return fmt.Errorf("generation failed: %w", err)
	}

	duration := time.Since(start)

	fmt.Println("✓ Completed in", duration)
	return nil
}
