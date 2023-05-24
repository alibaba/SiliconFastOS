package moby

import (
	"fmt"

	"github.com/containerd/containerd/reference"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	log "github.com/sirupsen/logrus"
	cache "oskit/cache"
	"oskit/docker"
	lktspec "oskit/spec"
)

func imagePullGet(ref *reference.Spec, trustedRef, architecture string, alwaysPull bool) (lktspec.ImageSource, error) {
	image := ref.String()
	pullImageName := image
	remoteOptions := []remote.Option{remote.WithAuthFromKeychain(authn.DefaultKeychain)}
	if trustedRef != "" {
		pullImageName = trustedRef
	}
	log.Debugf("ImagePull to cache %s trusted reference %s", image, pullImageName)

	log.Printf("Image %s not found in local cache, pulling", image)
	remoteRef, err := name.ParseReference(pullImageName)
	if err != nil {
		return cache.ImageSource{}, fmt.Errorf("invalid image name %s: %v", pullImageName, err)
	}

	_, err = remote.Get(remoteRef, remoteOptions...)
	if err != nil {
		return cache.ImageSource{}, fmt.Errorf("error getting manifest for trusted image %s: %v", pullImageName, err)
	}

	// ensure it includes our architecture
	return cache.NewSource(ref, architecture, nil), nil
}

// imagePull pull an image from the OCI registry to the cache.
// If the image root already is in the cache, use it, unless
// the option pull is set to true.
// if alwaysPull, then do not even bother reading locally
func imagePull(ref *reference.Spec, alwaysPull bool, cacheDir string, dockerCache bool, architecture string) (lktspec.ImageSource, error) {
	// several possibilities:
	// - alwaysPull: try to pull it down from the registry to linuxkit cache, then fail
	// - !alwaysPull && dockerCache: try to read it from docker, then try linuxkit cache, then try to pull from registry, then fail
	// - !alwaysPull && !dockerCache: try linuxkit cache, then try to pull from registry, then fail
	// first, try docker, if that is available
	if !alwaysPull && dockerCache {
		if err := docker.HasImage(ref); err == nil {
			return docker.NewSource(ref), nil
		}
		// docker is not required, so any error - image not available, no docker, whatever - just gets ignored
	}

	// next try the local cache
	if !alwaysPull {
		c, err := cache.NewProvider(cacheDir)
		if err != nil {
			return nil, err
		}
		if image, err := c.ValidateImage(ref, architecture); err == nil {
			return image, nil
		}
	}

	// if we made it here, we either did not have the image, or it was incomplete
	c, err := cache.NewProvider(cacheDir)
	if err != nil {
		return nil, err
	}
	return c.ImagePull(ref, ref.String(), architecture, alwaysPull)
}
