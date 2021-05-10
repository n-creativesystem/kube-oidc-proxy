package logout

import (
	"net/http"
	"net/url"

	"github.com/n-creativesystem/oidc-proxy/config"
	"github.com/n-creativesystem/oidc-proxy/routes/response"

	"github.com/gorilla/sessions"
)

func LogoutHandler(configuration config.GetConfiguration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conf := configuration()
		logoutUrl, err := url.Parse(conf.Oidc.Logout)

		if err != nil {
			response.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session, err := conf.Store.Get(r, conf.SessionName)
		if err != nil {
			response.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.Options = &sessions.Options{MaxAge: -1}
		session.Save(r, w)
		conf.Store.Delete(session)
		http.Redirect(w, r, logoutUrl.String(), http.StatusTemporaryRedirect)
	}
}
