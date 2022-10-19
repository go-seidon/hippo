package main

import (
	"context"

	"github.com/go-seidon/hippo/internal/app"
	grpc_app "github.com/go-seidon/hippo/internal/grpc-app"
	rest_app "github.com/go-seidon/hippo/internal/rest-app"
)

func main() {
	config, err := app.NewDefaultConfig()
	if err != nil {
		panic(err)
	}

	restApp, err := rest_app.NewRestApp(
		rest_app.WithConfig(config),
	)
	if err != nil {
		panic(err)
	}

	grpcApp, err := grpc_app.NewGrpcApp(
		grpc_app.WithConfig(config),
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
