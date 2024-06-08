package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func NewHTTPClient(method string, url string, timeOut time.Duration,
	body io.Reader, headers map[string]string,
) (*http.Response, error) {
	client := &http.Client{
		Timeout: timeOut,
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create newRequest, method: %s, url: %s, Reason: %v", method, url, err)
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	v, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed in client http request, method: %s, url: %s, Reason: %v", method, url, err)
	}
	return v, nil
}
