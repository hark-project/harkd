package server

import (
	"fmt"
	"net/http"

	"harkd/context"
	"harkd/routes"
)

// HarkdServer is a HTTP server providing hark's backend functionality.
type HarkdServer interface {
	Run() error
}

// Config is the config for a Hark server.
type Config struct {
	Port int `default:"8080"`
}

func (c Config) listenAddr() string {
	return fmt.Sprintf(":%d", c.Port)
}

// New constructs a new instance of HarkdServer.
func New(config Config, ctxFactory context.Factory) (HarkdServer, error) {
	router, err := routes.New(ctxFactory)
	if err != nil {
		return nil, err
	}
	return harkdServer{config, router}, nil
}

type harkdServer struct {
	Config
	routes.Router
}

// Run runs the server.
func (hds harkdServer) Run() error {
	listenAddr := hds.Config.listenAddr()
	fmt.Printf("harkd: listening on %s\n", listenAddr)
	return http.ListenAndServe(listenAddr, hds.Router)
}
