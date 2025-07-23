package main

import (
	"context"
	"github.com/aionbit/ixo-plugins/guard"
	"github.com/aionbit/ixo-plugins/plugin"
	"time"
)

var PluginInstance sign

type sign struct {
}

var exception = plugin.NewError("guard/signature/sign")

type request struct {
	SigningMethod string `mapstructure:"signing_method"`
	Secret        string `mapstructure:"secret"`
	Data          string `mapstructure:"data"`
}

func (s *sign) error(v ...any) error {
	return exception.Throw(v...)
}

func (s *sign) getRequest(input any) (*request, error) {
	req := request{}
	if err := plugin.DecodeInput(input, &req); err != nil {
		return nil, s.error("decode input error", err)
	}
	if req.SigningMethod == "" {
		return nil, s.error("signing_method is required")
	}
	if req.Secret == "" {
		return nil, s.error("secret is required")
	}
	return &req, nil
}

// Run Note
func (s *sign) Run(ctx context.Context, input any) (any, error) {
	req, err := s.getRequest(input)
	if err != nil {
		return nil, err
	}
	t := time.Now()
	sign, err := guard.GenerateSignature(t, req.SigningMethod, []byte(req.Data), []byte(req.Secret))
	if err != nil {
		return nil, s.error("generate error", err)
	}
	return map[string]any{
		"signature": sign,
		"timestamp": t.Unix(),
	}, nil
}
