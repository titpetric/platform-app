# Blog Module

A blog module for the platform that provides article management, search, and templating functionality.

## Features

- **Markdown-based articles** with YAML front matter
- **Full-text search** with SQLite FTS5
- **Template rendering** with vuego
- **Syntax highlighting** for code blocks with Chroma
- **HTML and JSON APIs** with content negotiation
- **Cache control** for optimal performance

## Quick Start

### Configuration

Set the database connection:

```bash
export PLATFORM_DB_BLOG="sqlite:///tmp/blog.db"  # or sqlite://:memory: for development
```

### Create Articles

Add markdown files to `data/` directory:

```markdown
---
title: My Article
description: A brief description
date: 2024-01-15
layout: post
ogImage: /path/to/image.png
---

# Article Content

Your markdown content here...
```

### Run

```bash
task          # Start the service
task test     # Run tests
task cover    # Generate coverage report
```

## API Endpoints

| Method | Path                        | Response               |
|--------|-----------------------------|------------------------|
| GET    | `/api/blog/articles`        | JSON array of articles |
| GET    | `/api/blog/articles/{slug}` | Single article JSON    |
| GET    | `/api/blog/search?q=query`  | Search results         |
| GET    | `/blog/`                    | Article list (HTML)    |
| GET    | `/blog/{slug}`              | Article detail (HTML)  |

## Architecture

The module consists of:

- **Module** (`blog.go`) - Lifecycle management and route registration
- **Handlers** (`handlers.go`) - HTTP request handling
- **Storage** (`storage/`) - Database operations
- **Models** (`model/`) - Data types
- **Templates** (`template/`, `theme/`) - HTML rendering

## Documentation

- [Architecture](docs/ARCHITECTURE.md) - System design and patterns
- [Implementation](docs/IMPLEMENTATION.md) - Module implementation details
- [Setup Guide](docs/SETUP.md) - Configuration and getting started
- [Markdown](docs/MARKDOWN.md) - Syntax highlighting and rendering
- [Porting Guide](docs/PORTING.md) - Migration from Eleventy/WebC

## Testing

```bash
task test                 # Run all tests
task test:unit           # Unit tests only
task test:integration    # Integration tests only
task cover               # Generate coverage report
task bench               # Run benchmarks
```

## Dependencies

- Go 1.25.4+
- SQLite 3.x
- [platform](https://github.com/titpetric/platform)
- [vuego](https://github.com/titpetric/vuego)
- [blackfriday](https://github.com/russross/blackfriday)
- [chroma](https://github.com/alecthomas/chroma)
