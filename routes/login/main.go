package login

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/n-creativesystem/oidc-proxy/auth"
	"github.com/n-creativesystem/oidc-proxy/config"
	"github.com/n-creativesystem/oidc-proxy/routes/response"
)

func LoginHandler(configuration config.GetConfiguration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conf := configuration()
		if r.Method != http.MethodGet {
			return
		}
		// Generate random state
		b := make([]byte, 32)
		_, err := rand.Read(b)
		if err != nil {
			response.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		state := base64.StdEncoding.EncodeToString(b)

		session, err := conf.Store.Get(r, conf.SessionName)
		if err != nil {
			response.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.Values["state"] = state
		err = session.Save(r, w)
		if err != nil {
			response.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		authenticator, err := auth.NewAuthenticator(conf)
		if err != nil {
			response.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, authenticator.Config.AuthCodeURL(state), http.StatusTemporaryRedirect)
	}
}
