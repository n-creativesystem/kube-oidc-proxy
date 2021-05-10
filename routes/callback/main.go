package callback

import (
	"context"
	"fmt"
	"net/http"

	"github.com/n-creativesystem/oidc-proxy/auth"
	"github.com/n-creativesystem/oidc-proxy/config"
	"github.com/n-creativesystem/oidc-proxy/routes/response"

	"github.com/coreos/go-oidc"
)

func CallbackHandler(configuration config.GetConfiguration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conf := configuration()
		if r.Method != http.MethodGet {
			return
		}
		session, err := conf.Store.Get(r, conf.SessionName)
		if err != nil {
			response.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if r.URL.Query().Get("state") != session.Values["state"] {
			http.Redirect(w, r, conf.Login, http.StatusTemporaryRedirect)
			// response.Error(w, "Invalid state parameter", http.StatusBadRequest)
			return
		}

		authenticator, err := auth.NewAuthenticator(conf)
		if err != nil {
			response.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token, err := authenticator.Config.Exchange(context.TODO(), r.URL.Query().Get("code"))
		if err != nil {
			conf.Log.Critical(fmt.Sprintf("no token found: %v", err))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		rawIDToken, ok := token.Extra("id_token").(string)
		if !ok {
			response.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
			return
		}

		oidcConfig := &oidc.Config{
			ClientID: conf.Oidc.ClientId,
		}

		_, err = authenticator.Provider.Verifier(oidcConfig).Verify(context.TODO(), rawIDToken)

		if err != nil {
			response.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// var profile map[string]interface{}
		// if err := idToken.Claims(&profile); err != nil {
		// 	response.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		// session.Values["profile"] = profile
		auth.SetTokenSession(session, token)
		err = session.Save(r, w)
		if err != nil {
			response.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to logged in page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
