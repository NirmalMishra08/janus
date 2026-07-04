package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewHandler(target string) (http.Handler, error) {

	targetUrl, err := url.Parse(target)
	if err != nil {
		return nil , err
	}

	proxy := httputil.NewSingleHostReverseProxy(targetUrl)

	return proxy, nil 
}
