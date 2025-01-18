package cmd

import (
	"context"
	"fmt"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/infra/logger"
	"log"
	"os"
	"runtime/debug"

	"github.com/urfave/cli/v3"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/infra/config"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/platform/server"
)

func Execute() {
	var configPath string

	cmd := &cli.Command{
		Name:        "snapp-object-store",
		Description: "This is the backend of S3 object storage panel on snappcloud unified panel",
		Commands: []*cli.Command{
			{
				Name:  "snapp-object-store",
				Usage: "run snapp-object-store backend",
				Action: func(_ context.Context, _ *cli.Command) error {
					cfg := config.Provide(configPath)
					loggerObj := logger.Provide(cfg.LoggerConfigs)
					err := server.StartServer(cfg, loggerObj)
					if err != nil {
						return err
					}
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "configPath",
						Value:       "./config.yaml",
						Usage:       "Path to config file",
						Destination: &configPath,
						OnlyOnce:    false,
					},
				},
			},
		},
		Version: func() string {
			revision := ""
			timestamp := ""
			modified := ""

			if info, ok := debug.ReadBuildInfo(); ok {
				for _, setting := range info.Settings {
					switch setting.Key {
					case "vcs.revision":
						revision = setting.Value
					case "vcs.time":
						timestamp = setting.Value
					case "vcs.modified":
						modified = setting.Value
					}
				}
			}

			if revision == "" {
				return ""
			}

			if modified == "true" {
				return fmt.Sprintf("%s (%s) [dirty]", revision, timestamp)
			}

			return fmt.Sprintf("%s (%s)", revision, timestamp)
		}(),
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
