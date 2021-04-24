package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"net/http/httputil"
	"net/url"

	"kube-oidc-proxy/app"
	"kube-oidc-proxy/env"
	"kube-oidc-proxy/logger"
	"kube-oidc-proxy/routes/response"
)

type proxy struct {
	targetURL  *url.URL
	hostHeader string
}

func (p *proxy) run() error {
	port := env.Env.ProxyPort
	if port != "" && port[:1] != ":" {
		port = ":" + port
	}
	logger.Log.Info(fmt.Sprintf("Server listening on http://localhost%s/", port))
	certExists := fileIsExists(env.Env.CertFile)
	keyExists := fileIsExists(env.Env.KeyFile)
	if certExists && keyExists {
		return http.ListenAndServeTLS(port, env.Env.CertFile, env.Env.KeyFile, nil)
	} else {
		return http.ListenAndServe(port, nil)
	}
}

func (p *proxy) reverse(w http.ResponseWriter, r *http.Request) {
	session, err := app.Store.Get(r, "auth-session")
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rawIdToken := ""
	idToken := session.Values["id_token"]
	if token, ok := idToken.(string); ok {
		rawIdToken = token
	} else {
		http.Redirect(w, r, "/oauth2/login", http.StatusTemporaryRedirect)
		return
	}
	if rawIdToken != "" {
		r.Header.Add("Authorization", "Bearer "+rawIdToken)
	} else {
		http.Redirect(w, r, "/oauth2/login", http.StatusTemporaryRedirect)
		return
	}

	director := func(req *http.Request) {
		req.URL.Scheme = p.targetURL.Scheme
		req.URL.Host = p.targetURL.Host
		req.Host = p.hostHeader
	}

	t := &transport{}
	trace := &httptrace.ClientTrace{
		GotConn: t.GotConn,
	}

	r = r.WithContext(httptrace.WithClientTrace(r.Context(), trace))

	reverse := &httputil.ReverseProxy{Director: director, Transport: t}
	reverse.FlushInterval = -1
	reverse.ServeHTTP(w, r)
}

type transport struct {
	current *http.Request
}

func (t *transport) GotConn(info httptrace.GotConnInfo) {
	logger.Log.Info(fmt.Sprintf("Connected to %v", t.current.URL))
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.current = req
	b, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}
	logger.Log.Debug(string(b))
	tr := http.DefaultTransport.(*http.Transport)
	tr.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	return tr.RoundTrip(req)
}
