package tests

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/infra/config"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/infra/logger"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/platform/server"
)

var conf config.Config
var confPath = "./../configs/configs-test.yaml"

func init() {
	cfg := config.Provide(confPath)
	conf = cfg
}

type BaseTestSuite struct {
	suite.Suite
	server  *server.Server
	context context.Context
	cancel  context.CancelFunc
	*require.Assertions
}

func (s *BaseTestSuite) SetupTest() {
	ctx, cancel := context.WithCancel(context.Background())
	s.context = ctx
	s.cancel = cancel
	s.Assertions = s.Require()
	testServer, err := server.NewServer(ctx, cancel, conf, logger.Provide(conf.LoggerConfigs))
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

func (s *BaseTestSuite) TearDownTest() {
	s.cancel()

	if err := s.server.ShutDown(); err != nil {
		s.FailNowf("failed to close server with err : ", err.Error())
	}
}
