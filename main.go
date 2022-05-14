package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func rewriteHost(root *url.URL, basePath string, p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// From https://gist.github.com/JalfResi/6287706
		r.Host = root.Host
		relPath, err := filepath.Rel(basePath, r.URL.Path)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		path := filepath.Join("/", relPath)
		r.URL.Path = path
		p.ServeHTTP(w, r)
	}
}

func printUsageAndExit() {
	fmt.Printf(`Usage: %s <remote url to proxy> <local URL or path>

Example: %s "http://example.com" "/"

The remote URL must be a full URL.

The local URL can be:

- A URL path starting with a slash, e.g. "/" or "/test/" (defaults to localhost)
- A URL with implied HTTPS, e.g. "domain.test/path/" (defaults to HTTPS)
- A full URL using HTTP or HTTPS, e.g. "http://domain.test:8000/path/"
`, os.Args[0], os.Args[0])
	os.Exit(1)
}

func main() {
	flag.Parse()
	if flag.NArg() < 2 {
		printUsageAndExit()
	}

	remoteRoot := flag.Arg(0)
	localRoot := flag.Arg(1)

	remoteRootURL, err := url.Parse(remoteRoot)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remoteRootURL)

	if strings.HasPrefix(localRoot, "/") {
		localRoot = fmt.Sprintf("localhost%s", localRoot)
	}
	if !strings.Contains(localRoot, "://") {
		localRoot = fmt.Sprintf("https://%s", localRoot)
	}
	localRootURL, err := url.Parse(localRoot)
	if err != nil {
		panic(err)
	}
	localRootDomain := localRootURL.Hostname()
	localRootPath := localRootURL.Path
	localRootPortString := func(defaultPortString string) string {
		if localRootURL.Port() != "" {
			return fmt.Sprintf(":%s", localRootURL.Port())
		}
		return defaultPortString
	}

	if len(localRootPath) == 0 {
		localRootURL.Path = "/"
	}
	http.HandleFunc(localRootURL.Path, rewriteHost(remoteRootURL, localRootURL.Path, proxy))

	fmt.Printf(`Serving from remote URL:

    %s

To local URL:

    %s

`, remoteRootURL, localRootURL)
	if localRootDomain != "localhost" {

		fmt.Printf(`Make sure the following domains is in /etc/hosts before connecting:

    %s

`, localRootDomain)
	}

	if localRootURL.Scheme == "http" {
		err = http.ListenAndServe(localRootPortString(":80"), nil)
		if err != nil {
			panic(err)
		}
	} else if localRootURL.Scheme == "https" {
		tryServe := func() error {
			return http.ListenAndServeTLS(localRootPortString(":443"),
				fmt.Sprintf("/Users/lgarron/.data/servecert/certs/%s/%s.pem", localRootDomain, localRootDomain),
				fmt.Sprintf("/Users/lgarron/.data/servecert/certs/%s/%s-key.pem", localRootDomain, localRootDomain),
				nil)
		}
		err = tryServe()
		if err != nil {
			if !strings.Contains(err.Error(), "no such file") {
				panic(err)
			}
			fmt.Printf(`--------
Could not find a certificate. Running:

    ./script/cert %s

`, localRootDomain)
			cmd := exec.Command("./script/cert", localRootDomain)
			err = cmd.Run()
			if err != nil {
				panic(err)
			}
			err := tryServe()
			if err != nil {
				panic(err)
			}
		}
	} else {
		panic(fmt.Sprintf("Unexpected scheme: %s", localRootURL.Scheme))
	}
}