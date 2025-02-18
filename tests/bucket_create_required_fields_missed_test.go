package tests

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/suite"
	language "gitlab.snapp.ir/platform/snapp_object_store/langs/en"
	"net/http"
	"testing"
)

func TestBucketCreateMissingRequiredFieldsTestSuite(t *testing.T) {
	suite.Run(t, new(BucketCreateMissingRequiredFieldsTestSuite))
}

type BucketCreateMissingRequiredFieldsTestSuite struct {
	BaseTestSuite
}

func (b *BucketCreateMissingRequiredFieldsTestSuite) TestBucketCreateMissingRequiredFieldsShouldFail() {
	path := "api/bucket/create"
	body := []byte("")
	request, err := http.NewRequest(http.MethodPost, FetchURL(b.server.Config.ServerConfigs, path), bytes.NewBuffer(body))
	b.Equal(err, nil)

	client := &http.Client{}
	res, errCall := client.Do(request)
	b.Equal(errCall, nil)

	msgErr := &HttpMessageError{}
	errDecode := json.NewDecoder(res.Body).Decode(msgErr)
	b.Equal(errDecode, nil)
	b.Equal(res.StatusCode, http.StatusBadRequest)
	b.Equal(msgErr.Message, language.AuthKeyNotProvided)
}
