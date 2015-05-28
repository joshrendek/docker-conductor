package conductor

import (
	"github.com/fsouza/go-dockerclient"
	"os"
	"strings"
)

type Conductor struct {
	Client *docker.Client
}

type ConductorContainer struct {
	Conductor *Conductor
	Container docker.APIContainers
}

type ConductorContainerConfig struct {
	Name        string
	Image       string
	PortMap     map[string]string
	Environment []string
}

func (c *ConductorContainer) ID() string {
	return c.Container.ID
}

func New(Host string) *Conductor {
	client, _ := docker.NewClient(Host)
	return &Conductor{Client: client}
}

func (c *Conductor) PullImage(image string) {
	parsed := strings.Split(image, "/")
	registry := parsed[0]
	image_and_tag := strings.Join(parsed[1:], "")
	parsed_image := strings.Split(image_and_tag, ":")
	repository := parsed_image[0]
	tag := parsed_image[1]
	opts := docker.PullImageOptions{
		Repository:   repository,
		Registry:     registry,
		Tag:          tag,
		OutputStream: os.Stdout,
	}

	c.Client.PullImage(opts, docker.AuthConfiguration{})
}

func (c *Conductor) CreateAndStartContainer(cfg ConductorContainerConfig) {

	portBindings := map[docker.Port][]docker.PortBinding{}

	for k, v := range cfg.PortMap {
		portBindings[docker.Port(k)] = []docker.PortBinding{{HostIP: "0.0.0.0", HostPort: v}}
	}

	hostConfig := &docker.HostConfig{PortBindings: portBindings}

	container, err := c.Client.CreateContainer(docker.CreateContainerOptions{
		Name:       cfg.Name,
		Config:     &docker.Config{Image: cfg.Image, Env: cfg.Environment},
		HostConfig: hostConfig,
	})

	if err != nil {
		panic(err)
	}

	c.Client.StartContainer(container.ID, hostConfig)
}

func (c *Conductor) RemoveContainer(id string) error {
	return c.Client.RemoveContainer(docker.RemoveContainerOptions{ID: id, Force: true})
}

func (c *Conductor) FindContainer(needle string) *ConductorContainer {
	containers, _ := c.Client.ListContainers(docker.ListContainersOptions{All: true})
	for _, container := range containers {
		for _, name := range container.Names {
			if name == "/"+needle {
				return &ConductorContainer{Conductor: c, Container: container}
			}
		}
	}
	return &ConductorContainer{}
}
