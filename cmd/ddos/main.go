package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"git.h2hsecure.com/ddos/waf/internal/repository/cache"
	"git.h2hsecure.com/ddos/waf/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	_, cancel := context.WithCancel(context.Background())

	go func() {
		cache := cache.NewMemcache("localhost:11211")
		server.CreateHttpServer("8081", cache)
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Blocking, press ctrl+c to continue...")
	<-done
	cancel()
}
