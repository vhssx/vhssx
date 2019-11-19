package main

import (
	"os"

	"github.com/urfave/cli"
	. "github.com/zhanbei/static-server/configs"
	"github.com/zhanbei/static-server/helpers/terminator"
	"github.com/zhanbei/static-server/libs"
)

var ops = new(ServerOptions)

func init() {
	ops.TrustProxyIp = true
}

// The primary program entrance.
// Support more custom built, like for lite/medium/heavy programs, for cli/gui(with different themes) modes, and for linux/windows/mac platforms.
// @see [Support multiple entrances and keep the current one as the primary. · Issue #6 · zhanbei/static-server](https://github.com/zhanbei/static-server/issues/6)
func main() {
	app := cli.NewApp()
	app.Name = "static-server"
	app.Usage = "A static server in Go, supporting hosting static files in the no-trailing-slash version."
	app.Version = "0.9.0"
	app.Description = "A static server in Go, supporting hosting static files in the no-trailing-slash version."
	app.UsageText = "static-server [global options] [<http-address>:]<http-port> <www-root-directory>"
	app.Flags = []cli.Flag{
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

	//fmt.Println("listening:", address, mUsingVirtualHost, mNoTrailingSlash)
	return libs.RealServer(&Configure{address, rootDir, ops, nil, nil, nil})
}
