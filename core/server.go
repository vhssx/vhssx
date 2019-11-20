package core

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/zhanbei/serve-static"
	"github.com/zhanbei/static-server/conf"
	"github.com/zhanbei/static-server/helpers/terminator"
	"github.com/zhanbei/static-server/helpers/writersHelper"
	"github.com/zhanbei/static-server/recorder"
)

func RealServer(cfg *conf.Configure, loggers recorder.IRecorders) error {
	ops := cfg.Server
	var handler http.Handler
	if !ops.NoTrailingSlash {
		// Hosting in the normal mode.
		handler = GetNoDirListingHandler(cfg.RootDir, ops.DirectoryListing)
	} else {
		fmt.Println("Hosting static files in the " + conf.OptionNameNoTrailingSlash + " mode.")
		if ops.UsingVirtualHost {
			fmt.Println("Enabled virtual hosting based on request.Host; @see https://en.wikipedia.org/wiki/Virtual_hosting.")
		}
		mStaticServer, err := servestatic.NewFileServer(cfg.RootDir, ops.UsingVirtualHost)
		if err != nil {
			terminator.ExitWithPreLaunchServerError(err, "ERROR: The specified www-root-directory does not exist: "+cfg.RootDir)
		}
		handler = mStaticServer
	}
	fmt.Println("Looking after directory:", cfg.RootDir)
	gor := cfg.GorillaOptions
	if gor != nil && gor.Enabled {
		target := writersHelper.StdoutVsFileWriter(gor.Stdout, gor.LogWriter)
		if target != nil {
			handler = handlers.CombinedLoggingHandler(target, handler)
		} else {
			fmt.Println("Warning: both the stdout and the target are nil!")
		}
	}

	server := &http.Server{
		Addr: cfg.Address,

		Handler: StructuredLoggingHandler(handler, cfg, loggers),
	}

	fmt.Println("Server is running at:", cfg.Address)
	return server.ListenAndServe()
}
