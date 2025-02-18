package tests

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestBucketCreateTestSuite(t *testing.T) {
	suite.Run(t, new(BucketCreateTestSuite))
}

type BucketCreateTestSuite struct {
	BaseTestSuite
}

func (u *BucketCreateTestSuite) TestBucketCreateQuotaExceedShouldFail() {}

func (u *BucketCreateTestSuite) TestBucketCreateShouldPass() {}
