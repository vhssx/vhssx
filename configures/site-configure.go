package configures

type SiteConfigure struct {
	DirectoryListing *bool `json:"directoryListing"`

	NoTrailingSlash *bool `json:"noTrailingSlash"`
	// Dynamic routing, from a group of paths to a fixed target.
	RouteMappings *RouteMappings `json:"routeMappings"`

	PrivatePages []string `json:"privatePages"`
	// Allow all robots or disallow.
	//RobotRules []string `json:"robotRules"`
}

type RouteMappings = []RouteMapping
type RouteMapping struct {
	Regexp string `json:"regexp"`

	Target string `json:"target"`
}
