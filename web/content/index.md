# Req - Test APIs with Terminal Velocity

A terminal-based API client built for the [Boot.dev Hackathon 2025](https://blog.boot.dev/news/hackathon-2025/).

## Features

- Terminal user interface with beautiful TUI
- Request collections and organization  
- Demo data generation with realistic APIs
- Request builder with tabs for body, headers, query params
- Production-ready logging system

## Tech Stack

The project uses:

1. **Go** for core logic and HTTP operations
2. **Bubble Tea** for terminal user interface
3. **SQLite** for file-based storage
4. **SQLC** for type-safe database operations
5. **Goose** for database migrations

## Installation

```bash
go install github.com/maniac-en/req@v0.1.0
req
```

## What's Implemented

- Collections CRUD operations (create, edit, delete, navigate)
- Request builder interface with tabbed editing
- Endpoint browsing with sidebar navigation
- Demo data generation (JSONPlaceholder, ReqRes, HTTPBin APIs)
- Beautiful warm color scheme with vim-like navigation
- Pagination and real-time search filtering

## Coming Soon

- HTTP request execution (core feature)
- Response viewer with syntax highlighting  
- Endpoint management (add/edit endpoints)
- Environment variables support
- Export/import functionality

## Try It Out

**GitHub**: https://github.com/maniac-en/req  
**Installation**: `go install github.com/maniac-en/req@v0.1.0`  
**Usage**: Just run `req` in your terminal!

The app works completely offline with no external dependencies required.

---

This blog is built with ❤️ using [pyssg](https://github.com/maniac-en/pyssg) - A guided learning project at [boot.dev](https://www.boot.dev/courses/build-static-site-generator)