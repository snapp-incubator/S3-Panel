package tests

import (
	"context"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.snapp.ir/platform/s3-panel/internal/api"
	"gitlab.snapp.ir/platform/s3-panel/internal/config"
	"gitlab.snapp.ir/platform/s3-panel/internal/logging"
)

var conf config.Config
var confPath = "./../configs/test-config.yaml"

func init() {
	cfg := config.Provide(confPath)
	conf = cfg
}

type HttpMessageError struct {
	Message string `json:"message"`
}

type BaseTestSuite struct {
	suite.Suite
	server  *api.Server
	context context.Context
	cancel  context.CancelFunc
	*require.Assertions
}

func (s *BaseTestSuite) SetupSuite() {
	ctx, cancel := context.WithCancel(context.Background())
	s.context = ctx
	s.cancel = cancel
	s.Assertions = s.Require()
	testServer, err := api.NewServer(ctx, cancel, conf, logging.Provide(conf.Logger))
	if err != nil {
		s.FailNowf("failed to setup application with err : ", err.Error())
	}

	s.server = testServer
	go func() {
		errStart := testServer.Start()
		if errStart != nil {
			s.FailNow("Failed to start router", errStart.Error())
		}
	}()
}

func (s *BaseTestSuite) TearDownSuite() {
	s.cancel()

	if err := s.server.ShutDown(); err != nil {
		s.FailNowf("failed to close server with err : ", err.Error())
	}
}
