package tests

import (
	"bytes"
	"encoding/json"
	language "gitlab.snapp.ir/platform/snapp_object_store/langs/en"
	"net/http"
)

func (s *ServerTestSuite) TestHTTPRequestWithNoAccessKeyShouldFail() {
	path := "api/bucket/create"
	body := []byte("")
	request, err := http.NewRequest(http.MethodPost, FetchURL(s.server.Config.ServerConfigs, path), bytes.NewBuffer(body))
	s.Equal(err, nil)

	client := &http.Client{}
	res, errCall := client.Do(request)
	s.Equal(errCall, nil)

	msgErr := &HttpMessageError{}
	errDecode := json.NewDecoder(res.Body).Decode(msgErr)
	s.Equal(errDecode, nil)
	s.Equal(res.StatusCode, http.StatusBadRequest)
	s.Equal(msgErr.Message, language.AuthKeyNotProvided)
}
