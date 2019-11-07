package libs

const (
	OptionNameEnableVirtualHosting = "enable-virtual-hosting"

	OptionNameNoTrailingSlash = "no-trailing-slash"

	OptionNameDirectoryListing = "enable-directory-listing"
)

// The options reading from configuration files, or environment.
type ServerOptions struct {
	DirectoryListing bool

	NoTrailingSlash bool

	UsingVirtualHost bool
}
