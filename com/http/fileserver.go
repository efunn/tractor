package http

import (
	"net/http"

	"github.com/manifold/tractor/pkg/manifold"
)

func init() {
	manifold.RegisterComponent(&FileServer{}, "")
}

type FileServer struct {
	FileSystem http.FileSystem
}

func (c *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: how to determine base path
	// http.StripPrefix("/files/", http.FileServer(c.FileSystem)).ServeHTTP(w, r)
	http.FileServer(c.FileSystem).ServeHTTP(w, r)
}
