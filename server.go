package main

import (
	"net/http"
	"net/url"
	"os"

	"kube-oidc-proxy/env"
	"kube-oidc-proxy/routes/callback"
	"kube-oidc-proxy/routes/login"
	"kube-oidc-proxy/routes/logout"
)

func StartServer() error {
	u, err := url.Parse(env.Env.KubernetesDashboardURL)
	if err != nil {
		return err
	}
	p := &proxy{targetURL: u, hostHeader: "localhost"}
	http.HandleFunc("/", p.reverse)
	http.HandleFunc("/oauth2/login", login.LoginHandler)
	http.HandleFunc("/oauth2/callback", callback.CallbackHandler)
	http.HandleFunc("/oauth2/logout", logout.LogoutHandler)
	return p.run()
}

func fileIsExists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}
