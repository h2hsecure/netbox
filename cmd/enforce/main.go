package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"git.h2hsecure.com/ddos/waf/internal/core/domain"
	"git.h2hsecure.com/ddos/waf/internal/repository/cache"
	"git.h2hsecure.com/ddos/waf/internal/repository/fsm"
	"git.h2hsecure.com/ddos/waf/internal/server"
	"git.h2hsecure.com/ddos/waf/internal/server/handler"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).
		With().
		Str("context", "encoder").
		Logger()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	errChan := make(chan error)

	go func() {
		mux := server.CreatePromServer()

		if err := http.ListenAndServe(os.Getenv("PROM_LISTEN"), mux); err != nil {
			errChan <- err
		}
	}()

	go func() {

		cache, err := cache.NewMemcache(os.Getenv("CACHE_SOCK"))
		if err != nil {
			errChan <- err
			return
		}

		machine := fsm.NewStateMachine(cache)

		clusterAddress, err := domain.ParseAddress(os.Getenv("CLUSTER_STR"))

		if err != nil {
			errChan <- fmt.Errorf("parse address: %w", err)
		}

		if len(clusterAddress) == 0 {
			errChan <- fmt.Errorf("no address found for grpc: %s", os.Getenv("CLUSTER_STR"))
		}

		myAddress, err := domain.ParseAddress(os.Getenv("MY_ADDRESS"))

		if err != nil {
			errChan <- fmt.Errorf("parse address: %w", err)
		}

		if len(clusterAddress) == 0 {
			errChan <- fmt.Errorf("no address found for grpc: %s", os.Getenv("myAddress"))
		}

		raft, err := server.NewRaft(myAddress[0], clusterAddress, machine)

		if err != nil {
			errChan <- err
			return
		}

		grpcServerHandler := handler.NewGrpcHandler(raft)

		log.Info().Msg("Enforce started, press ctrl+c to break it..")

		if err := server.CreateGrpcServer(os.Getenv("GRPC_SERVER_PORT"), grpcServerHandler); err != nil {
			errChan <- err
			return
		}

	}()

	select {
	case <-done:
	case err := <-errChan:
		log.Err(err).Msg("startup")
	}
}
