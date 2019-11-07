package libs

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zhanbei/serve-static"
)

// Disable directory listing with http.FileServer
func GetNoDirListingHandler(rootDir string, mDirectoryListing bool) http.Handler {
	// @see https://stackoverflow.com/questions/26559557/how-do-you-serve-a-static-html-file-using-a-go-web-server
	// @see https://groups.google.com/forum/#!msg/golang-nuts/bStLPdIVM6w/hidTJgDZpHcJ
	// @see https://www.alexedwards.net/blog/disable-http-fileserver-directory-listings
	handler := http.FileServer(http.Dir(rootDir))
	if mDirectoryListing {
		fmt.Println("Enabled directory listing.")
		return handler
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqPath := r.URL.Path
		if strings.HasSuffix(reqPath, "/") && reqPath != "/" {
			exists, _ := servestatic.IsFileRegular(rootDir, reqPath, servestatic.IndexDotHtml)
			if !exists {
				http.NotFound(w, r)
				return
			}
		}
		handler.ServeHTTP(w, r)
	})
}
