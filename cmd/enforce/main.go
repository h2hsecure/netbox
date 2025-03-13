package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"

	"git.h2hsecure.com/ddos/waf/cmd"
	"git.h2hsecure.com/ddos/waf/internal/core/domain"
	"git.h2hsecure.com/ddos/waf/internal/repository/cache"
	"git.h2hsecure.com/ddos/waf/internal/repository/fsm"
	"git.h2hsecure.com/ddos/waf/internal/server"
	"git.h2hsecure.com/ddos/waf/internal/server/handler"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	cfg, err := cmd.CurrentConfig()

	if err != nil {
		panic(fmt.Errorf("config builder: %w", err))
	}

	logFileName := path.Join(cfg.LogDir, "enforcer.log")
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
	}}).
		With().
		Str("context", "encoder").
		Logger()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	errChan := make(chan error)

	go func() {
		mux := server.CreatePromServer()

		if err := http.ListenAndServe(cfg.PromListener, mux); err != nil {
			errChan <- err
		}
	}()

	go func() {

		cache, err := cache.NewMemcache(cfg.Cache.Sock)
		if err != nil {
			errChan <- err
			return
		}

		machine := fsm.NewStateMachine(cfg.Cache, cache)

		clusterAddress, err := domain.ParseAddress(cfg.Enforcer.ClusterStr)

		if err != nil {
			errChan <- fmt.Errorf("parse address: %w", err)
		}

		if len(clusterAddress) == 0 {
			errChan <- fmt.Errorf("no address found for grpc: %s", cfg.Enforcer.ClusterStr)
		}

		myAddress, err := domain.ParseAddress(cfg.Enforcer.MyAddress)

		if err != nil {
			errChan <- fmt.Errorf("parse address: %w", err)
		}

		if len(clusterAddress) == 0 {
			errChan <- fmt.Errorf("no address found for grpc: %s", cfg.Enforcer.MyAddress)
		}

		raft, err := server.NewRaft(myAddress[0], clusterAddress, machine)

		if err != nil {
			errChan <- err
			return
		}

		grpcServerHandler := handler.NewGrpcHandler(raft)

		log.Info().
			Str("id", myAddress[0].GetId()).
			Str("grpc", myAddress[0].GrpcAddress()).
			Msg("Enforce started")

		if err := server.CreateGrpcServer(myAddress[0], grpcServerHandler); err != nil {
			errChan <- err
			return
		}
	}()

	select {
	case <-done:
	case err := <-errChan:
		log.Err(err).Msg("startup in enforcer")
		os.Exit(1)
	}
}
