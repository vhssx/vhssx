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
- [ ] Upgrading the Cli Package
- [ ] Add a cmd/test(-t) to validate the specific configures and report.
- [ ] Development in the virtual hosts mode.
	- Modify the request host following preferences.
- [ ] Site Configures per site in the virtual hosts modes
	- [ ] Fallthrough for all sites, like index.html, and favicon.ico, robots.txt
		- * -> _.default.com
	- [ ] Scan and Discover Sites for Configures and Cache
		- [ ] site-config.toml
		- [ ] no trailing slash
		- [ ] enable folder listing
	- [ ] Fallthrough for subdomains, like default favicon.ico `+++`
		- (*.)domain.com -> _.domain.com
		- Scan the folders for site configures.
		- Target/404 > Site/404.html > Scope > Global
	- [ ] Path Mapping
		- `^/[^.]+$` -> `index.html`
	- [ ] Custom Pages(404/50x) for different hosts. `+++`
		- [ ] Dynamic rendering?
		- [ ] Enable by default and for all sites?
		- [ ] pages/404.html
	- [ ] Parking and Counting `++`
		- [ ] Dynamic rendering!
	- [ ] Filters: hide private files, like "README.md", "site-config.toml". `+`
	- [ ] Robots: support allowed rules and disallowed rules. `---`

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
