package auth

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/oauth2"

	"kube-oidc-proxy/env"
	"kube-oidc-proxy/logger"

	oidc "github.com/coreos/go-oidc"
)

type Authenticator struct {
	Provider *oidc.Provider
	Config   oauth2.Config
	Ctx      context.Context
}

func NewAuthenticator() (*Authenticator, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, env.Env.OIDC_PROVIDER)
	if err != nil {
		logger.Log.Critical(fmt.Sprintf("failed to get provider: %v", err))
		return nil, err
	}
	scopes := []string{oidc.ScopeOpenID, oidc.ScopeOfflineAccess, "profile", "email"}
	ss := strings.Split(env.Env.OIDC_SCOPE, ",")
	scopes = append(scopes, ss...)
	conf := oauth2.Config{
		ClientID:     env.Env.OIDC_CLIENT_ID,
		ClientSecret: env.Env.OIDC_CLIENT_SECRET,
		RedirectURL:  env.Env.OIDC_REDIRECT_URL + "/oauth2/callback",
		Endpoint:     provider.Endpoint(),
		Scopes:       scopes,
	}

	return &Authenticator{
		Provider: provider,
		Config:   conf,
		Ctx:      ctx,
	}, nil
}
