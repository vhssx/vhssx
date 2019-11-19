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
