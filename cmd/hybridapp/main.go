package main

import (
	"context"

	"github.com/go-seidon/hippo/internal/app"
	"github.com/go-seidon/hippo/internal/grpcapp"
	"github.com/go-seidon/hippo/internal/restapp"
)

func main() {
	config, err := app.NewDefaultConfig()
	if err != nil {
		panic(err)
	}

	restApp, err := restapp.NewRestApp(
		restapp.WithConfig(config),
	)
	if err != nil {
		panic(err)
	}

	grpcApp, err := grpcapp.NewGrpcApp(
		grpcapp.WithConfig(config),
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
