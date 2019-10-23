package http

import (
	"net/http"

	"github.com/manifold/tractor/pkg/manifold"
	"github.com/urfave/negroni"
)

func init() {
	manifold.RegisterComponent(&Logger{}, "")
}

type Logger struct {
	logger *negroni.Logger
}

func (c *Logger) Initialize() {
	c.logger = negroni.NewLogger()
}

func (c *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	c.logger.ServeHTTP(w, r, next)
}
