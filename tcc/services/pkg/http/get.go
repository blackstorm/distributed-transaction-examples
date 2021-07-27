package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

func Post(ctx context.Context, url string, data interface{}) (*http.Response, error) {
	var reader io.Reader

	if data != nil {
		v, _ := json.Marshal(data)
		reader = bytes.NewReader(v)
	}

	if req, err := http.NewRequestWithContext(ctx, "Post", url, reader); err != nil {
		return nil, err
	} else {
		return http.DefaultClient.Do(req)
	}
}
