package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/aionbit/ixo-plugins/plugin"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	Host   string            `mapstructure:"host"`
	Method string            `mapstructure:"method"`
	Path   string            `mapstructure:"path"`
	Query  map[string]string `mapstructure:"query"`
	Header map[string]string `mapstructure:"header"`
	Body   any               `mapstructure:"body"`
}

func (r *Request) URL() string {
	path := r.Path
	if len(r.Query) > 0 {
		var queries = url.Values{}
		for k, v := range r.Query {
			queries.Set(k, v)
		}
		if strings.Contains(path, "?") {
			path += "&" + queries.Encode()
		} else {
			path += "?" + queries.Encode()
		}
	}
	return "http://" + r.Host + path
}

type Response struct {
	StatusCode int
	Header     map[string]string
	Body       any
}

func errorResponse(err error) *Response {
	if err == nil {
		err = errors.New("unknown error")
	}
	return &Response{
		StatusCode: http.StatusInternalServerError,
		Body: map[string]string{
			"error": throw(err).Error(),
		},
	}
}

func throw(err error) error {
	if err == nil {
		return nil
	}
	return errors.New("plugin:net/proxy " + err.Error())
}

var PluginInstance proxy

type proxy struct {
}

// Run Note
func (p *proxy) Run(ctx context.Context, input any) (any, error) {
	var req = &Request{}
	if err := plugin.DecodeInput(input, &req); err != nil {
		return nil, throw(err)
	}
	buf := bytes.NewBuffer(nil)
	switch {
	case strings.Contains(req.Header["Content-Type"], "application/json"):
		if err := json.NewEncoder(buf).Encode(req.Body); err != nil {
			return nil, throw(err)
		}
	case strings.Contains(req.Header["Content-Type"], "application/x-yaml"):
		if err := yaml.NewEncoder(buf).Encode(req.Body); err != nil {
			return nil, throw(err)
		}
	}
	r, err := http.NewRequestWithContext(ctx, req.Method, req.URL(), buf)
	if err != nil {
		return nil, throw(err)
	}
	for k, v := range req.Header {
		r.Header.Set(k, v)
	}
	d, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, throw(err)
	}
	defer d.Body.Close()
	resp := &Response{
		StatusCode: d.StatusCode,
		Header:     make(map[string]string),
	}
	for k, v := range d.Header {
		if len(v) > 0 {
			resp.Header[k] = v[0] // 只取第一个值
		}
	}
	body, err := io.ReadAll(d.Body)
	if err != nil {
		return nil, throw(err)
	}
	if len(body) > 0 {
		contentType := d.Header.Get("Content-Type")
		switch {
		case strings.Contains(contentType, "application/json"):
			if err := json.NewEncoder(buf).Encode(req.Body); err != nil {
				return nil, throw(err)
			}
		case strings.Contains(contentType, "application/x-yaml"):
			if err := yaml.NewEncoder(buf).Encode(req.Body); err != nil {
				return nil, throw(err)
			}
		default:
			return nil, throw(errors.New("unsupported response header Content-Type: " + contentType))
		}
	}
	return plugin.EncodeOutput(resp)
}
