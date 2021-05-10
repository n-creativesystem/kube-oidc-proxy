package auth

import (
	"context"
	"fmt"

	"github.com/n-creativesystem/oidc-proxy/config"

	oidc "github.com/coreos/go-oidc"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

type Authenticator struct {
	Provider *oidc.Provider
	Config   oauth2.Config
	Ctx      context.Context
}

func NewAuthenticator(conf config.Servers) (*Authenticator, error) {
	oidcConf := conf.Oidc
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, oidcConf.Provider)
	if err != nil {
		conf.Log.Critical(fmt.Sprintf("failed to get provider: %v", err))
		return nil, err
	}
	o2conf := oauth2.Config{
		ClientID:     oidcConf.ClientId,
		ClientSecret: oidcConf.ClientSecret,
		RedirectURL:  oidcConf.RedirectUrl,
		Endpoint:     provider.Endpoint(),
		Scopes:       oidcConf.Scopes,
	}

	return &Authenticator{
		Provider: provider,
		Config:   o2conf,
		Ctx:      ctx,
	}, nil
}

func SetTokenSession(session *sessions.Session, token *oauth2.Token) {
	rawIdToken, _ := token.Extra("id_token").(string)
	session.Values["id_token"] = rawIdToken
	session.Values["access_token"] = token.AccessToken
	session.Values["refresh_token"] = token.RefreshToken
}
