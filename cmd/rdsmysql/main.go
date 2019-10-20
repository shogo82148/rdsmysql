package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/shogo82148/rdsmysql/internal/config"
)

func main() {
	conf, err := config.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(run(conf))
}

func run(c *config.Config) int {
	dir, err := ioutil.TempDir("", "rdsmysql-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	conf := aws.NewConfig()
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            *conf,
		SharedConfigState: session.SharedConfigEnable,
	}))

	err = config.Generate(sess, dir, c)
	if err != nil {
		log.Fatal(err)
	}

	mysql, err := exec.LookPath("mysql")
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
	wg.Add(1)
	go func() {
		defer wg.Done()
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		for {
			select {
			case s := <-sig:
				cmd.Process.Signal(s)
			case <-done:
				return
			}
		}
	}()

	// password rotation
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				config.Generate(sess, dir, c)
			case <-done:
				return
			}
		}
	}()

	_ = cmd.Wait()
	close(done)
	wg.Wait()
	return cmd.ProcessState.ExitCode()
}
