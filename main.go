package main

import (
	"kube-oidc-proxy/app"
)

func main() {
	app.Init()
	StartServer()
}
