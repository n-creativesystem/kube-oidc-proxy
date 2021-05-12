package routes

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/n-creativesystem/oidc-proxy/config"
	"github.com/n-creativesystem/oidc-proxy/logger"
)

type MultiHost map[string]Handler

func (m MultiHost) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler := m[r.Host]; handler != nil {
		handler.ServeHTTP(w, r)
	} else {
		errorResponse(logger.Log)(w, r, errors.New(http.StatusText(http.StatusNotFound)))
	}
}

func New(configuration config.GetConfiguration) (Handler, error) {
	conf := configuration()
	router := new(conf)
	host := conf.GetHostname()
	for _, location := range conf.Locations {
		u, err := url.Parse(location.ProxyPass)
		if err != nil {
			return nil, err
		}
		for _, path := range location.Urls {
			router.Proxy(path.Path, u, host, path.Type, path.Token, location.IsProxySSLVerify())
		}
	}
	router.Login(conf.Login)
	router.Callback(conf.Callback)
	router.Logout(conf.Logout)
	return router, nil
}
