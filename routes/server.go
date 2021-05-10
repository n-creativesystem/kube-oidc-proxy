package routes

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-plugin"
	"github.com/n-creativesystem/oidc-proxy/config"
	"github.com/n-creativesystem/oidc-proxy/logger"
	"github.com/n-creativesystem/oidc-proxy/routes/callback"
	errRouter "github.com/n-creativesystem/oidc-proxy/routes/errors"
	"github.com/n-creativesystem/oidc-proxy/routes/login"
	"github.com/n-creativesystem/oidc-proxy/routes/logout"
	"github.com/n-creativesystem/oidc-proxy/routes/proxy"
)

type Handler interface {
	http.Handler
	Close()
}

type MultiHost map[string]Handler

func (m MultiHost) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler := m[r.Host]; handler != nil {
		handler.ServeHTTP(w, r)
	} else {
		errRouter.ErrorHandler(logger.Log)(w, r, errors.New(http.StatusText(http.StatusNotFound)))
	}
}

func (m MultiHost) Close() {
	for _, value := range m {
		if value != nil {
			value.Close()
		}
	}
}

type muxWrap struct {
	mux              *http.ServeMux
	GetConfiguration config.GetConfiguration
	cachePlugin      *plugin.Client
}
type confHandler func(configuration config.GetConfiguration) http.HandlerFunc

func (m *muxWrap) Close() {
	if m.cachePlugin != nil {
		m.cachePlugin.Kill()
	}
}

func (m *muxWrap) HandleFunc(pattern string, handler confHandler) {
	m.mux.HandleFunc(pattern, handler(m.GetConfiguration))
}

func (m *muxWrap) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}

func New(configuration config.GetConfiguration) (Handler, error) {
	conf := configuration()
	m := http.NewServeMux()
	mux := &muxWrap{
		mux:              m,
		GetConfiguration: configuration,
		cachePlugin:      conf.Store.GetPluginClient().(*plugin.Client),
	}
	host := conf.GetHostname()
	for _, location := range conf.Locations {
		u, err := url.Parse(location.ProxyPass)
		if err != nil {
			return nil, err
		}
		p := proxy.New(u, host)
		for _, path := range location.Urls {
			p.AddRoute(path.Path, func() string {
				return path.Token
			})
		}
		mux.HandleFunc("/", p.Reverse)
	}
	mux.HandleFunc(conf.Login, login.LoginHandler)
	mux.HandleFunc(conf.Callback, callback.CallbackHandler)
	mux.HandleFunc(conf.Logout, logout.LogoutHandler)
	return mux, nil
}
