package configures

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/zhanbei/serve-static"
)

const Page404Path = "/404.html"

func Serve404Page(w http.ResponseWriter, ss *servestatic.FileServer, page string) (bool, error) {
	exists, location := ss.GetFilePathFromStatics(page)
	if !exists {
		return false, nil
	}
	// FIX-ME Is it okay to read file and directly serve as 404 pages?
	// FIX-ME What about the response header of #Content-Type?
	bts, err := ioutil.ReadFile(location)
	if err != nil {
		// Found 404 page but failed to serve it, because of unknown error.
		fmt.Println("Found 404 page but failed to serve it, because of unknown error.")
		return false, err
	}
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write(bts)
	return true, nil
}

func ServeDynamicContents(w http.ResponseWriter, r *http.Request, cfg *SiteConfigure, ss *servestatic.FileServer, notFound func()) bool {
	if cfg == nil {
		return false
	}
	if cfg.IsPrivate(r.URL.Path) {
		// 1. Filters for private pages to protect whitelist(hidden resources).
		// Responding the custom 404.
		notFound()
		return true
	}
	target := cfg.GetPotentialMappedTarget(r.URL.Path)
	if target != "" {
		// 2. Custom resources mapping, like  for *Single Page Application* to serve *dynamical* contents.
		exists, location := ss.GetFilePathFromStatics(target)
		if exists {
			http.ServeFile(w, r, location)
		} else {
			notFound()
		}
		return true
	}
	return false
}
