package tests

import (
	"fmt"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/infra/config"
	"net/http"
	"strings"
)

func FetchURL(conf config.ServerConfig, path string) string {
	return fmt.Sprintf("http://%s:%s/%s", conf.Address, conf.Port, path)
}

func isServerAuthEnabled(s string) bool {
	// consider it would be disabled by default
	if strings.ToLower(s) == "true" {
		return true
	}
	return false
}

func AddAuthorizationHeader(req *http.Request, conf config.ServerConfig) {
	if isServerAuthEnabled(conf.AuthEnabled) {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", conf.AuthToken))
	}
}

func AddContentTypeHeader(req *http.Request) {
	req.Header.Add("Content-Type", "application/json")
}

func DoHttpClientRequest(r *http.Request) (*http.Response, error) {
	client := &http.Client{}
	res, errCall := client.Do(r)
	defer res.Body.Close()
	return res, errCall
}
