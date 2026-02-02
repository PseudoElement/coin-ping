package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

func Get(urlStr string, successResp any, errResp any, params [][2]string, headers [][2]string) (err error, success bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if len(params) > 0 {
		p := url.Values{}
		for _, pair := range params {
			p.Add(pair[0], pair[1])
		}

		urlStr += "?" + p.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		panic(err)
	}

	if len(headers) > 0 {
		for _, pair := range headers {
			req.Header.Set(pair[0], pair[1])
		}
	}

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err, false
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		err = json.Unmarshal(body, &errResp)
		return err, false
	}

	err = json.Unmarshal(body, &successResp)

	return err, err == nil
}
