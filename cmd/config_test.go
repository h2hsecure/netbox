package cmd_test

import (
	"testing"

	"github.com/h2hsecure/netbox/cmd"
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
