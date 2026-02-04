package main

import (
	"fmt"
	"runtime"

	"github.com/pulumi/pulumi/sdk/v3"
	"github.com/neopilot-ai/neo/v3/cmd/neo/cli"
	"github.com/neopilot-ai/neo/v3/pkg/global"
)

var CmdVersion = &cli.Command{
	Name: "version",
	Description: cli.Description{
		Short: "Print the version of the CLI",
		Long:  `Prints the current version of the CLI.`,
	},
	Run: func(cli *cli.Cli) error {
		fmt.Println("neo", version)
		if cli.Bool("verbose") {
			fmt.Println("pulumi", sdk.Version)
			fmt.Println("config", global.ConfigDir())
			fmt.Println("GOARCH", runtime.GOARCH)
			fmt.Println("GOOS", runtime.GOOS)
		}
		return nil
	},
}
