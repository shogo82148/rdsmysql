package main

import (
	"testing"

	"github.com/shogo82148/rdsmysql/internal/config"
	"github.com/shogo82148/rdsmysql/internal/testutils"
)

func TestRun(t *testing.T) {
	testutils.Setup(t)

	cfg := &config.Config{
		User: testutils.User,
		Host: testutils.Host,
		Args: []string{"-e", "SELECT 1"},
	}
	if got := run(cfg); got != 0 {
		t.Errorf("want 0, but got %d", got)
	}
}
