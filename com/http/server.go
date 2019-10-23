package http

import (
	"log"
	"net"
	"net/http"

	"github.com/manifold/tractor/pkg/manifold"
	frontend "github.com/manifold/tractor/pkg/session"
	"github.com/urfave/negroni"
)

func init() {
	manifold.RegisterComponent(&Server{}, "")
}

type Server struct {
	Listener net.Listener `com:"singleton"`
	Handler  http.Handler `com:"singleton"`
	// Middleware []negroni.Handler `com:"extpoint"`

	s *http.Server
}

func (c *Server) InspectorButtons() []frontend.Button {
	return []frontend.Button{{
		Name: "Serve",
	}}
}

func (c *Server) Serve() {
	log.Println("starting http server")
	n := negroni.New()
	// for _, handler := range c.Middleware {
	// 	n.Use(handler)
	// }
	n.UseHandler(c.Handler)
	c.s = &http.Server{
		Handler: n,
	}
	go func() {
		if err := c.s.Serve(c.Listener); err != nil {
			log.Fatal(err)
		}
	}()
}
