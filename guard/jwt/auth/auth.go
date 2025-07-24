package main

import (
	"context"
	"github.com/aionbit/ixo-plugins/guard"
	"github.com/aionbit/ixo-plugins/plugin"
	"github.com/golang-jwt/jwt/v5"
)

var PluginInstance auth

type auth struct {
}

var exception = plugin.NewError("guard/jwt/auth")

type request struct {
	Token  string `mapstructure:"token"`
	Secret string `mapstructure:"secret"`
}

func (a *auth) error(v ...any) error {
	return exception.Throw(v...)
}

func (a *auth) getRequest(input any) (*request, error) {
	var req request
	if err := plugin.DecodeInput(input, &req); err != nil {
		return nil, a.error("decode input error", err)
	}
	if req.Token == "" {
		return nil, a.error("token is required")
	}
	return &req, nil
}

// Run Note
func (a *auth) Run(ctx context.Context, input any) (any, error) {
	req, err := a.getRequest(input)
	if err != nil {
		return nil, err
	}
	token, err := jwt.ParseWithClaims(req.Token, &guard.AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(req.Secret), nil
	})
	if err != nil {
		return nil, a.error("parse token error", err)
	}
	if claims, ok := token.Claims.(*guard.AuthClaims); ok && token.Valid {
		data := map[string]any{
			"user_id": claims.UserID,
			"issuer":  claims.Issuer,
		}
		if claims.ExpiresAt != nil {
			data["expires_at"] = claims.ExpiresAt.String()
		}
		if claims.IssuedAt != nil {
			data["issued_at"] = claims.IssuedAt.String()
		}
		return data, nil
	} else {
		return nil, a.error(jwt.ErrTokenInvalidClaims)
	}
}
