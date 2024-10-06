package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

		raft, err := server.NewRaft(os.Getenv("CLUSTER_ID"), os.Getenv("MY_ADDRESS"), machine)

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
