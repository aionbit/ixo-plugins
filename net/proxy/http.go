package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	Host   string            `json:"host"`
	Method string            `json:"method"`
	Path   string            `json:"path"`
	Query  map[string]string `json:"query"`
	Header map[string]string `json:"header"`
	Body   any               `json:"body"`
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
	StatusCode int               `json:"status_code"`
	Header     map[string]string `json:"header"`
	Body       any               `json:"body"`
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
		Header:     map[string]string{"Content-Type": "application/json"},
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
	b, err := json.Marshal(req.Body)
	if err != nil {
		return errorResponse(err), nil
	}
	r, err := http.NewRequestWithContext(ctx, req.Method, req.URL(), bytes.NewReader(b))
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
		if err := json.Unmarshal(body, &resp.Body); err != nil {
			return errorResponse(err), nil
		}
	}
	return resp, nil
}
