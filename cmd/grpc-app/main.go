package main

import (
	"context"

	"github.com/go-seidon/hippo/internal/app"
	grpc_app "github.com/go-seidon/hippo/internal/grpc-app"
)

func main() {
	config, err := app.NewDefaultConfig()
	if err != nil {
		panic(err)
	}

	app, err := grpc_app.NewGrpcApp(
		grpc_app.WithConfig(config),
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
