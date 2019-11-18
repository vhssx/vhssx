# Static Server

<!-- > 2018-07-05T00:23:46+0800 -->

<!-- Titles: *Static Server*, *Http Server*. -->

A static server in Go, supporting hosting static files in the no-trailing-slash version.

## Installation

Get and install the repository using `go get`.

```bash
go get -u -v github.com/zhanbei/static-server
```

## Applications

There are multiple entrances for several standalone applications,
lightweight and special-purposed,
or heavy and powerful.

1. A/The Simple Static Site Server `static-server` `assss` `zss`
1. Lightweight VHSS `vhssl`
1. VHSS with MongoDB Driver `vhssdb`

## Configures and Options

Support multiple targets(file/stdout/files/mongodb) in multiple(possible) formats(text/json/custom).

| Formats\\Targets | Stdout | File(s) | MongoDB |
| :---: | :---: | :---: | :---: |
| **Text** | `+++` | `++` | `---` |
| **JSON** | `+` | `+++` | `---` |
| **BSON** | `---` | `---` | `+++` |

- Server
	- Virtual Hosting `false`
		- Whether to enable virtual hosting; @see https://en.wikipedia.org/wiki/Virtual_hosting
	- No Trailing Slash `false`
		- Hosting static files in the no-trailing-slash mode.
	- Directory Listing `false`
		- Listing files of a directory if the `index.html` is not found when in the normal mode.
	- TrustProxyIp `false`
		- This server will be running behind a reverse proxy,
		 and prefer to fetch the remote ip from header [ `X-Remote-Addr` > `X-Forwarded-For` > `IP` ] over `ctx.Request.Ip`.
- Loggers
	- Enabled `false`
	- Format `"text" | "json"`
	- Per Host `false`
		- A log file per host.
	- Target `"stdout" | ${file}` `${dir}`
		- The value should be `"stdout"|${file}` if the `perHost` option is false.
		- The value should be `${dir}` if the `perHost` option is true.
- Mongo
	- Enabled `false`
	- Connection `"mongodb://${hostname}:${port=27017}"`
	- Database `string`
	- Collection Prefix `"logging.vhss"`

Libraries to receive arguments and options.

1. Command-line Arguments
	- [urfave/cli: A simple, fast, and fun package for building command line apps in Go](https://github.com/urfave/cli)
	- [go-micro/config/source/cli at master · micro/go-micro](https://github.com/micro/go-micro/tree/master/config/source/cli)
1. Configuration File
	- [go-micro/config/source/file at master · micro/go-micro](https://github.com/micro/go-micro/tree/master/config/source/file)


## Examples of Usage

```bash
# Serving the current folder on a specific port.
static-server 1234

# Serving a specific folder on a specific port.
static-server 1234 ./some-folder/site-root

# Serving a folder in the no-trailing-slash mode.
static-server --no-trailing-slash 1234 ./some-folder/site-root

# Serving a folder in the no-trailing-slash mode with virtual hosting.
static-server --no-trailing-slash --enable-virtual-hosting 1234 ./some-folder/site-root
```

Run `static-server --help` to get the help text, like shown below:

```text
NAME:
   static-server - A static server in Go, supporting hosting static files in the no-trailing-slash version.

USAGE:
   static-server [global options] [<http-address>:]<http-port> <www-root-directory>

VERSION:
   0.9.0

DESCRIPTION:
   A static server in Go, supporting hosting static files in the no-trailing-slash version.

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --enable-virtual-hosting  Whether to enable virtual hosting; @see https://en.wikipedia.org/wiki/Virtual_hosting
   --no-trailing-slash       Hosting static files in the no-trailing-slash mode.
   --help, -h                show help
   --version, -v             print the version
```
