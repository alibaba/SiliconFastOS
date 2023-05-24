package pkglib

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-version"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"oskit/docker"
	"oskit/moby"

	"oskit/util"
)

const (
	buildersEnvVar = "LINUXKIT_BUILDERS"
	envVarCacheDir = "LINUXKIT_CACHE"
	// this is the most recent manifest pointed to by moby/buildkit:master as of 2022-07-22, so it includes
	// our required commit. Once there is a normal semver tag later than this, we should switch to it.
	defaultBuilderImage = util.DOCKER_REPO + "/silicon_fast/buildkit:1.0.0"
)

// Contains fields settable in the build.yml
type pkgInfo struct {
	Image        string            `yaml:"image"`
	Org          string            `yaml:"org"`
	Arches       []string          `yaml:"arches"`
	ExtraSources []string          `yaml:"extra-sources"`
	GitRepo      string            `yaml:"gitrepo"` // ??
	Network      bool              `yaml:"network"`
	DisableCache bool              `yaml:"disable-cache"`
	Config       *moby.ImageConfig `yaml:"config"`
	BuildArgs    *[]string         `yaml:"buildArgs,omitempty"`
	Release      string            `yaml:"release"`
	Depends      struct {
		DockerImages struct {
			TargetDir string   `yaml:"target-dir"`
			Target    string   `yaml:"target"`
			FromFile  string   `yaml:"from-file"`
			List      []string `yaml:"list"`
		} `yaml:"docker-images"`
	} `yaml:"depends"`
}

// Pkg encapsulates information about a package's source
type Pkg struct {
	// These correspond to pkgInfo fields
	image     string
	org       string
	arches    []string
	network   bool
	trust     bool
	cache     bool
	config    *moby.ImageConfig
	buildArgs *[]string

	push bool

	force          bool
	path           string
	hash           string
	builderIns     string
	builderImage   string
	builderRestart bool
	release        string
	nobuild        bool
	md5            string
	target         string
}

// NewFromCLI creates a range of Pkg from a set of CLI arguments. Calls fs.Parse()
func NewFromCLI(fs *flag.FlagSet, args ...string) ([]Pkg, error) {
	// Defaults
	piBase := pkgInfo{
		Org:          util.DOCKER_REPO,
		Arches:       []string{"arm64"},
		Network:      true,
		DisableCache: false,
	}

	// TODO(ijc) look for "$(git rev-parse --show-toplevel)/.build-defaults.yml"?

	// Ideally want to look at every directory from root to `pkg`
	// for this file but might be tricky to arrange ordering-wise.
	force := fs.Bool("force", false, "Force rebuild even if image is in local cache")
	//ignoreCache := fs.Bool("no-cache", false, "Ignore cached intermediate images, always pulling from registry")
	builders := fs.String("builders", "", "Which builders instance to use for which platforms")
	builderImage := fs.String("builder-image", defaultBuilderImage, "buildkit builder container image to use")
	builderArgs := fs.String("builder-arg", "", "buildkit builder container image to use")
	target := fs.String("target", "", "buildkit builder container image to use")
	//buildArg := flags.String("build-arg", "", "buildkit builder container image to use")
	builderRestart := fs.Bool("builder-restart", false, "force restarting builder, even if container with correct name and image exists")
	all := fs.Bool("all", false, "build all build yml")
	push := fs.Bool("push", false, "build all build yml")

	var (
		release    *string
		nobuild    *bool
		nobuildRef = false
	)

	nobuild = &nobuildRef
	release = fs.String("release", "", "Release the given version")
	nobuild = fs.Bool("nobuild", false, "Skip building the image before pushing, conflicts with -force")

	// These override fields in pi below, bools are in both forms to allow user overrides in either direction.
	// These will apply to all packages built.
	argDisableCache := fs.Bool("disable-cache", piBase.DisableCache, "Disable build cache")
	argEnableCache := fs.Bool("enable-cache", !piBase.DisableCache, "Enable build cache")
	argNoNetwork := fs.Bool("nonetwork", !piBase.Network, "Disallow network use during build")
	argNetwork := fs.Bool("network", piBase.Network, "Allow network use during build")

	argOrg := fs.String("org", piBase.Org, "Override the hub org")

	// Other arguments
	var buildYML, hash, hashPath string
	var devMode bool

	fs.StringVar(&buildYML, "build-yml", "build.yml", "Override the name of the yml file")
	fs.StringVar(&hash, "hash", "", "Override the image hash (default is to query git for the package's tree-sh)")
	fs.StringVar(&hashPath, "hash-path", "", "Override the directory to use for the image hash, must be a parent of the package dir (default is to use the package dir)")
	fs.BoolVar(&devMode, "dev", false, "Force org and hash to $USER and \"dev\" respectively")

	util.AddLoggingFlags(fs)

	_ = fs.Parse(args)

	util.SetupLogging()

	if fs.NArg() < 1 {
		return nil, fmt.Errorf("at least one pkg directory is required")
	}

	if *nobuild && *force {
		fmt.Fprint(os.Stderr, "flags -force and -nobuild conflict")
		os.Exit(1)
	}

	if *all && buildYML != "build.yml" {
		return nil, errors.New("flag -all and -build-yml conflict")
	}

	var pkgs []Pkg
	for _, pkg := range fs.Args() {
		var (
			pkgHashPath string
			pkgHash     = hash
		)
		pkgPath, err := filepath.Abs(pkg)
		if err != nil {
			return nil, err
		}

		if pkgHash == "" {
			pkgHash, err = util.DirMd5(pkgPath)
			if err != nil {
				return nil, fmt.Errorf("pkg hash failed")
			}
		}
		if hashPath == "" {
			pkgHashPath = pkgPath
		} else {
			pkgHashPath, err = filepath.Abs(hashPath)
			if err != nil {
				return nil, err
			}

			if !strings.HasPrefix(pkgPath, pkgHashPath) {
				return nil, fmt.Errorf("Hash path is not a prefix of the package path")
			}

			// TODO(ijc) pkgPath and hashPath really ought to be in the same git tree too...
		}

		if devMode {
			var newA []string
			for _, cliArg := range strings.Split(*builderArgs, ",") {
				if cliArg == "" {
					continue
				}
				newA = append(newA, cliArg)
			}

			// If --org is also used then this will be overwritten
			// by argOrg when we iterate over the provided options
			// in the fs.Visit block below.
			p := Pkg{
				image:          filepath.Base(pkgPath),
				org:            os.Getenv("USER"),
				hash:           pkgHash,
				network:        true,
				cache:          *argDisableCache,
				buildArgs:      &newA,
				path:           pkgPath,
				force:          *force,
				builderIns:     *builders,
				builderImage:   *builderImage,
				builderRestart: *builderRestart,
				release:        *release,
				nobuild:        *nobuild,
				target:         *target,
				push:           *push,
			}
			if err = p.build(); err != nil {
				fmt.Errorf(err.Error())
			}
			continue
		}

		var buildYmls []string
		if *all {
			buildYmls, err = filepath.Glob(filepath.Join(pkgPath, "build*.yml"))
			if err != nil {
				log.Fatal(err)
			}
		} else {
			buildYmls = append(buildYmls, filepath.Join(pkgPath, buildYML))
		}

		for _, file := range buildYmls {

			// make our own copy of piBase. We could use some deepcopy library, but it is just as easy to marshal/unmarshal
			pib, err := yaml.Marshal(&piBase)
			if err != nil {
				return nil, err
			}
			var pi pkgInfo
			if err := yaml.Unmarshal(pib, &pi); err != nil {
				return nil, err
			}

			b, err := os.ReadFile(file)
			if err != nil {
				return nil, err
			}

			if err := yaml.Unmarshal(b, &pi); err != nil {
				return nil, err
			}

			if pi.Image == "" {
				return nil, fmt.Errorf("image field is required")
			}

			// Go's flag package provides no way to see if a flag was set
			// apart from Visit which iterates over only those which were
			// set. This must be run here, rather than earlier, because we need to
			// have read it from the build.yml file first, then override based on CLI.
			fs.Visit(func(f *flag.Flag) {
				switch f.Name {
				case "disable-cache":
					pi.DisableCache = *argDisableCache
				case "enable-cache":
					pi.DisableCache = !*argEnableCache
				case "network":
					pi.Network = *argNetwork
				case "nonetwork":
					pi.Network = !*argNoNetwork
				case "org":
					pi.Org = *argOrg
				}
			})

			var newA []string
			for _, cliArg := range strings.Split(*builderArgs, ",") {
				if cliArg == "" {
					continue
				}
				exist := false
				if pi.BuildArgs != nil {
					argArr := strings.SplitN(cliArg, "=", 2)
					k := argArr[0]
					v := argArr[1]
					for i, s := range *pi.BuildArgs {
						a := strings.SplitN(s, "=", 2)
						k1 := a[0]
						v1 := a[1]
						if k1 == k {
							exist = true
							if v != v1 {
								(*pi.BuildArgs)[i] = cliArg
							}
						}
					}
				}
				if !exist {
					newA = append(newA, cliArg)
				}
			}

			if pi.BuildArgs == nil {
				pi.BuildArgs = &newA
			} else {
				*pi.BuildArgs = append(*pi.BuildArgs, newA...)
			}

			var r string
			if *all {
				r = pi.Release
			} else {
				r = *release
			}

			p := Pkg{
				image:          pi.Image,
				org:            pi.Org,
				hash:           pkgHash,
				arches:         pi.Arches,
				network:        pi.Network,
				cache:          !pi.DisableCache,
				config:         pi.Config,
				buildArgs:      pi.BuildArgs,
				path:           pkgPath,
				force:          *force,
				builderIns:     *builders,
				builderImage:   *builderImage,
				builderRestart: *builderRestart,
				release:        r,
				nobuild:        *nobuild,
				target:         *target,
				push:           *push,
			}
			if err = p.build(); err != nil {
				fmt.Errorf(err.Error())
			}
		}
	}
	return pkgs, nil
}

func (p Pkg) build() error {
	if !p.force {
		if b := docker.TagExist(p.image, p.hash); b {
			fmt.Printf("Image exist in docker")
			return nil
		}
	}

	localVersion, err := docker.Version()
	if err != nil {
		return err
	}
	lv, _ := version.NewVersion(localVersion)
	buildxVersion, _ := version.NewVersion("19.03")

	var cmds []exec.Cmd
	if lv.Compare(buildxVersion) <= 0 {
		if p.nobuild {
			fmt.Printf("Pushing %q without building", p.Tag())
		} else {
			fmt.Printf("Pushing %q with building", p.Tag())
			args := append([]string{"docker", "build", "--force-rm"}, p.netFlag()...)
			args = append(args, p.buildArg()...)
			args = append(args, p.releaseFlag()...)
			args = append(args, ".")
			cmds = append(cmds, *exec.CommandContext(context.Background(), "/bin/sh", "-c", strings.Join(args, " ")))
		}
		if p.push {
			if p.release != "" {
				releaseRef, _ := p.ReleaseTag(p.release)
				cmds = append(cmds, *exec.CommandContext(context.Background(), "/bin/sh", "-c", "docker push "+releaseRef))
			} else {
				cmds = append(cmds, *exec.CommandContext(context.Background(), "/bin/sh", "-c", "docker push "+p.FullTag()))
			}
		}
	} else {
		args := append([]string{"docker", "buildx", "build"}, p.builder()...)
		args = append(args, p.netFlag()...)
		args = append(args, p.buildArg()...)
		args = append(args, p.cacheFlag())
		args = append(args, p.pushFlag())
		args = append(args, p.targetFlag()...)
		args = append(args, p.releaseFlag()...)
		args = append(args, ".")
		cmds = append(cmds, *exec.CommandContext(context.Background(), "/bin/sh", "-c", strings.Join(args, " ")))

	}

	for _, exeCmd := range cmds {
		exeCmd.Dir = p.path
		exeCmd.Stdout = os.Stdout
		exeCmd.Stderr = os.Stderr
		exeCmd.Stdin = os.Stdin
		exeCmd.Env = os.Environ()

		log.Debugf("Executing: %v", exeCmd.Args)

		err := exeCmd.Run()
		if err != nil {
			if isExecErrNotFound(err) {
				return fmt.Errorf("linuxkit pkg requires docker to be installed")
			}
			return err
		}
	}
	return nil
}

// Hash returns the hash of the package
func (p Pkg) Hash() string {
	return p.hash
}

// ReleaseTag returns the tag to use for a particular release of the package
func (p Pkg) ReleaseTag(release string) (string, error) {
	if release == "" {
		return "", fmt.Errorf("a release tag is required")
	}
	tag := p.org + "/" + p.image + ":" + release
	return tag, nil
}

// Tag returns the tag to use for the package
func (p Pkg) Tag() string {
	t := p.hash
	if t == "" {
		t = "latest"
	}
	return p.org + "/" + p.image + ":" + t
}

// FullTag returns a reference expanded tag
func (p Pkg) FullTag() string {
	return util.ReferenceExpand(p.Tag())
}

// TrustEnabled returns true if trust is enabled
func (p Pkg) TrustEnabled() bool {
	return p.trust
}

// Arches which arches this can be built for
func (p Pkg) Arches() []string {
	return p.arches
}

//nolint:unused // will be used when linuxkit cache is eliminated and we return to docker image cache
func (p Pkg) archSupported(want string) bool {
	for _, supp := range p.arches {
		if supp == want {
			return true
		}
	}
	return false
}

func (p Pkg) builder() []string {
	return []string{"--builder", "buildkit"}
}

func (p Pkg) buildArg() []string {
	var buildArg []string
	for _, arg := range *p.buildArgs {
		buildArg = append(buildArg, "--build-arg", arg)
	}
	return buildArg
}

func (p Pkg) pushFlag() string {
	if p.push {
		return "--push"
	}
	return "--load"
}

func (p Pkg) netFlag() []string {
	if p.network {
		return []string{"--network", "host"}
	}
	return []string{}
}

func (p Pkg) targetFlag() []string {
	if p.target != "" {
		return []string{"--target ", p.target}
	}
	return []string{}
}

func (p Pkg) cacheFlag() string {
	if p.cache {
		return ""
	}
	return " --no-cache"
}

func (p Pkg) releaseFlag() []string {
	s := []string{"-t", p.FullTag()}
	if p.release != "" {
		r, _ := p.ReleaseTag(p.release)
		s = append(s, "-t", r)
	}
	return s
}

// Expands path from relative to abs against base, ensuring the result is within base, but is not base itself. Field is the fieldname, to be used for constructing the error.
func makeAbsSubpath(field, base, path string) (string, error) {
	if path == "" {
		return "", nil
	}

	if filepath.IsAbs(path) {
		return "", fmt.Errorf("%s must be relative to package directory", field)
	}

	p, err := filepath.Abs(filepath.Join(base, path))
	if err != nil {
		return "", err
	}

	if p == base {
		return "", fmt.Errorf("%s must not be exactly the package directory", field)
	}

	if !strings.HasPrefix(p, base) {
		return "", fmt.Errorf("%s must be within package directory", field)
	}

	return p, nil
}
