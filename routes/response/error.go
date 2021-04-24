package response

import (
	"net/http"

	"kube-oidc-proxy/logger"
)

func Error(w http.ResponseWriter, err string, code int) {
	logger.Log.Critical(err)
	http.Error(w, err, code)
}
