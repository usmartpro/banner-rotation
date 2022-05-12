package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/usmartpro/banner-rotation/internal/app"
	"github.com/usmartpro/banner-rotation/internal/config"
	"github.com/usmartpro/banner-rotation/internal/logger"
	"github.com/usmartpro/banner-rotation/internal/mq"
	internalhttp "github.com/usmartpro/banner-rotation/internal/server/http"
	"github.com/usmartpro/banner-rotation/internal/storage"
)

func main() {
	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	configuration, err := config.LoadConfiguration()
	if err != nil {
		log.Fatalf("Error read configuration: %s", err)
	}
	logg, err := logger.New(configuration.Logger)
	if err != nil {
		log.Println("error create logger: " + err.Error())
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	rabbitClient, err := mq.NewRabbit(
		ctx,
		configuration.Rabbit.Dsn,
		configuration.Rabbit.Exchange,
		configuration.Rabbit.Queue,
		logg)
	if err != nil {
		cancel()
		log.Fatalf("error create rabbit client: %s", err) //nolint:gocritic
	}

	storageConf := storage.New(ctx, configuration.Storage.Dsn).Connect(ctx)
	bannerRotation := app.New(logg, storageConf, rabbitClient)

	// HTTP
	server := internalhttp.NewServer(logg, bannerRotation, configuration.HTTP.Host, configuration.HTTP.Port)

	go func() {
		if err := server.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
		}
	}()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("banner rotation is running...")

	<-ctx.Done()
}
