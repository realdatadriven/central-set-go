package main

import (
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/realdatadriven/central-set-go/assets"
)

// NewFallbackFileServer creates a file server that tries the filesystem first,
// then falls back to embedded files. This allows dynamic uploads while serving
// static assets from either disk or the embedded filesystem.
func NewFallbackFileServer() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Clean(r.URL.Path)

		// Prevent directory traversal
		if filepath.IsAbs(path) || path == "" || path == "." {
			http.NotFound(w, r)
			return
		}

		// Try filesystem first (for dynamic/development files)
		if tryFilesystem(w, r, path) {
			return
		}

		// Fall back to embedded files
		tryEmbedded(w, r, path)
	})
}

// tryFilesystem attempts to serve from the filesystem
// Returns true if file was served successfully
func tryFilesystem(w http.ResponseWriter, r *http.Request, path string) bool {
	fsPath := filepath.Join("static", path)

	file, err := os.Open(fsPath)
	if err != nil {
		return false
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return false
	}

	// Handle directories
	if info.IsDir() {
		indexFile, err := os.Open(filepath.Join(fsPath, "index.html"))
		if err != nil {
			return false
		}
		defer indexFile.Close()
		indexInfo, err := indexFile.Stat()
		if err != nil {
			return false
		}
		http.ServeContent(w, r, "index.html", indexInfo.ModTime(), indexFile)
		return true
	}

	http.ServeContent(w, r, info.Name(), info.ModTime(), file)
	return true
}

func ServeFSFile(fsys fs.FS, name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f, err := fsys.Open(name)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer f.Close()

		http.ServeContent(w, r, name, time.Time{}, f.(io.ReadSeeker))
	}
}

// tryEmbedded attempts to serve from the embedded filesystem
func tryEmbedded(w http.ResponseWriter, r *http.Request, path string) {
	// Get the embedded static directory

	staticFS, err := fs.Sub(assets.EmbeddedFiles, "static")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	file, err := staticFS.Open(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Handle directories
	if info.IsDir() {
		/*indexFile, err := staticFS.Open(filepath.Join(path, "index.html"))
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer indexFile.Close()
		indexInfo, err := indexFile.Stat()
		if err != nil {
			http.NotFound(w, r)
			return
		}*/
		// http.ServeContent(w, r, "index.html", indexInfo.ModTime(), indexFile)
		ServeFSFile(staticFS, filepath.Join(path, "index.html"))(w, r)
		return
	}
	ServeFSFile(staticFS, path)(w, r)
	// http.ServeContent(w, r, info.Name(), info.ModTime(), file)
}

// ServeStaticFile serves a single file from static directory (filesystem first, then embedded)
func ServeStaticFile(w http.ResponseWriter, r *http.Request, filePath string) {
	path := filepath.Clean(filePath)

	// Try filesystem first
	if tryFilesystem(w, r, path) {
		return
	}

	// Fall back to embedded
	tryEmbedded(w, r, path)
}
