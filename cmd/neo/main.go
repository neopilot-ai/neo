package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/neopilot-ai/neo/cmd/neo/cli"
	"github.com/neopilot-ai/neo/cmd/neo/mosaic/errors"
	"github.com/neopilot-ai/neo/cmd/neo/mosaic/ui"
	"github.com/neopilot-ai/neo/internal/util"
	"github.com/neopilot-ai/neo/pkg/flag"
	"github.com/neopilot-ai/neo/pkg/global"
	"github.com/neopilot-ai/neo/pkg/process"
	"github.com/neopilot-ai/neo/pkg/project"
	"github.com/neopilot-ai/neo/pkg/telemetry"
)

var version = "dev"

func main() {
	// check if node_modules/.bin/neo exists
	nodeModulesBinPath := filepath.Join("node_modules", ".bin", "neo")
	binary, _ := os.Executable()
	if _, err := os.Stat(nodeModulesBinPath); err == nil && !strings.Contains(binary, "node_modules") && os.Getenv("NEO_SKIP_LOCAL") != "true" && version != "dev" {
		// forward command to node_modules/.bin/neo
		fmt.Println(ui.TEXT_WARNING_BOLD.Render("Warning: ") + "You are using a global installation of NEO but you also have a local installation specified in your package.json. The local installation will be used but you should typically run it through your package manager.")
		cmd := process.Command(nodeModulesBinPath, os.Args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "NEO_SKIP_LOCAL=true")
		if err := cmd.Run(); err != nil {
			os.Exit(1)
		}
		return
	}
	telemetry.SetVersion(version)
	defer telemetry.Close()
	defer process.Cleanup()
	telemetry.Track("cli.start", map[string]interface{}{
		"args": os.Args[1:],
	})
	err := run()
	if err != nil {
		err := errors.Transform(err)
		errorMessage := err.Error()
		truncated := errorMessage
		if len(errorMessage) > 255 {
			truncated = errorMessage[:255]
		}
		telemetry.Track("cli.error", map[string]interface{}{
			"error": truncated,
		})
		if readableErr, ok := err.(*util.ReadableError); ok {
			slog.Error("exited with error", "err", readableErr.Unwrap())
			msg := readableErr.Error()
			if msg != "" {
				ui.Error(readableErr.Error())
				if readableErr.IsHinted() {
					fmt.Println("   " + ui.TEXT_DIM.Render(readableErr.Unwrap().Error()))
				}
			}
		} else {
			slog.Error("exited with error", "err", err)
			// check if context cancelled error
			if err != context.Canceled {
				ui.Error("Unexpected error occurred. Please run with --print-logs or check .neo/log/neo.log if available.")
			}
		}
		telemetry.Close()
		os.Exit(1)
		return
	}
	telemetry.Track("cli.success", map[string]interface{}{})
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-interruptChannel
		slog.Info("interrupted")
		cancel()
	}()
	c, err := cli.New(ctx, cancel, root, version)
	if err != nil {
		return err
	}
	_, err = user.Current()
	if err != nil {
		return err
	}

	if !flag.NEO_SKIP_DEPENDENCY_CHECK {
		spin := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		spin.Suffix = "  Download dependencies..."
		if global.NeedsPulumi() {
			spin.Suffix = "  Installing pulumi..."
			spin.Start()
			err := global.InstallPulumi(ctx)
			if err != nil {
				spin.Stop()
				return util.NewHintedError(err, "Could not install pulumi")
			}
		}
		if global.NeedsBun() {
			spin.Suffix = "  Installing bun..."
			spin.Start()
			err := global.InstallBun(ctx)
			if err != nil {
				spin.Stop()
				return util.NewHintedError(err, "Could not install bun")
			}
		}
		spin.Stop()
	}
	return c.Run()
}

var root = &cli.Command{
	Name: "neo",
	Description: cli.Description{
		Short: "deploy anything",
		Long: strings.Join([]string{
			"The CLI helps you manage your NEO apps.",
			"",
			"If you are using NEO as a part of your Node project, we recommend installing it locally.",
			"```bash",
			"npm install neo",
			"```",
			"---",
			"If you are not using Node, you can install the CLI globally.",
			"```bash",
			"curl -fsSL https://neo.dev/install | bash",
			"```",
			"",
			":::note",
			"The CLI currently supports macOS, Linux, and WSL. Windows support is in beta.",
			":::",
			"",
			"To install a specific version.",
			"",
			"```bash \"VERSION=0.0.403\"",
			"curl -fsSL https://neo.dev/install | VERSION=0.0.403 bash",
			"```",
			"---",
			"#### With a package manager",
			"",
			"You can also use a package manager to install the CLI.",
			"",
			"- **macOS**",
			"",
			"  The CLI is available via a Homebrew Tap, and as downloadable binary in the [releases](https://github.com/neopilot-ai/neo/releases/latest).",
			"",
			"  ```bash",
			"  brew tap neo/tap",
			"  brew install neo",
			"",
			"  # Upgrade",
			"  brew upgrade neo",
			"  ```",
			"",
			"  You might have to run `brew upgrade neo`, before the update.",
			"",
			"- **Linux**",
			"",
			"  The CLI is available as downloadable binaries in the [releases](https://github.com/neopilot-ai/neo/releases/latest). Download the `.deb` or `.rpm` and install with `sudo dpkg -i` and `sudo rpm -i`.",
			"",
			"  For Arch Linux, it's available in the [aur](https://aur.archlinux.org/packages/neo-bin).",
			"---",
			"#### Usage",
			"",
			"Once installed you can run the commands using.",
			"",
			"```bash",
			"neo [command]",
			"```",
			"",
			"The CLI takes a few global flags. For example, the deploy command takes the `--stage` flag",
			"",
			"```bash",
			"neo deploy --stage production",
			"```",
			"---",
			"#### Environment variables",
			"",
			"You can access any environment variables set in the CLI in your `neo.config.ts` file. For example, running:",
			"",
			"```bash",
			"ENV_VAR=123 neo deploy",
			"```",
			"",
			"Will let you access `ENV_VAR` through `process.env.ENV_VAR`.",
		}, "\n"),
	},
	Flags: []cli.Flag{
		{
			Name: "stage",
			Type: "string",
			Description: cli.Description{
				Short: "The stage to deploy to",
				Long: strings.Join([]string{
					"Set the stage the CLI is running on.",
					"",
					"```bash frame=\"none\"",
					"neo [command] --stage production",
					"```",
					"",
					"The stage is a string that is used to prefix the resources in your app. This allows you to have multiple _environments_ of your app running in the same account.",
					"",
					":::tip",
					"Changing the stage will redeploy your app to a new stage with new resources. The old resources will still be around in the old stage.",
					":::",
					"",
					"You can also use the `NEO_STAGE` environment variable.",
					"```bash frame=\"none\"",
					"NEO_STAGE=dev neo [command]",
					"```",
					"This can also be declared in a `.env` file or in the CLI session.",
					"",
					"If the stage is not passed in, then the CLI will:",
					"",
					"1. Use the username on the local machine.",
					"   - If the username is `root`, `admin`, `prod`, `dev`, `production`, then it will prompt for a stage name.",
					"2. Store this in the `.neo/stage` file and reads from it in the future.",
					"",
					"This stored stage is called your **personal stage**.",
				}, "\n"),
			},
		},
		{
			Name: "verbose",
			Type: "bool",
			Description: cli.Description{
				Short: "Enable verbose logging",
				Long: strings.Join([]string{
					"",
					"Prints extra information to the log files in the `.neo/` directory.",
					"",
					"```bash",
					"neo [command] --verbose",
					"```",
					"",
					"To also view this on the screen, use the `--print-logs` flag.",
					"",
				}, "\n"),
			},
		},
		{
			Name: "print-logs",
			Type: "bool",
			Description: cli.Description{
				Short: "Print logs to stderr",
				Long: strings.Join([]string{
					"",
					"Print the logs to the screen. These are logs that are written to the `.neo/` directory.",
					"",
					"```bash",
					"neo [command] --print-logs",
					"```",
					"It can also be set using the `NEO_PRINT_LOGS` environment variable.",
					"",
					"```bash",
					"NEO_PRINT_LOGS=1 neo [command]",
					"```",
					"This is useful when running in a CI environment.",
					"",
				}, "\n"),
			},
		},
		{
			Name: "config",
			Type: "string",
			Description: cli.Description{
				Short: "Path to the config file",
				Long: strings.Join([]string{
					"",
					"Optionally, pass in a path to the NEO config file. This default to",
					"`neo.config.ts` in the current directory.",
					"",
					"```bash",
					"neo --config path/to/config.ts [command]",
					"```",
					"",
					"This is useful when your monorepo has multiple NEO apps in it.",
					"You can run the NEO CLI for a specific app by passing in the path to",
					"its config file.",
				}, "\n"),
			},
		},
		{
			Name: "help",
			Type: "bool",
			Description: cli.Description{
				Short: "Print help",
				Long: strings.Join([]string{
					"Prints help for the given command.",
					"",
					"```sh frame=\"none\"",
					"neo [command] --help",
					"```",
					"",
					"Or the global help.",
					"",
					"```sh frame=\"none\"",
					"neo --help",
					"```",
				}, "\n"),
			},
		},
	},
	Children: []*cli.Command{
		{
			Name: "init",
			Description: cli.Description{
				Short: "Initialize a new project",
				Long: strings.Join([]string{
					"Initialize a new project in the current directory. This will create a `neo.config.ts` and `neo install` your providers.",
					"",
					"If this is run in a Next.js, Remix, Astro, or SvelteKit project, it'll init NEO in drop-in mode.",
					"",
					"To skip the interactive confirmation after detecting the framework.",
					"",
					"```bash frame=\"none\"",
					"neo init --yes",
					"```",
				}, "\n"),
			},
			Run: CmdInit,
			Flags: []cli.Flag{
				{
					Name: "yes",
					Type: "bool",
					Description: cli.Description{
						Short: "Skip interactive confirmation",
						Long:  "Skip interactive confirmation for detected framework.",
					},
				},
			},
		},
		{
			Name:   "ui",
			Hidden: true,
			Run:    CmdUI,
			Flags: []cli.Flag{
				{
					Name: "filter",
					Type: "string",
					Description: cli.Description{
						Short: "Filter events",
						Long:  "Filter events.",
					},
				},
			},
		},
		{
			Name: "dev",
			Description: cli.Description{
				Short: "Run in development mode",
				Long: strings.Join([]string{
					"Run your app in dev mode. By default, this starts a multiplexer with processes that",
					" deploy your app, run your functions, and start your frontend.",
					"",
					":::note",
					"The tabbed terminal UI is only available on Linux/macOS and WSL.",
					":::",
					"",
					"Each process is run in a separate tab that you can click on in the sidebar.",
					"",
					"![neo dev multiplexer mode](../../../../assets/docs/cli/neo-dev-multiplexer-mode.png)",
					"",
					"The multiplexer makes it so that you won't have to start your frontend or",
					"your container applications separately.",
					"",
					"<VideoAside title=\"Watch a video about dev mode\" href=\"https://youtu.be/mefLc137EB0\" />",
					"",
					"Here's what happens when you run `neo dev`.",
					"",
					"- Deploy most of your resources as-is.",
					"- Except for components that have a `dev` prop.",
					"  - `Function` components are run [_Live_](/docs/live/) in the **Functions** tab.",
					"  - `Task` components have their _stub_ versions deployed that proxy the task",
					"    and run their `dev.command` in the **Tasks** tab.",
					"  - Frontends like `Nextjs`, `Remix`, `Astro`, `StaticSite`, etc. have their dev",
					"    servers started in a separate tab and are not deployed.",
					"  - `Service` components are not deployed, and instead their `dev.command` is",
					"    started in a separate tab.",
					"  - `Postgres`, `Aurora`, and `Redis` link to a local database if the `dev` prop is",
					"    set.",
					"- Start an [`neo tunnel`](#tunnel) session in a new tab if your app has a `Vpc`"",
					"  with `bastion` enabled.",
					"- Load any [linked resources](/docs/linking) in the environment.",
					"- Start a watcher for your `neo.config.ts` and redeploy any changes.",
					"",
					":::note",
					"The `Service` component and the frontends like `Nextjs` or `StaticSite` are not",
					"deployed by `neo dev`.",
					":::",
					"",
					"Optionally, you can disable the multiplexer and not spawn any child",
					"processes by running `neo dev` in basic mode.",
					"",
					"```bash frame=\"none\"",
					"neo dev --mode=basic",
					"```",
					"",
					"This will only deploy your app and run your functions. If you are coming from NEO",
					"v2, this is how `neo dev` used to work.",
					"",
					"However in `basic` mode, you'll need to start your frontend separately by running",
					"`neo dev` in a separate terminal session by passing in the command. For example:",
					"",
					"```bash frame=\"none\"",
					"neo dev next dev",
					"```",
					"",
					"By wrapping your command, it'll load your [linked resources](/docs/linking) in the",
					"environment.",
					"",
					"To pass in a flag to the command, use `--`.",
					"",
					"```bash frame=\"none\"",
					"neo dev -- next dev --turbo",
					"```",
					"",
					"You can also disable the tabbed terminal UI by running `neo dev` in",
					"mono mode.",
					"",
					"```bash frame=\"none\"",
					"neo dev --mode=mono",
					"```",
					"",
					"Unlike `basic` mode, this'll spawn child processes. But instead of",
					"a tabbed UI it'll show their outputs in a single stream.",
					"",
					"This is used by default in Windows.",
				}, "\n"),
			},
			Flags: []cli.Flag{
				{
					Name: "mode",
					Type: "string",
					Description: cli.Description{
						Short: "mode=mono to turn off multiplexer. mode=basic to not spawn any child processes",
						Long:  "Defaults to using `multi` mode. Use `mono` to get a single stream of all child process logs or `basic` to not spawn any child processes.",
					},
				},
			},
			Args: []cli.Argument{
				{
					Name: "command",
					Description: cli.Description{
						Short: "The command to run",
					},
				},
			},
			Examples: []cli.Example{
				{
					Content: "neo dev",
					Description: cli.Description{
						Short: "Brings up your entire app - should be all you need",
					},
				},
				{
					Content: "neo dev next dev",
					Description: cli.Description{
						Short: "Start a command connected to a running neo dev session",
					},
				},
				{
					Content: "neo dev -- next dev --turbo",
					Description: cli.Description{
						Short: "Use -- to pass flags to the command",
					},
				},
			},
			Run: CmdMosaic,
		},
		CmdDeploy,
		CmdDiff,
		{
			Name: "add",
			Description: cli.Description{
				Short: "Add a new provider",
				Long: strings.Join([]string{
					"Adds and installs the given provider. For example,",
					"",
					"```bash frame=\"none\"",
					"neo add aws",
					"```",
					"",
					"This command will:",
					"",
					"1. Installs the package for the AWS provider.",
					"2. Add `aws` to the globals in your `neo.config.ts`.",
					"3. And, add it to your `providers`.",
					"",
					"```ts title=\"neo.config.ts\"",
					"{",
					"  providers: {",
					"    aws: \"6.27.0\"",
					"  }",
					"}",
					"```",
					"",
					"You can use any provider listed in the [Directory](/docs/all-providers#directory).",
					"",
					":::note",
					"Running `neo add aws` above is the same as manually adding the provider to your config and running `neo install`.",
					":::",
					"",
					"By default, the latest version of the provider is installed. If you want to use a specific version, you can change it in your config.",
					"",
					"```ts title=\"neo.config.ts\"",
					"{",
					"  providers: {",
					"    aws: {",
					"      version: \"6.26.0\"",
					"    }",
					"  }",
					"}",
					"```",
					"",
					"You'll need to run `neo install` if you update the `providers` in your config.",
					"",
					"By default, these packages are fetched from the NPM registry. If you want to use a different registry, you can set the `NPM_REGISTRY` environment variable.",
					"",
					"```bash frame=\"none\"",
					"NPM_REGISTRY=https://my-registry.com neo add aws",
					"```",
				}, "\n"),
			},
			Args: []cli.Argument{
				{
					Name:     "provider",
					Required: true,
					Description: cli.Description{
						Short: "The provider to add",
						Long:  "The provider to add.",
					},
				},
			},
			Run: func(cli *cli.Cli) error {
				pkg := cli.Positional(0)
				spin := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
				spin.Suffix = "  Adding provider..."
				spin.Start()
				defer spin.Stop()
				cfgPath, err := cli.Discover()
				if err != nil {
					return err
				}
				stage, err := cli.Stage(cfgPath)
				if err != nil {
					return err
				}
				p, err := project.New(&project.ProjectConfig{
					Version: version,
					Config:  cfgPath,
					Stage:   stage,
				})
				if err != nil {
					return err
				}
				if !p.CheckPlatform(version) {
					err := p.CopyPlatform(version)
					if err != nil {
						return err
					}
				}
				entry, err := project.FindProvider(pkg, "latest")
				if err != nil {
					return util.NewReadableError(err, "Could not find provider "+pkg)
				}
				err = p.Add(entry.Name, entry.Version)
				if err != nil {
					return err
				}
				spin.Suffix = "  Downloading provider..."
				p, err = project.New(&project.ProjectConfig{
					Version: version,
					Config:  cfgPath,
					Stage:   stage,
				})
				if err != nil {
					return err
				}
				err = p.Install()
				if err != nil {
					return err
				}
				spin.Stop()
				ui.Success(fmt.Sprintf("Added provider \"%s\". You can create resources with `new %s.SomeResource()`.", entry.Alias, entry.Alias))
				return nil
			},
		},
		{
			Name: "install",
			Description: cli.Description{
				Short: "Install all the providers",
				Long: strings.Join([]string{
					"Installs the providers in your `neo.config.ts`. You'll need this command when:",
					"",
					"1. You add a new provider to the `providers` or `home` in your config.",
					"2. Or, when you want to install new providers after you `git pull` some changes.",
					"",
					":::tip",
					"The `neo install` command is similar to `npm install`.",
					":::",
					"",
					"Behind the scenes, it installs the packages for your providers and adds the providers to your globals.",
					"",
					"If you don't have a version specified for your providers in your `neo.config.ts`, it'll install their latest versions.",
				}, "\n"),
			},
			Run: func(cli *cli.Cli) error {
				cfgPath, err := cli.Discover()
				if err != nil {
					return err
				}

				stage, err := cli.Stage(cfgPath)
				if err != nil {
					return err
				}

				p, err := project.New(&project.ProjectConfig{
					Version: version,
					Config:  cfgPath,
					Stage:   stage,
				})
				if err != nil {
					return err
				}

				spin := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
				defer spin.Stop()
				spin.Suffix = "  Installing providers..."
				spin.Start()
				if !p.CheckPlatform(version) {
					err := p.CopyPlatform(version)
					if err != nil {
						return err
					}
				}

				err = p.Install()
				if err != nil {
					return err
				}
				spin.Stop()
				ui.Success("Installed providers")
				return nil
			},
		},
		{
			Name: "secret",
			Description: cli.Description{
				Short: "Manage secrets",
				Long: strings.Join([]string{
					"Manage the secrets in your app defined with `neo.Secret`.",
					"",
					"<VideoAside title=\"Watch a video on how secrets work\" href=\"https://youtu.be/7tW2L3P6LKw\" />",
					"",
					"The `--fallback` flag can be used to manage the fallback values of a secret.",
					"",
					"Applies to all the sub-commands in `neo secret`.",
				}, "\n"),
			},
			Flags: []cli.Flag{
				{
					Name: "fallback",
					Type: "bool",
					Description: cli.Description{
						Short: "Manage the fallback values of secrets",
						Long:  "Manage the fallback values of secrets.",
					},
				},
			},
			Children: []*cli.Command{
				CmdSecretSet,
				CmdSecretRemove,
				CmdSecretLoad,
				CmdSecretList,
			},
		},
		{
			Name: "shell",
			Args: []cli.Argument{
				{
					Name: "command",
					Description: cli.Description{
						Short: "A command to run",
						Long:  "A command to run.",
					},
				},
			},
			Flags: []cli.Flag{
				{
					Name: "target",
					Description: cli.Description{
						Short: "Run it only for a component",
						Long:  "Only run it for the given component.",
					},
				},
			},
			Description: cli.Description{
				Short: "Run a command with linked resources",
				Long: strings.Join([]string{
					"Run a command in the environment of your app with all the linked resources loaded.",
					"",
					"For example, let's say you have the following resources in your app.",
					"",
					"```js title=\"neo.config.ts\" {5,9}",
					"new neo.aws.Bucket(\"MyMainBucket\");",
					"new neo.aws.Bucket(\"MyAdminBucket\");",
					"```",
					"",
					"We can now write a script that'll can access both these resources with the [JS SDK](/docs/reference/sdk/).",
					"",
					"```js title=\"my-script.js\" \"Resource.MyMainBucket.name\" \"Resource.MyAdminBucket.name\"",
					"import { Resource } from \"neo\";",
					"",
					"console.log(Resource.MyMainBucket.name, Resource.MyAdminBucket.name);",
					"```",
					"",
					"And run the script with `neo shell`.",
					"",
					"```bash frame=\"none\" frame=\"none\"",
					"neo shell node my-script.js",
					"```",
					"",
					"This'll have access to all the buckets from above.",
					"",
					":::tip",
					"Run the command with `--` to pass arguments to it.",
					":::",
					"",
					"To pass arguments into the script, you'll need to prefix it using `--`.",
					"",
					"```bash frame=\"none\" frame=\"none\" /--(?!a)/",
					"neo shell -- node my-script.js --arg1 --arg2",
					"```",
					"",
					"If no command is passed in, it opens a shell session with the linked resources.",
					"",
					"```bash frame=\"none\" frame=\"none\"",
					"neo shell",
					"```",
					"",
					"This is useful if you want to run multiple commands, all while accessing the resources in your app.",
					"",
					"Optionally, you can run this for a specific component by passing in the name of the component.",
					"",
					"```bash frame=\"none\" frame=\"none\"",
					"neo shell --target MyComponent",
					"```",
					"",
					"Here the linked resources for `MyComponent` and its environment variables are available.",
				}, "\n"),
			},
			Examples: []cli.Example{
				{
					Content: "neo shell",
					Description: cli.Description{
						Short: "Open a shell session",
					},
				},
			},
			Run: CmdShell,
		},
		{
			Name: "remove",
			Description: cli.Description{
				Short: "Remove your application",
				Long: strings.Join([]string{
					"Removes your application. By default, it removes your personal stage.",
					"",
					":::tip",
					"The resources in your app are removed based on the `removal` setting in your `neo.config.ts`.",
					":::",
					"",
					"This does not remove the NEO _state_ and _bootstrap_ resources in your account as these might still be in use by other apps. You can remove them manually if you want to reset your account, [learn more](/docs/state/#reset).",
					"",
					"Optionally, remove your app from a specific stage.",
					"",
					"```bash frame=\"none\" frame=\"none\"",
					"neo remove --stage production",
					"```",
					"You can also remove a specific component by passing in the name of the component from your `neo.config.ts`.",
					"",
					"```bash frame=\"none\"",
					"neo remove --target MyComponent",
					"```",
				}, "\n"),
			},
			Flags: []cli.Flag{
				{
					Name: "target",
					Type: "string",
					Description: cli.Description{
						Short: "Run it only for a component",
						Long:  "Only run it for the given component.",
					},
				},
			},
			Run: CmdRemove,
		},
		{
			Name: "unlock",
			Description: cli.Description{
				Short: "Clear any locks on the app state",
				Long: strings.Join([]string{
					"When you run `neo deploy`, it acquires a lock on your state file to prevent concurrent deploys.",
					"",
					"However, if something unexpectedly kills the `neo deploy` process, or if you manage to run `neo deploy` concurrently, the lock might not be released.",
					"",
					"This should not usually happen, but it can prevent you from deploying. You can run `neo unlock` to release the lock.",
				}, "\n"),
			},
			Run: func(c *cli.Cli) error {
				p, err := c.InitProject()
				if err != nil {
					return err
				}
				defer p.Cleanup()

				err = p.ForceUnlock()
				if err != nil {
					return err
				}
				color.New(color.FgGreen, color.Bold).Print("✓ ")
				color.New(color.FgWhite).Print(" Unlocked the app state for: ")
				color.New(color.FgWhite, color.Bold).Println(p.App().Name, "/", p.App().Stage)
				return nil
			},
		},
		CmdVersion,
		{
			Name: "upgrade",
			Description: cli.Description{
				Short: "Upgrade the CLI",
				Long: strings.Join([]string{
					"Upgrade the CLI to the latest version. Or optionally, pass in a version to upgrade to.",
					"",
					"```bash frame=\"none\"",
					"neo upgrade 0.10",
					"```",
				}, "\n"),
			},
			Args: cli.ArgumentList{
				{
					Name: "version",
					Description: cli.Description{
						Short: "A version to upgrade to",
						Long:  "A version to upgrade to.",
					},
				},
			},
			Run: CmdUpgrade,
		},
		{
			Name: "telemetry", Description: cli.Description{
				Short: "Manage telemetry settings",
				Long: strings.Join([]string{
					"Manage telemetry settings.",
					"",
					"NEO collects completely anonymous telemetry data about general usage. We track:",
					"- Version of NEO in use",
					"- Command invoked, `neo dev`, `neo deploy`, etc.",
					"- General machine information, like the number of CPUs, OS, CI/CD environment, etc.",
					"",
					"This is completely optional and can be disabled at any time.",
					"",
					"You can also opt-out by setting an environment variable: `NEO_TELEMETRY_DISABLED=1` or `DO_NOT_TRACK=1`.",
				}, "\n"),
			},
			Children: []*cli.Command{
				{
					Name: "enable",
					Description: cli.Description{
						Short: "Enable telemetry",
						Long:  "Enable telemetry.",
					},
					Run: func(cli *cli.Cli) error {
						return telemetry.Enable()
					},
				},
				{
					Name: "disable",
					Description: cli.Description{
						Short: "Disable telemetry",
						Long:  "Disable telemetry.",
					},
					Run: func(cli *cli.Cli) error {
						return telemetry.Disable()
					},
				},
			},
		},
		{
			Name:   "introspect",
			Hidden: true,
			Run: func(cli *cli.Cli) error {
				data, err := json.MarshalIndent(cli.Path()[0], "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
				return nil
			},
		},
		{
			Name:   "print-and-not-quit",
			Hidden: true,
			Run: func(cli *cli.Cli) error {
				lines := strings.Split(os.Getenv("NEO_DEV_COMMAND_MESSAGE"), "\n")
				for _, line := range lines {
					fmt.Println(line)
				}
				<-cli.Context.Done()
				return nil
			},
		},
		{
			Name:   "common-errors",
			Hidden: true,
			Run: func(cli *cli.Cli) error {
				data, err := json.MarshalIndent(project.CommonErrors, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
				return nil
			},
		},
		{
			Name: "refresh",
			Description: cli.Description{
				Short: "Refresh the local app state",
				Long: strings.Join([]string{
					"Compares your local state with the state of the resources in the cloud provider. Any changes that are found are adopted into your local state. It will:",
					"",
					"1. Go through every single resource in your state.",
					"2. Make a call to the cloud provider to check the resource.",
					"   - If the configs are different, it'll **update the state** to reflect the change.",
					"   - If the resource doesn't exist anymore, it'll **remove it from the state**.",
					"",
					":::note",
					"The `neo refresh` does not make changes to the resources in the cloud provider.",
					":::",
					"You can also refresh a specific component by passing in the name of the component.",
					"",
					"```bash frame=\"none\"",
					"neo refresh --target MyComponent",
					"```",
					"",
					"This is useful for cases where you want to ensure that your local state is in sync with your cloud provider. [Learn more about how state works](/docs/providers/#how-state-works).",
				}, "\n"),
			},
			Flags: []cli.Flag{
				{
					Name: "target",
					Type: "string",
					Description: cli.Description{
						Short: "Run it only for a component",
						Long:  "Only run it for the given component.",
					},
				},
			},
			Run: CmdRefresh,
		},
		CmdState,
		CmdCert,
		CmdTunnel,
		CmdDiagnostic,
	},
}
