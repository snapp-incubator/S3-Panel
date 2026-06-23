package logging

import (
	"log"

	"gitlab.snapp.ir/platform/s3-panel/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Provide creates a zap logger for console.
func Provide(cfg config.LoggerConfig) *zap.Logger {
	var lvl zapcore.Level
	if err := lvl.Set(cfg.Level); err != nil {
		log.Printf("cannot parse log level %s: %s", cfg.Level, err)

		lvl = zapcore.WarnLevel
	}

	zapCfg := zap.NewDevelopmentConfig()
	zapCfg.Level.SetLevel(lvl)

	logger, err := zapCfg.Build()
	if err != nil {
		log.Fatalf("logger creation failed %s", err)
	}

	return logger
}
