package conf

import (
	"fmt"

	"github.com/zhanbei/static-server/helpers/terminator"
)

const (
	OptionNameEnableVirtualHosting = "enable-virtual-hosting"

	OptionNameNoTrailingSlash = "no-trailing-slash"

	OptionNameDirectoryListing = "enable-directory-listing"
)

// The options reading from configuration files, or environment.
type ServerOptions struct {
	DirectoryListing bool `json:"directoryListing"`

	NoTrailingSlash bool `json:"noTrailingSlash"`

	UsingVirtualHost bool `json:"virtualHosting"`
	// Header [ X-Remote-Addr > X-Forwarded-For > IP ]
	TrustProxyIp bool `json:"trustProxyIp"`
}

// Get a default options for server.
func NewDefaultServerOptions() *ServerOptions {
	return &ServerOptions{false, false, false, true}
}

func (m *ServerOptions) IsValid() bool {
	if !m.NoTrailingSlash && m.UsingVirtualHost {
		return false
	}
	return true
}

func (m *ServerOptions) ValidateOrExit() {
	if !m.IsValid() {
		fmt.Println("ERROR: Sorry, currently virtual hosting is supported only in the " + OptionNameNoTrailingSlash + " mode.")
		terminator.ExitWithConfigError(nil, "You may add the --"+OptionNameNoTrailingSlash+" option to use --"+OptionNameEnableVirtualHosting+" option.")
	}
}
