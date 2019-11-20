package core

import (
	"net/http"

	"github.com/zhanbei/serve-static"
)

var mStaticServer *servestatic.FileServer

// Another pattern is to create server for all existed sites, with standalone configurations.
func VirtualHostStaticHandler(ss *servestatic.FileServer) http.Handler {
	mStaticServer = ss
	return ss
}
