package main

import (
	"fmt"
	"os"

	"harkd/context"
	"harkd/server"

	"github.com/ceralena/envconf"
)

const exitStatusFail = 1

func loadConfig() (*server.Config, error) {
	cfg := new(server.Config)
	err := envconf.ReadConfigEnv(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err)
	os.Exit(exitStatusFail)
}

func main() {
	serverCfg, err := loadConfig()
	if err != nil {
		fail(err)
	}

	contextFactory, err := context.HomeDirFactory()
	if err != nil {
		fail(err)
	}

	s, err := server.New(*serverCfg, contextFactory)
	if err != nil {
		fail(err)
	}

	err = s.Run()
	if err != nil {
		fail(err)
	}

}
