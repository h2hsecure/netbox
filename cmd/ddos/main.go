package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"

	"git.h2hsecure.com/ddos/waf/internal/core/domain"
	"git.h2hsecure.com/ddos/waf/internal/repository/cache"
	"git.h2hsecure.com/ddos/waf/internal/repository/grpc"
	"git.h2hsecure.com/ddos/waf/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	logFileName := path.Join(os.Getenv("LOG_DIR"), "ddos.log")
	f, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(fmt.Errorf("unable to create log file: %s", logFileName))
	}

	defer f.Close()

	if err := os.Chmod(logFileName, os.FileMode(0644)); err != nil {
		panic(fmt.Errorf("chmod failed for log file: %s", logFileName))
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: &lumberjack.Logger{
		Filename:   logFileName,
		MaxBackups: 10, // files
		MaxSize:    5,  // megabytes
		MaxAge:     10, // days
		Compress:   true,
	}}).
		With().
		Str("context", "ddos").
		Logger()

	_, cancel := context.WithCancel(context.Background())

	defer cancel()

	errChan := make(chan error)
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		cache, err := cache.NewMemcache(os.Getenv("CACHE_SOCK"))
		if err != nil {
			errChan <- err
			return
		}

		address, err := domain.ParseAddress(os.Getenv("CLUSTER_STR"))

		if err != nil {
			errChan <- fmt.Errorf("parse address: %w", err)
		}

		if len(address) == 0 {
			errChan <- fmt.Errorf("no address found for grpc: %s", os.Getenv("CLUSTER_STR"))
		}

		mq, err := grpc.NewEnforceClient(address)
		if err != nil {
			errChan <- err
			return
		}

		engine := server.CreateHttpServer(cache, mq)

		syscall.Unlink(os.Getenv("INTERNAL_SOCK"))

		listener, err := net.Listen("unix", os.Getenv("INTERNAL_SOCK"))
		if err != nil {
			log.Err(err).Msg("listen socket")
			errChan <- err
			return
		}

		if err := os.Chown(os.Getenv("INTERNAL_SOCK"), 101, 101); err != nil {
			log.Err(err).Msg("chown socket")
			errChan <- err
			return
		}

		if err := os.Chmod(os.Getenv("INTERNAL_SOCK"), 0644); err != nil {
			log.Err(err).Msg("chown socket")
			errChan <- err
			return
		}

		log.Info().Msgf("Server is running on port %v\n", listener)
		if err := http.Serve(listener, engine); err != nil {
			errChan <- err
		}
	}()

	log.Info().Msg("Ddos backend started, press ctrl+c to break it...")
	select {
	case <-done:
		log.Info().Msg("signal recieved")
	case err := <-errChan:
		log.Err(err).Send()
		os.Exit(1)
	}
}
