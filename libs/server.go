package libs

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/zhanbei/serve-static"
	"github.com/zhanbei/static-server/conf"
	"github.com/zhanbei/static-server/helpers/terminator"
	"github.com/zhanbei/static-server/utils"
)

func RealServer(cfg *conf.Configure) error {
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
		if gor.Stdout || utils.NotEmpty(gor.Target) {
			if gor.Stdout {
				handler = handlers.CombinedLoggingHandler(os.Stdout, handler)
			}
			if utils.NotEmpty(gor.Target) {
				// FIX-ME
				writer, err := os.Open(gor.Target)
				if err != nil {
					fmt.Println("Failed to open file to write logs:", err)
				} else {
					handler = handlers.CombinedLoggingHandler(writer, handler)
				}
			}
		} else {
			fmt.Println("Warning: both the stdout and the target are nil!")
		}
	}

	server := &http.Server{
		Addr: cfg.Address,

		Handler: StructuredLoggingHandler(handler, cfg),
	}

	fmt.Println("Server is running at:", cfg.Address)
	return server.ListenAndServe()
}
