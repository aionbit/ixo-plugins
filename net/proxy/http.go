package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/aionbit/ixo-plugins/plugin"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var exception = plugin.NewError("net/proxy")

type request struct {
	Host    string            `mapstructure:"host"`
	Method  string            `mapstructure:"method"`
	Path    string            `mapstructure:"path"`
	Query   map[string]string `mapstructure:"query"`
	Header  map[string]string `mapstructure:"header"`
	Body    any               `mapstructure:"body"`
	Timeout string            `mapstructure:"timeout"`
	timeout time.Duration
}

func (r *request) URL() string {
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
	Header     map[string]any
	Body       io.ReadCloser
}

func (resp *Response) Encode() (map[string]any, error) {
	return map[string]any{
		"status_code": resp.StatusCode,
		"header":      resp.Header,
		"body":        resp.Body,
	}, nil
}

var PluginInstance proxy

type proxy struct {
}

func (p *proxy) error(v ...any) error {
	return exception.Throw(v...)
}

func (p *proxy) getRequest(input any) (*request, error) {
	var req request
	if err := plugin.DecodeInput(input, &req); err != nil {
		return nil, p.error("decode input error", err)
	}
	if req.Host == "" {
		return nil, p.error("host is required")
	}
	if req.Method == "" {
		return nil, p.error("method is required")
	}
	if req.Path == "" {
		return nil, p.error("path is required")
	}
	if req.Timeout != "" {
		d, err := time.ParseDuration(req.Timeout)
		if err != nil {
			return nil, p.error("parse timeout error", d)
		}
		if d <= 0 {
			return nil, p.error("invalid timeout")
		}
		req.timeout = d
	}
	return &req, nil
}

// Run Note
func (p *proxy) Run(ctx context.Context, input any) (any, error) {
	req, err := p.getRequest(input)
	if err != nil {
		return nil, err
	}
	var reqBody io.Reader
	if b, ok := req.Body.(io.Reader); ok {
		reqBody = b
	} else {
		buf := bytes.NewBuffer(nil)
		switch {
		case strings.Contains(req.Header["Content-Type"], "application/json"):
			if err := json.NewEncoder(buf).Encode(req.Body); err != nil {
				return nil, p.error("json encode error", err)
			}
		case strings.Contains(req.Header["Content-Type"], "application/x-yaml"):
			if err := yaml.NewEncoder(buf).Encode(req.Body); err != nil {
				return nil, p.error("yaml encode error", err)
			}
		}
		reqBody = buf
	}
	if req.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, req.timeout)
		defer cancel()
	}
	r, err := http.NewRequestWithContext(ctx, req.Method, req.URL(), reqBody)
	if err != nil {
		return nil, p.error("create new request error", err)
	}
	for k, v := range req.Header {
		r.Header.Set(k, v)
	}
	d, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, p.error("send http request error", err)
	}
	resp := &Response{
		StatusCode: d.StatusCode,
		Header:     make(map[string]any),
		Body:       newAutoCloseReadCloser(d.Body),
	}
	return resp.Encode()
}

func newAutoCloseReadCloser(rc io.ReadCloser) *autoCloseReadCloser {
	return &autoCloseReadCloser{rc: rc}
}

type autoCloseReadCloser struct {
	rc     io.ReadCloser
	closed bool
}

func (a *autoCloseReadCloser) Read(p []byte) (int, error) {
	n, err := a.rc.Read(p)
	if err == io.EOF && !a.closed {
		a.closed = true
		_ = a.rc.Close()
	}
	return n, err
}

func (a *autoCloseReadCloser) Close() error {
	if a.closed {
		return nil
	}
	a.closed = true
	_, _ = io.Copy(io.Discard, a.rc)
	return a.rc.Close()
}
