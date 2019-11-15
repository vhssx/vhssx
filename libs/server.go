package libs

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/zhanbei/serve-static"
	"github.com/zhanbei/static-server/helpers/terminator"
)

func RealServer(ops *ServerOptions, address, rootDir string, recorder IRecorder) error {
	var handler http.Handler
	if !ops.NoTrailingSlash {
		// Hosting in the normal mode.
		handler = GetNoDirListingHandler(rootDir, ops.DirectoryListing)
	} else {
		fmt.Println("Hosting static files in the " + OptionNameNoTrailingSlash + " mode.")
		if ops.UsingVirtualHost {
			fmt.Println("Enabled virtual hosting based on request.Host; @see https://en.wikipedia.org/wiki/Virtual_hosting.")
		}
		mStaticServer, err := servestatic.NewFileServer(rootDir, ops.UsingVirtualHost)
		if err != nil {
			terminator.ExitWithPreLaunchServerError(err, "ERROR: The specified www-root-directory does not exist: "+rootDir)
		}
		handler = mStaticServer
	}
	fmt.Println("Looking after directory:", rootDir)
	handler = handlers.CombinedLoggingHandler(os.Stdout, handler)

	server := &http.Server{
		Addr: address,

		Handler: StructuredLoggingHandler(handler, ops, recorder),
	}

	fmt.Println("Server is running at:", address)
	return server.ListenAndServe()
}
