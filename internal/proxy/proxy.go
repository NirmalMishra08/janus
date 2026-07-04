package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewHandler(target string) (http.Handler, error) {

	targetUrl, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(targetUrl)
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = targetUrl.Scheme
		req.URL.Host = targetUrl.Host
		// We keep the original path (/users) as it is
		// No stripping for now
	}

	return proxy, nil
}
