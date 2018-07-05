# Static Server

<!-- > 2018-07-05T00:23:46+0800 -->

<!-- Titles: *Static Server*, *Http Server*. -->

A static server in Go, supporting hosting static files in the no-trailing-slash version.

## Installation

Get and install the repository using `go get`.

```bash
go get -u -v github.com/zhanbei/static-server
```

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
