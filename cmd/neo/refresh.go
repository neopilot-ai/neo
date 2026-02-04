package main

import (
	"strings"

	"github.com/neopilot-ai/neo/cmd/neo/cli"
	"github.com/neopilot-ai/neo/cmd/neo/mosaic/ui"
	"github.com/neopilot-ai/neo/pkg/bus"
	"github.com/neopilot-ai/neo/pkg/project"
	"github.com/neopilot-ai/neo/pkg/server"
	"golang.org/x/sync/errgroup"
)

func CmdRefresh(c *cli.Cli) error {
	p, err := c.InitProject()
	if err != nil {
		return err
	}
	defer p.Cleanup()

	target := []string{}
	if c.String("target") != "" {
		target = strings.Split(c.String("target"), ",")
	}

	var wg errgroup.Group
	defer wg.Wait()
	ui := ui.New(c.Context)
	events := bus.SubscribeAll()
	defer close(events)
	wg.Go(func() error {
		for evt := range events {
			ui.Event(evt)
		}
		return nil
	})
	s, err := server.New()
	if err != nil {
		return err
	}
	wg.Go(func() error {
		defer c.Cancel()
		return s.Start(c.Context, p)
	})
	defer ui.Destroy()
	defer c.Cancel()
	err = p.Run(c.Context, &project.StackInput{
		Command:    "refresh",
		Target:     target,
		ServerPort: s.Port,
		Verbose:    c.Bool("verbose"),
	})
	if err != nil {
		return err
	}
	return nil
}
