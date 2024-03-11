package main

import (
	"testing"

	"github.com/shogo82148/rdsmysql/v2/internal/config"
	"github.com/shogo82148/rdsmysql/v2/internal/testutils"
)

func TestRun(t *testing.T) {
	testutils.Setup(t)

	cfg := &config.Config{
		User: testutils.User,
		Host: testutils.Host,
		Port: 3306,
		Args: []string{"-e", "SELECT 1"},
	}
	if got := run(cfg); got != 0 {
		t.Errorf("want 0, but got %d", got)
	}
}
