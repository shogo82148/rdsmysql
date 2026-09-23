package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/shogo82148/rdsmysql/v2/internal/config"
)

func main() {
	conf, err := config.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(run(conf))
}

func run(c *config.Config) int {
	if c.Version {
		showVersion()
		return 0
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dir, err := os.MkdirTemp("", "rdsmysql-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir) //nolint:errcheck // clean up

	cfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = config.Generate(ctx, cfg, dir, c)
	if err != nil {
		log.Fatal(err)
	}

	mysql, err := exec.LookPath("mysqldump")
	if err != nil {
		log.Fatal(err)
	}

	args := append([]string{fmt.Sprintf("--defaults-extra-file=%s", filepath.Join(dir, "my.conf"))}, c.Args...)
	cmd := exec.Command(mysql, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	done := make(chan struct{})

	// transfer signals.
	wg.Go(func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		for {
			select {
			case s := <-sig:
				_ = cmd.Process.Signal(s)
			case <-done:
				return
			}
		}
	})

	// password rotation
	wg.Go(func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := config.Generate(ctx, cfg, dir, c); err != nil {
					log.Printf("failed to re-generate auth token: %v", err)
				}
			case <-done:
				return
			}
		}
	})

	_ = cmd.Wait()
	close(done)
	wg.Wait()
	return cmd.ProcessState.ExitCode()
}
