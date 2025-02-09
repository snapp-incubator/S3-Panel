package tests

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestUserCreatesExcessBucketTestSuite(t *testing.T) {
	suite.Run(t, new(UserCreatesExcessBucketTestSuite))
}

type UserCreatesExcessBucketTestSuite struct {
	BaseTestSuite
}

func (u *UserCreatesExcessBucketTestSuite) TestUserShouldFailToCreateExcessBucket() {}
