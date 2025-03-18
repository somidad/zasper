//go:build !apiserver
// +build !apiserver

package main

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"log"
	"strings"
)

//go:embed ui/build/*
var staticFiles embed.FS

type spaHandler struct {
	staticFS   embed.FS
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Treat the request path as relative and join it with the static directory
    path := filepath.Join(h.staticPath, r.URL.Path)
       log.Println(path)

    // Normalize the path (replace backslashes with forward slashes)
    path = strings.ReplaceAll(path, "\\", "/")

    // Open the file
    _, err := h.staticFS.Open(path)
    if os.IsNotExist(err) {
        // File does not exist, serve index.html
        index, err := h.staticFS.ReadFile(filepath.Join(h.staticPath, h.indexPath))
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        w.WriteHeader(http.StatusAccepted)
        w.Write(index)
        return
    } else if err != nil {
        // If we got an error (other than file not existing), return 500
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // File exists, serve it
    statics, err := fs.Sub(h.staticFS, h.staticPath)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.FileServer(http.FS(statics)).ServeHTTP(w, r)
}

func getSpaHandler() http.Handler {
	return spaHandler{staticFS: staticFiles, staticPath: "ui/build", indexPath: "index.html"}
}
