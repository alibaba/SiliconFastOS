package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"oskit/docker"
	"oskit/util"
	"oskit/version"

	"gopkg.in/yaml.v2"
)

// GlobalConfig is the global tool configuration
type GlobalConfig struct {
	Pkg PkgConfig `yaml:"pkg"`
}

// PkgConfig is the config specific to the `pkg` subcommand
type PkgConfig struct {
}

var (
	// Config is the global tool configuration
	Config = GlobalConfig{}
)

func printVersion() {
	fmt.Printf("%s version %s\n", filepath.Base(os.Args[0]), version.Version)
	if version.GitCommit != "" {
		fmt.Printf("commit: %s\n", version.GitCommit)
	}
	os.Exit(0)
}

func readConfig() {
	cfgPath := filepath.Join(os.Getenv("HOME"), ".moby", "linuxkit", "config.yml")
	cfgBytes, err := os.ReadFile(cfgPath)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		fmt.Printf("Failed to read %q\n", cfgPath)
		os.Exit(1)
	}
	if err := yaml.Unmarshal(cfgBytes, &Config); err != nil {
		fmt.Printf("Failed to parse %q\n", cfgPath)
		os.Exit(1)
	}
}

func main() {
	flag.Usage = func() {
		fmt.Printf("USAGE: %s [options] COMMAND\n\n", filepath.Base(os.Args[0]))
		fmt.Printf("Commands:\n")
		fmt.Printf("  build       Build an image from a YAML file\n")
		fmt.Printf("  cache       Manage the local cache\n")
		fmt.Printf("  lookup      Lookup a docker image\n")
		fmt.Printf("  pkg         Package building\n")
		fmt.Printf("  run         Run a VM image on a local hypervisor\n")
		fmt.Printf("  serve       Run a local http server (for iPXE booting)\n")
		fmt.Printf("  version     Print version information\n")
		fmt.Printf("  help        Print this message\n")
		fmt.Printf("\n")
		fmt.Printf("Run '%s COMMAND --help' for more information on the command\n", filepath.Base(os.Args[0]))
		fmt.Printf("\n")
		fmt.Printf("Options:\n")
		flag.PrintDefaults()
	}

	readConfig()

	// Set up logging
	util.AddLoggingFlags(nil)
	flag.Parse()
	util.SetupLogging()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Printf("Please specify a command.\n\n")
		flag.Usage()
		os.Exit(1)
	}

	switch args[0] {
	case "run":
		run(args[1:])
	case "build":
		build(args[1:])
	case "cache":
		cache(args[1:])
	case "serve":
		serve(args[1:])
	case "lookup":
		lookup(args[1:])
	case "version":
		printVersion()
	case "pkg":
		pkgBuildPush(args[1:])
	case "image":
		docker.Version()
	case "help":
		flag.Usage()
	default:
		fmt.Printf("%q is not valid command.\n\n", args[0])
		flag.Usage()
		os.Exit(1)
	}
}
