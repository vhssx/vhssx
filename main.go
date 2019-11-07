package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/urfave/cli"
	. "github.com/zhanbei/static-server/libs"
)

var mUsingVirtualHost = false
var mNoTrailingSlash = true
var mDirectoryListing = false

func Action(c *cli.Context) error {
	if !mNoTrailingSlash && mUsingVirtualHost {
		fmt.Println("ERROR: Sorry, currently virtual hosting is supported only in the " + OptionNameNoTrailingSlash + " mode.")
		log.Fatal("You may add the --" + OptionNameNoTrailingSlash + " option to use --" + OptionNameEnableVirtualHosting + " option.")
	}

	if c.NArg() <= 0 {
		log.Fatal("Please specify a port, like `static-server 8080`.")
	}
	address := c.Args().Get(0)
	port, err := strconv.Atoi(address)
	if err != nil {
		// Check the address.
	} else {
		// The address is only a port.
		if port < 1 || 65535 < port {
			log.Fatal("ERROR: unavailable port[" + strconv.Itoa(port) + "]; make sure http port is number and is limited to <0-65535>.")
		}
		if port <= 1024 {
			fmt.Println("WARNING: the port[" + strconv.Itoa(port) + "] specified is not bigger than 1024; root privileges may be needed!")
		}
		address = ":" + strconv.Itoa(port)
	}

	rootDir := "."
	if c.NArg() > 1 {
		rootDir = c.Args().Get(1)
	}
	rootDir, err = filepath.Abs(rootDir)
	if err != nil {
		fmt.Println("ERROR: The specified www-root-directory is invalid:" + rootDir)
		log.Fatal(err)
	}

	//fmt.Println("listening:", address, mUsingVirtualHost, mNoTrailingSlash)
	return RealServer(&ServerOptions{mDirectoryListing, mNoTrailingSlash, mUsingVirtualHost}, address, rootDir)
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

			Destination: &mUsingVirtualHost,
		},
		cli.BoolFlag{
			Name: OptionNameNoTrailingSlash,

			Usage: "Hosting static files in the " + OptionNameNoTrailingSlash + " mode.",

			Destination: &mNoTrailingSlash,
		},
		cli.BoolFlag{
			Name: OptionNameDirectoryListing,

			Usage: "Listing files of a directory if the index.html is not found when in the normal mode.",

			Destination: &mDirectoryListing,
		},
	}
	app.Action = Action

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
