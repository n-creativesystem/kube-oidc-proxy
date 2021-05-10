package response

import (
	"net/http"

	"github.com/n-creativesystem/oidc-proxy/logger"
)

func Error(w http.ResponseWriter, err string, code int) {
	logger.Log.Critical(err)
	http.Error(w, err, code)
}
