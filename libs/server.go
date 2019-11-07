package libs

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/zhanbei/serve-static"
)

func RealServer(ops *ServerOptions, address, rootDir string) error {
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
			fmt.Println("ERROR: The specified www-root-directory is invalid:" + rootDir)
			log.Fatal(err)
		}
		handler = mStaticServer
	}
	fmt.Println("Looking after directory:", rootDir)
	handler = handlers.CombinedLoggingHandler(os.Stdout, handler)

	server := &http.Server{
		Addr: address,

		Handler: StructuredLoggingHandler(handler, ops),
	}

	fmt.Println("Server is running at:", address)
	return server.ListenAndServe()
}
