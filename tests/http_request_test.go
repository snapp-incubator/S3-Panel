package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/snapp-incubator/S3-Panel/internal/health"
	"github.com/snapp-incubator/S3-Panel/internal/messages"
)

func (s *ServerTestSuite) TestHealthEndpointShouldPass() {
	path := "/health"
	body := []byte("")
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, FetchURL(s.server.Config.Server, path), bytes.NewBuffer(body))
	AddContentTypeHeader(request)
	ctx := s.server.Router.NewContext(request, recorder)
	errCall := health.HandleHealth(ctx)
	s.Nil(errCall)

	var appHealth health.ApplicationHealth
	errUnmarshal := json.NewDecoder(recorder.Body).Decode(&appHealth)
	s.Nil(errUnmarshal)
	s.Equal(appHealth.Status, messages.ApplicationHealthy)
}

func (s *ServerTestSuite) TestHTTPRequestWithNoAccessKeyShouldFail() {
	path := "api/bucket/create"
	body := []byte("")

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, FetchURL(s.server.Config.Server, path), bytes.NewBuffer(body))
	AddContentTypeHeader(request)
	AddAuthorizationHeader(request, s.server.Config.Server)
	ctx := s.server.Router.NewContext(request, recorder)
	errCall := s.server.HandleBucketCreate(ctx)
	s.Nil(errCall)
	s.Equal(recorder.Code, http.StatusBadRequest)
}
