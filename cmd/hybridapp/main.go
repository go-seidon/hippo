package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

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

	go func() {
		err := restApp.Run(context.Background())
		if err != nil {
			log.Fatalf("failed running rest app %v", err)
		}
	}()

	go func() {
		err := grpcApp.Run(context.Background())
		if err != nil {
			log.Fatalf("failed running grpc app %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	listener := make(chan error, 2)
	go func() {
		listener <- restApp.Stop(ctx)
	}()
	go func() {
		listener <- grpcApp.Stop(ctx)
	}()

	var lerr error
	for i := 0; i < 2; i++ {
		err := <-listener
		if err != nil {
			lerr = err
		}
	}
	if lerr != nil {
		log.Fatalf("failed stopping app %v", lerr)
	}
}
