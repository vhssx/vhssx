# Designs

<!-- > 2019-11-18T15:08:24+0800 -->

## Work-flows and Error Handling

1. Initialize Application
1. Loading Configures and Options
1. Validate the Options Syntax
1. Panic on Pre-launch Errors
1. Take Care of Runtime Errors

## Topics

- Warnings and Logs
	- Panics for invalid configures;
	- Warning for abnormal configures;
	- Logging out the configurations resolved to confirm;
- Runtime Reload
	- No Panics
	- Hot Loading
- Failed Records
	- Writing to file or Stdout?
- Log Format
	- How to note down the request performance.

## Packages and Modules

- Applications/Entrances
	- vhssl
	- vhssdb
- Core Libraries
	- MongoDB
	- Real Server
		- Built-in Recorder
- Configuration Options
	- Loggers
		- Basic Logger
		- Gorilla Logger
	- MongoDB Options
- Other Libs
	- Helpers
	- utils

## Default Logger

- Console Logger in the Combined Log Format
	- No Gorilla Logger, Custom Loggers, and MongoDB Logger.

## Schedules

- [x] Upgrading MongoDB
- [x] Combine Gorilla Writers
- [x] Extract IRecorder
- [x] Upgrade recorder times.
- [x] Integrating MongoDB
- [x] Upgrade the logger formats(combined/extended).
- [x] Support _.other.sites
- [ ] 404 Status Code
- [ ] Upgrading the Cli Package
- [ ] Add a cmd/test(-t) to validate the specific configures and report.
- [ ] Development in the virtual hosts mode.
	- Modify the request host following preferences.
- [ ] Site Configures per site in the virtual hosts modes
	- [x] Scan and Discover Sites for Configures and Cache
		- [ ] site-config.toml
		- [ ] no trailing slash
		- [ ] enable folder listing
	- [ ] Path Mapping
		- `^/[^.]+$` -> `index.html`
	- [ ] Parking and Counting `++`
		- [ ] Dynamic rendering!
	- [ ] Filters: hide private files, like "README.md", "site-config.toml". `+`
	- [ ] Robots: support allowed rules and disallowed rules. `---`

## Fallthrough and Custom 404 Page

> Implement with powerful features and optimize later?

Powerful vs Performance?

- [x] Fallthrough for resources, like:
	- index.html
		- Site for all sub sites, or use nginx to direct to the right.
		- Only for very few special usages.
	- favicon.ico
		- This is common because all sub sites may share the same favicon.
		- Do not copy it again and again.
	- robots.txt
	- Other Resources(Scripts and Styles)
		- May not be copied nor fallthrough, but better be referred cross site.
		- May be used to as the custom *404* pages, corresponding with the route mapping.
- [x] Fallthrough by domains `+++`
	- (*.)domain.com --> _.domain.com --> _.default.com
	- Scan the folders for site configures.
	- Site/404.html --> Scope/404 --> Global/404
- [ ] Custom Pages(404/50x) for different hosts. `+++`
	- [ ] Dynamic rendering?
	- [ ] Enable by default and for all sites?
	- [ ] pages/404.html

## Use Cases

Use `stdout` target for development or by default, while use the `file` target for production.

1. Minimal `+`
	- `Gorilla/Combined -> Stdout/File`
1. Standard `+`
	- `Gorilla/Combined -> Stdout/File`
	- `JSON -> File`
1. Extended `+++`
	- `Builtin/Extended -> Stdout/File`
	- `JSON -> File`
1. Mongo `+++++`
	- `Builtin/Extended -> Stdout/File`
	- `BSON -> MongoDB`
1. Development `+++++`
	- `Gorilla/Combined -> File`
	- `Builtin/Extended -> File`
	- `JSON -> File`
	- `BSON -> MongoDB`
