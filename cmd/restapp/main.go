package main

import (
	"context"

	"github.com/go-seidon/hippo/internal/app"
	"github.com/go-seidon/hippo/internal/restapp"
)

func main() {
	config, err := app.NewDefaultConfig()
	if err != nil {
		panic(err)
	}

	app, err := restapp.NewRestApp(
		restapp.WithConfig(config),
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
