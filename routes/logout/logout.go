package logout

import (
	"net/http"
	"net/url"

	"kube-oidc-proxy/app"
	"kube-oidc-proxy/env"
	"kube-oidc-proxy/routes/response"
	"github.com/gorilla/sessions"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	domain := env.Env.OIDC_LOGOUT
	logoutUrl, err := url.Parse(domain)

	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session, err := app.Store.Get(r, "auth-session")
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Options = &sessions.Options{MaxAge: -1, Path: "/"}
	session.Save(r, w)
	http.Redirect(w, r, logoutUrl.String(), http.StatusTemporaryRedirect)
}
