package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-seidon/hippo/internal/app"
	"github.com/go-seidon/hippo/internal/config"
	grpc_app "github.com/go-seidon/hippo/internal/grpc-app"
)

func main() {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "local"
	}

	appConfig := &app.Config{AppEnv: appEnv}

	cfgFileName := fmt.Sprintf("config/%s.toml", appConfig.AppEnv)
	tomlConfig, err := config.NewViperConfig(
		config.WithFileName(cfgFileName),
	)
	if err != nil {
		panic(err)
	}

	err = tomlConfig.LoadConfig()
	if err != nil {
		panic(err)
	}

	err = tomlConfig.ParseConfig(appConfig)
	if err != nil {
		panic(err)
	}

	app, err := grpc_app.NewGrpcApp(
		grpc_app.WithConfig(appConfig),
	)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	err = app.Run(ctx)
	if err != nil {
		panic(err)
	}
}
