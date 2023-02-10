package main

import (
	"fmt"
	"github.com/metacatdud/netrau/internal/hub"
	"gitlab.com/macroscope-lab/atomika/log"
	"os"

	"github.com/metacatdud/netrau/internal/member"
	"gitlab.com/macroscope-lab/atomika"
	"gitlab.com/macroscope-lab/atomika/transport"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	log.NewWithAutoConf()
	trans := transport.NewWithAutoConf()

	// Runtime section - Attempt to load configurations from environment
	// Defaults, config files, env vars
	var (
		memberCfg = &member.Runtime{}
	)

	if err := memberCfg.Configure("netrau.members"); err != nil {
		return err
	}

	// Setup services
	hubService, err := hub.New(
		hub.WithResendLimit(memberCfg.ResendLimit),
		hub.WithLocalAddr(memberCfg.LocalAddr),
		hub.WithHubAddr(memberCfg.Join),
	)
	if err != nil {
		return err
	}

	// Boot
	app := atomika.New()
	app.RegisterService([]atomika.Service{
		trans,
		hubService,
	})

	if err = app.Boot(); err != nil {
		return err
	}

	return nil
}
