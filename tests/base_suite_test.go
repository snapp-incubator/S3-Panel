package tests

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

type ServerTestSuite struct {
	BaseTestSuite
}
