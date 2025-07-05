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

// Run Note
func Run(ctx context.Context, input interface{}) (interface{}, error) {
	req, ok := input.(*Request)
	if !ok {
		return nil, errors.New("plugin:net/proxy input is not a Request type")
	}
	b, err := json.Marshal(req.Body)
	if err != nil {
		return nil, errors.New("plugin:net/proxy marshal body error: " + err.Error())
	}
	r, err := http.NewRequestWithContext(ctx, req.Method, req.URL(), bytes.NewReader(b))
	if err != nil {
		return nil, errors.New("plugin:net/proxy create request error: " + err.Error())
	}
	for k, v := range req.Header {
		r.Header.Set(k, v)
	}
	d, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, errors.New("plugin:net/proxy send http request error: " + err.Error())
	}
	defer d.Body.Close()
	resp := &Response{
		StatusCode: d.StatusCode,
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
			return nil, errors.New("plugin:net/proxy unmarshal response body error: " + err.Error())
		}
	}
	return resp, nil
}
