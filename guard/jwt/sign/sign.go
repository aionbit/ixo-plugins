package main

import (
	"context"
	"github.com/aionbit/ixo-plugins/guard"
	"github.com/aionbit/ixo-plugins/plugin"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var PluginInstance sign

type sign struct {
}

var exception = plugin.NewError("guard/jwt/sign")

type request struct {
	SigningMethod string `mapstructure:"signing_method"`
	Secret        string `mapstructure:"secret"`
	UserID        string `mapstructure:"user_id"`
	TTL           string `mapstructure:"ttl"`
	Issuer        string `mapstructure:"issuer"`

	signingMethod jwt.SigningMethod
	ttl           time.Duration
}

func (s *sign) error(v ...any) error {
	return exception.Throw(v...)
}

func (s *sign) getRequest(input any) (*request, error) {
	var req request
	if err := plugin.DecodeInput(input, &req); err != nil {
		return nil, s.error("decode input error", err)
	}
	if req.SigningMethod == "" {
		req.SigningMethod = "HS256"
	}
	if m := jwt.GetSigningMethod(req.SigningMethod); m == nil {
		return nil, s.error("invalid signing_method " + req.SigningMethod)
	} else {
		req.signingMethod = m
	}
	if req.Secret == "" {
		return nil, s.error("secret is required")
	}
	if req.UserID == "" {
		return nil, s.error("user_id is required")
	}
	if req.TTL == "" {
		return nil, s.error("ttl is required")
	} else {
		d, err := time.ParseDuration(req.TTL)
		if err != nil {
			return nil, s.error("parse ttl error", d)
		}
		if d <= 0 {
			return nil, s.error("invalid ttl")
		}
		req.ttl = d
	}
	if req.Issuer == "" {
		req.Issuer = "IxO-Gateway"
	}
	return &req, nil
}

// Run Note
func (s *sign) Run(ctx context.Context, input any) (any, error) {
	req, err := s.getRequest(input)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	claims := guard.AuthClaims{
		UserID: req.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(req.ttl)), // 1小时过期
			Issuer:    req.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(req.signingMethod, claims)
	t, err := token.SignedString([]byte(req.Secret))
	if err != nil {
		return nil, s.error("token signed string error", err)
	}
	return map[string]any{
		"token": t,
	}, nil
}
