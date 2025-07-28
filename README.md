# Req - Test APIs with Terminal Velocity

[![tests](https://github.com/maniac-en/req/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/maniac-en/req/actions/workflows/go.yml)
![GitHub repo](https://img.shields.io/badge/built%20at-Boot.dev%20Hackathon-blueviolet)

## About

`req` is a lightweight, terminal-based API client built for the
[Boot.dev Hackathon 2025](https://github.com/maniac-en/req?tab=License-1-ov-file).
It features a fast and minimal text user interface and lets you create, send,
and inspect HTTP requests interactively from the command line. It is ideal for
testing APIs without leaving your terminal.

Read more about `req` over here -
[Announcement Blog](https://maniac-en.github.io/req/)

## Installation

### You can install `req` using `go install`:

To install a specific release

```
go install github.com/maniac-en/req@latest
```

Replace `latest` with the specific version you want.

### Requirements

- Go version 1.24.4

## Usage

After installing `req`, you can run it using this command.

```
req
```

## Libraries Used

### Terminal UI (by Charm.sh)

- [bubbletea](https://github.com/charmbracelet/bubbletea) — A powerful, fun TUI
  framework for Go
- [bubbles](https://github.com/charmbracelet/bubbles) — Pre-built components for
  TUI apps
- [lipgloss](https://github.com/charmbracelet/lipgloss) — Terminal style/layout
  DSL

## License

This project is licensed under the
[MIT License](https://github.com/maniac-en/req?tab=License-1-ov-file).

```
1. Mudassir Bilal (mailto:mughalmudassir966@gmail.com)
2. Shivam Mehta (mailto:sm.cse17@gmail.com)
3. Yash Ranjan (mailto:yash.ranjan25@gmail.com)

MIT License

Copyright (c) 2025 Mudassir Bilal, Shivam Mehta, Yash Ranjan

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
```
