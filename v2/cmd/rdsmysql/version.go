package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

// the version is set by goreleaser
var version = ""

func getVersion() string {
	if version != "" {
		return version
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	return info.Main.Version
}

func showVersion() {
	fmt.Printf("rdsmysql version %s built with %s %s/%s\n", getVersion(), runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
