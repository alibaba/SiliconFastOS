package docker

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/containerd/containerd/reference"
	"github.com/docker/cli/cli/connhelper"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

var (
	clientOnce     sync.Once
	memoizedClient *client.Client
	errClient      error
)

// Client get a docker client.
func Client() (*client.Client, error) {
	clientOnce.Do(func() {
		memoizedClient, errClient = createClient()
	})
	return memoizedClient, errClient
}

func createClient() (*client.Client, error) {
	options := []client.Opt{
		client.WithAPIVersionNegotiation(),
		client.WithTLSClientConfigFromEnv(),
		client.WithHostFromEnv(),
	}

	// Support connection over ssh.
	if host := os.Getenv(client.EnvOverrideHost); host != "" {
		helper, err := connhelper.GetConnectionHelper(host)
		if err != nil {
			return nil, err
		}
		if helper != nil {
			options = append(options, client.WithDialContext(helper.Dialer))
		}
	}

	return client.NewClientWithOpts(options...)
}

// HasImage check if the provided ref is available in the docker cache.
func HasImage(ref *reference.Spec) error {
	log.Debugf("docker inspect image: %s", ref)
	cli, err := Client()
	if err != nil {
		return err
	}
	_, err = InspectImage(cli, ref)

	return err
}

// InspectImage inspect the provided ref.
func InspectImage(cli *client.Client, ref *reference.Spec) (dockertypes.ImageInspect, error) {
	log.Debugf("docker inspect image: %s", ref)

	inspect, _, err := cli.ImageInspectWithRaw(context.Background(), ref.String())
	if err != nil {
		return dockertypes.ImageInspect{}, err
	}

	log.Debugf("docker inspect image: %s...Done", ref)

	return inspect, nil
}

// Create create a container from the given image in docker, returning the full hash ID
// of the created container. Does not start the container.
func Create(image string, withNetwork bool) (string, error) {
	log.Debugf("docker create: %s", image)
	cli, err := Client()
	if err != nil {
		return "", errors.New("could not initialize Docker API client")
	}
	// we do not ever run the container, so /dev/null is used as command
	config := &container.Config{
		Cmd:             []string{"/dev/null"},
		Image:           image,
		NetworkDisabled: !withNetwork,
	}

	respBody, err := cli.ContainerCreate(context.Background(), config, nil, nil, nil, "")
	if err != nil {
		return "", err
	}

	log.Debugf("docker create: %s...Done", image)
	return respBody.ID, nil
}

// Export export the provided container ID from docker using `docker export`.
// The container must already exist.
func Export(container string) (io.ReadCloser, error) {
	log.Debugf("docker export: %s", container)
	cli, err := Client()
	if err != nil {
		return nil, errors.New("could not initialize Docker API client")
	}
	return cli.ContainerExport(context.Background(), container)
}

// Save save the provided image ref.
func Save(image string) (io.ReadCloser, error) {
	log.Debugf("docker save: %s", image)
	cli, err := Client()
	if err != nil {
		return nil, errors.New("could not initialize Docker API client")
	}
	return cli.ImageSave(context.Background(), []string{image})
}

// Rm remove the given container from docker.
func Rm(container string) error {
	log.Debugf("docker rm: %s", container)
	cli, err := Client()
	if err != nil {
		return errors.New("could not initialize Docker API client")
	}
	if err = cli.ContainerRemove(context.Background(), container, dockertypes.ContainerRemoveOptions{}); err != nil {
		return err
	}
	log.Debugf("docker rm: %s...Done", container)
	return nil
}

func Rmi(image string) error {
	log.Debugf("docker rm: %s", image)
	cli, err := Client()
	if err != nil {
		return errors.New("could not initialize Docker API client")
	}
	if rets, err := cli.ImageRemove(context.Background(), image, dockertypes.ImageRemoveOptions{}); err != nil {
		return err
	} else {
		for _, ret := range rets {
			log.Debugf("%s \n %s", ret.Deleted, ret.Untagged)
		}
	}
	log.Debugf("docker rm: %s...Done", image)
	return nil
}

func Ls() error {
	log.Debugf("docker ls:")
	cli, err := Client()
	if err != nil {
		return errors.New("could not initialize Docker API client")
	}
	var imageList []dockertypes.ImageSummary
	if imageList, err = cli.ImageList(context.Background(), dockertypes.ImageListOptions{}); err != nil {
		log.Printf(err.Error())
		return err
	}

	for _, image := range imageList {
		if len(image.RepoTags) == 0 {
			Rmi(image.ID)
		}
		for _, tag := range image.RepoTags {
			if strings.Split(tag, ":")[0] == "<none>" {
				Rmi(image.ID)
			}
		}
	}
	log.Debugf("docker ls: %s...Done")
	return nil
}

func TagExist(ref string, tag string) bool {
	cli, err := Client()
	if err != nil {
		log.Errorf("could not initialize Docker API client")
		return false
	}
	var imageList []dockertypes.ImageSummary
	if imageList, err = cli.ImageList(context.Background(), dockertypes.ImageListOptions{}); err != nil {
		log.Errorf(err.Error())
		return false
	}

	for _, image := range imageList {
		for _, tag := range image.RepoTags {
			img := strings.Split(tag, ":")[0]
			t := strings.Split(tag, ":")[1]
			if img == ref && t == tag {
				return true
			}
		}
	}
	log.Debugf("Do not find same image named %s:%s", ref, tag)
	return false
}

func Version() (string, error) {
	cli, err := Client()
	if err != nil {
		return "", errors.New("could not initialize Docker API client")
	}
	version, _ := cli.ServerVersion(context.Background())
	log.Debugf("Docker version: %s", version.Version)
	return version.Version, nil
}
