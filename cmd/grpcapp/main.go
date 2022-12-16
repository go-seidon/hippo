package main

import (
	"context"

	"github.com/go-seidon/hippo/internal/app"
	"github.com/go-seidon/hippo/internal/grpcapp"
)

func main() {
	config, err := app.NewDefaultConfig()
	if err != nil {
		panic(err)
	}

	app, err := grpcapp.NewGrpcApp(
		grpcapp.WithConfig(config),
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
