package main

import (
	"flag"
	"fmt"

	cachepkg "oskit/cache"

	log "github.com/sirupsen/logrus"
)

func cacheList(args []string) {
	flags := flag.NewFlagSet("list", flag.ExitOnError)

	cacheDir := flagOverEnvVarOverDefaultString{def: defaultLinuxkitCache(), envVar: envVarCacheDir}
	flags.Var(&cacheDir, "cache", fmt.Sprintf("Directory for caching and finding cached image, overrides env var %s", envVarCacheDir))

	if err := flags.Parse(args); err != nil {
		log.Fatal("Unable to parse args")
	}

	// list all of the images and content in the cache
	p, err := cachepkg.Get(cacheDir.String())
	if err != nil {
		log.Fatalf("unable to read a local cache: %v", err)
	}
	images, err := cachepkg.ListImages(p)
	if err != nil {
		log.Fatalf("error reading image names: %v", err)
	}
	log.Printf("%-80s %s", "image name", "root manifest hash")
	for name, hash := range images {
		log.Printf("%-80s %s", name, hash)
	}
}
