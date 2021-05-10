package proxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/n-creativesystem/oidc-proxy/auth"
	"github.com/n-creativesystem/oidc-proxy/config"
	"github.com/n-creativesystem/oidc-proxy/logger"
	"github.com/n-creativesystem/oidc-proxy/routes/errors"
	errRouter "github.com/n-creativesystem/oidc-proxy/routes/errors"
	"github.com/n-creativesystem/oidc-proxy/routes/response"
	"golang.org/x/oauth2"
)

type Proxy interface {
	Reverse(configuration config.GetConfiguration) http.HandlerFunc
	AddRoute(path string, handle Handle)
}

type proxy struct {
	targetURL  *url.URL
	hostHeader string
	node       *node
}

var _ Proxy = &proxy{}

func New(targetURL *url.URL, hostHeader string) Proxy {
	return &proxy{
		targetURL:  targetURL,
		hostHeader: hostHeader,
		node:       &node{},
	}
}

func (p *proxy) Reverse(configuration config.GetConfiguration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conf := configuration()
		tokenKey := ""
		handler, _, _ := p.node.getValue(r.URL.Path, nil)
		if handler != nil {
			tokenKey = handler()
		}
		log := conf.Log
		session, err := conf.Store.Get(r, conf.SessionName)
		if err != nil {
			response.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rawIdToken, ok := session.Values["id_token"].(string)
		if !ok {
			errors.UnAuthorizedResponse(w, conf.Login)
			// http.Redirect(w, r, conf.Login, http.StatusTemporaryRedirect)
			return
		}
		if rawIdToken != "" {
			authenticator, err := auth.NewAuthenticator(conf)
			oidcConfig := &oidc.Config{
				ClientID: conf.Oidc.ClientId,
			}
			// IDトークンの検証
			_, err = authenticator.Provider.Verifier(oidcConfig).Verify(context.TODO(), rawIdToken)
			if err != nil {
				// トークンの更新
				refreshToken := session.Values["refresh_token"].(string)
				ts := authenticator.Config.TokenSource(r.Context(), &oauth2.Token{
					RefreshToken: refreshToken,
				})
				token, _ := ts.Token()
				auth.SetTokenSession(session, token)
				session.Save(r, w)
			}
			// プロキシ先へ転送するトークンを取得
			rawToken, ok := session.Values[tokenKey].(string)
			if !ok {
				rawToken = rawIdToken
			}
			r.Header.Add("Authorization", "Bearer "+rawToken)
		} else {
			errors.UnAuthorizedResponse(w, conf.Login)
			// http.Redirect(w, r, conf.Login, http.StatusTemporaryRedirect)
			return
		}

		director := func(req *http.Request) {
			req.URL.Scheme = p.targetURL.Scheme
			req.URL.Host = p.targetURL.Host
			req.Host = p.hostHeader
		}

		t := &transport{
			log: log,
		}
		trace := &httptrace.ClientTrace{
			GotConn: t.GotConn,
		}

		*r = *r.WithContext(httptrace.WithClientTrace(r.Context(), trace))

		reverse := &httputil.ReverseProxy{
			Director:     director,
			Transport:    t,
			ErrorHandler: errRouter.ErrorHandler(log),
		}
		reverse.FlushInterval = -1
		reverse.ServeHTTP(w, r)
	}
}

func (p *proxy) AddRoute(path string, handle Handle) {
	p.node.addRoute(path, handle)
}

type transport struct {
	current *http.Request
	log     logger.ILogger
}

func (t *transport) GotConn(info httptrace.GotConnInfo) {
	t.log.Info(fmt.Sprintf("Connected to %v", t.current.URL))
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.current = req
	b, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}
	dumps := strings.Split(string(b), "\n")
	for _, dump := range dumps {
		t.log.Debug(dump)
	}
	tr := http.DefaultTransport.(*http.Transport)
	tr.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	return tr.RoundTrip(req)
}
