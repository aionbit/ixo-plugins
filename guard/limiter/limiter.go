package main

import (
	"context"
	"errors"
	"github.com/aionbit/ixo-plugins/plugin"
	"github.com/bluele/gcache"
	"golang.org/x/time/rate"
	"math"
	"time"
)

var limiters = gcache.New(1000000).LRU().Build()

var PluginInstance limiter

type limiter struct {
}

var exception = plugin.NewError("guard/limiter")

type request struct {
	Name       string     `mapstructure:"name"`
	Limit      rate.Limit `mapstructure:"limit"`
	Burst      int        `mapstructure:"burst"`
	Expiration string     `mapstructure:"expiration"`
	Wait       bool       `mapstructure:"wait"`
	expiration time.Duration
}

func (lim *limiter) error(v ...any) error {
	return exception.Throw(v...)
}

func (lim *limiter) getRequest(input any) (*request, error) {
	var req request
	if err := plugin.DecodeInput(input, &req); err != nil {
		return nil, lim.error("decode input error", err)
	}
	if req.Name == "" {
		return nil, lim.error("name is required")
	}
	if req.Expiration == "" {
		return nil, lim.error("expiration is required")
	} else {
		d, err := time.ParseDuration(req.Expiration)
		if err != nil || d <= 0 {
			return nil, lim.error("invalid expiration")
		}
		req.expiration = d
	}
	if math.IsNaN(float64(req.Limit)) {
		return nil, lim.error("limit is NaN")
	}
	return &req, nil
}

// Run Note
func (lim *limiter) Run(ctx context.Context, input any) (any, error) {
	req, err := lim.getRequest(input)
	if err != nil {
		return nil, err
	}
	if req.Burst < 0 || req.Limit < 0 {
		return nil, nil
	}
	v, err := limiters.GetIFPresent(req.Name)
	if err != nil {
		if !errors.Is(err, gcache.KeyNotFoundError) {
			return nil, lim.error("get limiter from cache error", err)
		}
		v = rate.NewLimiter(req.Limit, req.Burst)
		if err = limiters.SetWithExpire(req.Name, v, req.expiration); err != nil {
			return nil, lim.error("set limiter to cache error", err)
		}
	}
	if req.Wait {
		if err = v.(*rate.Limiter).WaitN(ctx, 1); err != nil {
			return nil, lim.error("wait token error", err)
		}
	} else {
		if !v.(*rate.Limiter).Allow() {
			return nil, lim.error("too many requests")
		}
	}
	return nil, nil
}
