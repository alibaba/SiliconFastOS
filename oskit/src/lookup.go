package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	util "oskit/util"

	log "github.com/sirupsen/logrus"
)

func lookup(args []string) {
	runCmd := flag.NewFlagSet("lookup", flag.ExitOnError)
	runCmd.Usage = func() {
		fmt.Printf("USAGE: %s run [options] <img-name>\n\n", os.Args[0])
		fmt.Printf("Options:\n\n")
		runCmd.PrintDefaults()
	}

	local := runCmd.Bool("local", false, "Release the given version")
	_ = runCmd.Parse(args)

	for _, img := range runCmd.Args() {
		qemuImg := fmt.Sprintf("%s/silicon_fast/dive:latest", util.DOCKER_REPO)

		var c []string
		if *local {
			c = []string{}
		} else {
			c = []string{"docker", "run", "--network", "host", "-v", "/var/run/docker.sock:/var/run/docker.sock", "-it", qemuImg}
		}
		c = append(c, "dive", img)

		exec := exec.CommandContext(context.Background(), "/bin/sh", "-c", strings.Join(c, " "))

		exec.Stdout = os.Stdout
		exec.Stderr = os.Stderr
		exec.Stdin = os.Stdin
		exec.Env = os.Environ()

		log.Debugf("Executing: %v", exec.Args)

		err := exec.Run()
		if err != nil {
			if isExecErrNotFound(err) {
				fmt.Errorf("linuxkit pkg requires docker to be installed")
			}
		}

	}
}
