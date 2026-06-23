package tests

import (
	"fmt"
	"net/http"
	"strings"

	"gitlab.snapp.ir/platform/s3-panel/internal/config"
)

func FetchURL(conf config.ServerConfig, path string) string {
	return fmt.Sprintf("http://%s:%s/%s", conf.Address, conf.Port, path)
}

func isServerAuthEnabled(s string) bool {
	// auth is disabled by default and only enabled when the value is exactly "true"
	return strings.EqualFold(strings.TrimSpace(s), "true")
}

func AddAuthorizationHeader(req *http.Request, conf config.ServerConfig) {
	if isServerAuthEnabled(conf.AuthEnabled) {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", conf.AuthToken))
	}
}

func AddContentTypeHeader(req *http.Request) {
	req.Header.Add("Content-Type", "application/json")
}

// DoHttpClientRequest performs the request and returns the response. The caller
// is responsible for closing the response body.
func DoHttpClientRequest(r *http.Request) (*http.Response, error) {
	client := &http.Client{}
	res, errCall := client.Do(r)
	if errCall != nil {
		return nil, errCall
	}
	return res, nil
}
