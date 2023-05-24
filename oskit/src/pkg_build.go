package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	pkglib "oskit/pkg"
)

func pkgBuildPush(args []string) {
	flags := flag.NewFlagSet("pkg", flag.ExitOnError)
	flags.Usage = func() {
		invoked := filepath.Base(os.Args[0])
		name := "push"
		fmt.Fprintf(os.Stderr, "USAGE: %s pkg %s [options] path\n\n", name, invoked)
		fmt.Fprintf(os.Stderr, "'path' specifies the path to the package source directory.\n")
		fmt.Fprintf(os.Stderr, "\n")
		flags.PrintDefaults()
	}

	// some logic clarification:
	// pkg                    - builds unless is in cache or published in registry
	// pkg --force           - always builds even if is in cache or published in registry
	// pkg --force --pull    - always builds even if is in cache or published in registry; --pull ignored
	// pkg                    - always builds unless is in cache
	// pkg --force            - always builds even if is in cache
	// pkg --nobuild          - skips build; if not in cache, fails
	// pkg --nobuild --force  - nonsensical

	_, err := pkglib.NewFromCLI(flags, args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
