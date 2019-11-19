package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli"
	. "github.com/zhanbei/static-server/conf"
	"github.com/zhanbei/static-server/configs"
	"github.com/zhanbei/static-server/db"
	"github.com/zhanbei/static-server/helpers/terminator"
	"github.com/zhanbei/static-server/libs"
	"github.com/zhanbei/static-server/utils"
)

var ops = NewDefaultServerOptions()

var OptionConfiguresFile = ""

// The primary program entrance.
// (Cli Arguments Receiver + Configuration File Parser + MongoDB Driver)
// Support more custom built, like for lite/medium/heavy programs, for cli/gui(with different themes) modes, and for linux/windows/mac platforms.
// @see [Support multiple entrances and keep the current one as the primary. · Issue #6 · zhanbei/static-server](https://github.com/zhanbei/static-server/issues/6)
func main() {
	app := cli.NewApp()
	app.Name = "static-server"
	app.Usage = "A static server in Go, supporting hosting static files in the no-trailing-slash version."
	app.Version = "0.9.1"
	app.Description = "A static server in Go, supporting hosting static files in the no-trailing-slash version."
	app.UsageText = "static-server [global options] [<http-address>:]<http-port> <www-root-directory>"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "configure",

			Usage: "The configuration file to be used.",

			Destination: &OptionConfiguresFile,
		},
		cli.BoolFlag{
			Name: OptionNameEnableVirtualHosting,

			Usage: "Whether to enable virtual hosting; @see https://en.wikipedia.org/wiki/Virtual_hosting",

			Destination: &ops.UsingVirtualHost,
		},
		cli.BoolFlag{
			Name: OptionNameNoTrailingSlash,

			Usage: "Hosting static files in the " + OptionNameNoTrailingSlash + " mode.",

			Destination: &ops.NoTrailingSlash,
		},
		cli.BoolFlag{
			Name: OptionNameDirectoryListing,

			Usage: "Listing files of a directory if the index.html is not found when in the normal mode.",

			Destination: &ops.DirectoryListing,
		},
	}
	app.Action = Action

	err := app.Run(os.Args)
	if err != nil {
		terminator.ExitWithPreLaunchServerError(err, "Loading configures from environment variables failed!")
	}
}

func Action(c *cli.Context) error {
	if utils.NotEmpty(OptionConfiguresFile) {
		return ActionConfigurationFile(c, OptionConfiguresFile)
	} else {
		return ActionCliArguments(c, ops)
	}
}

// FIX-ME Use a default configuration file, like `vhss.(yaml|toml|json)`.
func ActionConfigurationFile(c *cli.Context, confFile string) error {
	cfg, err := configs.LoadServerConfigures(confFile)
	if err != nil {
		terminator.ExitWithConfigError(err, "Loading and validating the configures failed!")
	}
	err = cfg.ValidateFile()
	if err != nil {
		terminator.ExitWithPreLaunchServerError(err, "Validating the required resources following configures failed!")
	}
	bts, err := json.Marshal(cfg)
	fmt.Println("Loading configures:", string(bts))
	fmt.Println(cfg, cfg.Server, cfg.Loggers, cfg.MongoDbOptions, confFile)

	mon := cfg.MongoDbOptions
	if mon != nil && mon.Enabled {
		err = db.ConnectToMongoDb(cfg.MongoDbOptions)
		if err != nil {
			terminator.ExitWithPreLaunchServerError(err, "Connecting to mongodb failed!")
		}
	}

	return libs.RealServer(cfg)
}

func ActionCliArguments(c *cli.Context, ops *ServerOptions) error {
	ops.ValidateOrExit()

	if c.NArg() <= 0 {
		terminator.ExitWithConfigError(nil, "Please specify a port, like `static-server 8080`.")
	}
	address := c.Args().Get(0)
	address, _ = ValidateArgAddressOrExit(address)

	rootDir := "."
	if c.NArg() > 1 {
		rootDir = c.Args().Get(1)
	}
	rootDir = ValidateArgRootDirOrExit(rootDir)

	fmt.Println("Loading arguments:", address, rootDir, ops)
	//fmt.Println("listening:", address, mUsingVirtualHost, mNoTrailingSlash)
	return libs.RealServer(&Configure{rootDir, address, ops, nil, nil, nil})
}
