package handler_test

import (
	"fmt"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/h2hsecure/netbox/cmd"
	"github.com/h2hsecure/netbox/internal/server/handler"
)

func Test_Config_Handler(t *testing.T) {

	hostname := "test_hostname"
	referer := "test_referer"

	os.Setenv("DOMAIN", hostname)
	os.Setenv("SYSTEM_ID", "Test123")

	gin.SetMode(gin.ReleaseMode)
	mux := gin.New()
	cfg, err := cmd.CurrentConfig()

	fmt.Printf("dump config: %v", cfg)

	if err != nil {
		t.Fatalf("it shouldn't return an error: %v", err)
	}

	err = handler.NewConfigHandler(mux, cfg)

	if err != nil {
		t.Fatalf("it shouldn't return an error: %v", err)
	}

	listener, err := net.Listen("tcp", ":64001")
	if err != nil {
		t.Fatalf("it shouldn't return an error: %v", err)
	}

	go func() {
		mux.RunListener(listener)
	}()

	response, err := resty.New().R().
		SetHeader("Referer", referer).
		Get("http://localhost:64001/ntb_dds")

	if err != nil || response.IsError() {
		t.Fatalf("it shouldn't return an error: %v, %v", err, response)
	}

	if !strings.Contains(string(response.Body()), hostname) {
		t.Fatalf("don't find related string=%s inside=%s", hostname, string(response.Body()))
	}

	if !strings.Contains(string(response.Body()), referer) {
		t.Fatalf("don't find related string=%s inside=%s", referer, string(response.Body()))
	}

	if !strings.Contains(string(response.Body()), "Verify you are human") {
		t.Fatalf("don't find related string=%s inside=%s", "Verify you are human", string(response.Body()))
	}

	response, err = resty.New().R().
		SetHeader("Referer", referer).
		SetHeader("Accept-Language", "nl-NL").
		Get("http://localhost:64001/ntb_dds")

	if err != nil || response.IsError() {
		t.Fatalf("it shouldn't return an error: %v, %v", err, response)
	}

	if !strings.Contains(string(response.Body()), "Verifieer dat u een mens") {
		t.Fatalf("don't find related string=%s inside=%s", "Verifieer dat u een mens", string(response.Body()))
	}

	listener.Close()
}
