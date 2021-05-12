package main

import (
	"os"

	"github.com/n-creativesystem/oidc-proxy/adapter"
	"github.com/n-creativesystem/oidc-proxy/logger"
)

func main() {
	app := adapter.New(Version, Revision)
	if err := app.Run(os.Args); err != nil {
		logger.Log.Error(err)
	}
}
