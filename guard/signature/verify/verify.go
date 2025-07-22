package main

import (
	"context"
	"github.com/aionbit/ixo-plugins/guard"
	"github.com/aionbit/ixo-plugins/plugin"
	"time"
)

var PluginInstance verify

type verify struct {
}

var exception = plugin.NewError("guard/signature/verify")

type request struct {
	Signature     string `mapstructure:"signature"`
	Timestamp     int    `mapstructure:"timestamp"`
	SigningMethod string `mapstructure:"signing_method"`
	Secret        string `mapstructure:"secret"`
	Data          string `mapstructure:"data"`
}

func (ver *verify) error(v ...any) error {
	return exception.Throw(v...)
}

func (ver *verify) getRequest(input any) (*request, error) {
	req := request{}
	if err := plugin.DecodeInput(input, &req); err != nil {
		return nil, ver.error("decode input error", err)
	}
	if req.Signature == "" {
		return nil, ver.error("signature is required")
	}
	if req.Timestamp <= 0 {
		return nil, ver.error("invalid timestamp")
	}
	if req.SigningMethod == "" {
		return nil, ver.error("signing_method is required")
	}
	if req.Secret == "" {
		return nil, ver.error("secret is required")
	}
	return &req, nil
}

// Run Note
func (ver *verify) Run(ctx context.Context, input any) (any, error) {
	req, err := ver.getRequest(input)
	if err != nil {
		return nil, err
	}
	tm := time.Unix(int64(req.Timestamp), 0)
	if time.Since(tm) > 5*time.Minute || time.Until(tm) > 5*time.Minute {
		return nil, ver.error("expired timestamp")
	}
	sign, err := guard.GenerateSignature(tm, req.SigningMethod, []byte(req.Data), []byte(req.Secret))
	if err != nil {
		return nil, ver.error("generate error", err)
	}
	if sign != req.Signature {
		return nil, ver.error("invalid signature")
	}
	return nil, nil
}
