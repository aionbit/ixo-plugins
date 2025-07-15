package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

// Input 输入
func Input() interface{} {
	return &Request{}
}

func errorResponse(err error) *Response {
	if err == nil {
		err = errors.New("unknown error")
	}
	return &Response{
		StatusCode: http.StatusInternalServerError,
		Body: map[string]string{
			"error": "plugin:net/proxy " + err.Error(),
		},
	}
}

// Run Note
func Run(ctx context.Context, input interface{}) (interface{}, error) {
	req, ok := input.(*Request)
	if !ok {
		return errorResponse(errors.New("input is not a Request type")), nil
	}
	buf := bytes.NewBuffer(nil)
	switch {
	case strings.Contains(req.Header["Content-Type"], "application/json"):
		if err := json.NewEncoder(buf).Encode(req.Body); err != nil {
			return errorResponse(err), nil
		}
	case strings.Contains(req.Header["Content-Type"], "application/x-yaml"):
		if err := yaml.NewEncoder(buf).Encode(req.Body); err != nil {
			return errorResponse(err), nil
		}
	}
	r, err := http.NewRequestWithContext(ctx, req.Method, req.URL(), buf)
	if err != nil {
		return errorResponse(err), nil
	}
	for k, v := range req.Header {
		r.Header.Set(k, v)
	}
	d, err := http.DefaultClient.Do(r)
	if err != nil {
		return errorResponse(err), nil
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
		return nil, err
	}
	if len(body) > 0 {
		contentType := d.Header.Get("Content-Type")
		switch {
		case strings.Contains(contentType, "application/json"):
			if err := json.NewEncoder(buf).Encode(req.Body); err != nil {
				return errorResponse(err), nil
			}
		case strings.Contains(contentType, "application/x-yaml"):
			if err := yaml.NewEncoder(buf).Encode(req.Body); err != nil {
				return errorResponse(err), nil
			}
		default:
			return errorResponse(errors.New("unsupported response header Content-Type: " + contentType)), nil
		}
	}
	return resp, nil
}
