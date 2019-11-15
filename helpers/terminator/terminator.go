package terminator

import (
	"log"
)

const (
	ExitCodeSuccess = 0
	// Unknown error.
	ExitCodeUnknownError = 1
	// The configuration options are invalid.
	ExitCodeConfigError = 2
	// Ready to start the server, but failed because of insufficient resources required or conflicts.
	ExitCodePreLaunchFatalError = 3
)

// Exit with exit error.
func exitWithError(code int, err error, message string) {
	if err != nil {
		log.Fatal(message, err)
	} else {
		log.Fatal(message)
	}
	//_, _ = fmt.Fprintln(os.Stderr, message)
	// The exit code is not yet used!
	//os.Exit(code)
}

func ExitWithConfigError(err error, message string) {
	exitWithError(ExitCodeConfigError, err, "[Configures] "+message)
}

func ExitWithPreLaunchServerError(err error, message string) {
	exitWithError(ExitCodePreLaunchFatalError, err, "[Server] "+message)
}
