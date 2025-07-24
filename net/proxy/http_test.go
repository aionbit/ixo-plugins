package main

import (
	"context"
	"io"
	"testing"
)

func Test_proxy_Run(t *testing.T) {
	req := map[string]any{
		"host":   "ixo.io",
		"method": "GET",
		"path":   "/api/user/fav/news",
		"query": map[string]any{
			"user_id": "1",
		},
		"header": map[string]any{
			"Gw-Namespace": "fake",
		},
	}
	resp, err := PluginInstance.Run(context.Background(), req)
	if err != nil {
		t.Error(err)
		return
	}
	body := resp.(map[string]any)["body"].(io.ReadCloser)
	data, err := io.ReadAll(body)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(data))
}
