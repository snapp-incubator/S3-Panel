package tests

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

type ServerTestSuite struct {
	BaseTestSuite
}
