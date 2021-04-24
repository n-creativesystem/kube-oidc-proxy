package app

import (
	"encoding/gob"

	"github.com/gorilla/sessions"
)

var (
	Store *sessions.CookieStore
)

func init() {
	gob.Register(&map[string]interface{}{})
}

func Init() error {
	Store = sessions.NewCookieStore([]byte("something-very-secret"))
	Store.Options.Secure = true
	return nil
}
