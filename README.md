# URL Shortener

A simple URL shortener service built with Go and SQLite.
[![Live Demo](https://img.shields.io/badge/Demo-Live-green?style=for-the-badge)](https://your-deployed-link.com)

## Features
- Shorten long URLs into compact, shareable links.
- Redirect short URLs to their original destinations.
- Lightweight and easy to deploy.

## Prerequisites
- Go 1.20 or higher
- SQLite (embedded, no setup required)

## Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/urlshortener.git
   cd urlshortener
2. Install Dependencies:
    ```bash
    go mod download
3. Run the application:
    ```bash
    go run main.go
4. The server will start on http://localhost:8080