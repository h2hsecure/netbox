package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"git.h2hsecure.com/ddos/waf/internal/core/domain"
	"git.h2hsecure.com/ddos/waf/internal/repository/cache"
	"git.h2hsecure.com/ddos/waf/internal/repository/grpc"
	"git.h2hsecure.com/ddos/waf/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).
		With().
		Str("context", "ddos").
		Logger()

	_, cancel := context.WithCancel(context.Background())

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
	}
	cancel()
}
