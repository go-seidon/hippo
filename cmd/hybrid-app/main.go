package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-seidon/hippo/internal/app"
	"github.com/go-seidon/hippo/internal/config"
	grpc_app "github.com/go-seidon/hippo/internal/grpc-app"
	rest_app "github.com/go-seidon/hippo/internal/rest-app"
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

	restApp, err := rest_app.NewRestApp(
		rest_app.WithConfig(appConfig),
	)
	if err != nil {
		panic(err)
	}

	grpcApp, err := grpc_app.NewGrpcApp(
		grpc_app.WithConfig(appConfig),
	)
	if err != nil {
		panic(err)
	}

	listener := make(chan error, 2)

	go func() {
		ctx := context.Background()
		listener <- restApp.Run(ctx)
	}()

	go func() {
		ctx := context.Background()
		listener <- grpcApp.Run(ctx)
	}()

	err = <-listener
	if err != nil {
		panic(err)
	}
}
