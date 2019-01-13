package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
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
	run(conf)
}

func run(c *config.Config) {
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

	// transfer signals.
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		for s := range sig {
			cmd.Process.Signal(s)
		}
	}()

	// password rotation
	go func() {
		d := 5 * time.Minute
		for range time.Tick(d) {
			func() {
				config.Generate(sess, dir, c)
			}()
		}
	}()

	if err := cmd.Wait(); err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			if s, ok := e.Sys().(syscall.WaitStatus); ok {
				os.Exit(s.ExitStatus())
			}
		}
		log.Println("Unimplemented for system where exec.ExitError.Sys() is not syscall.WaitStatus.")
		os.Exit(111)
	}
}
