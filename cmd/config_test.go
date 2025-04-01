package cmd_test

import (
	"testing"

	"git.h2hsecure.com/ddos/waf/cmd"
)

func TestConfigurationParse(t *testing.T) {
	cfg, err := cmd.CurrentConfig()

	if err != nil {
		t.Fatalf("error shouldn't returned")
	}

	if cfg.User.CounterFreq != 100 {
		t.Fatalf("config default test failed: expected: %d result: %d", 100, cfg.User.CounterFreq)
	}
}
