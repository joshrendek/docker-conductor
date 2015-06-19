package conductor

import (
	"github.com/fsouza/go-dockerclient"
	//	"strings"
)

type Conductor struct {
	Client *docker.Client
}

type ConductorContainer struct {
	Conductor *Conductor
	Container *docker.Container
}

type ConductorContainerConfig struct {
	Name        string
	Image       string
	PortMap     map[string]string
	Environment []string
	Volumes     []string
	Dns         []string
}

func (c *ConductorContainer) ID() string {
	return c.Container.ID
}

func New(Host string) *Conductor {
	client, _ := docker.NewClient(Host)
	return &Conductor{Client: client}
}

func (c *Conductor) PullImage(image string) (string, error) {
	opts := docker.PullImageOptions{
		Repository:    image,
		RawJSONStream: true,
	}
	err := c.Client.PullImage(opts, docker.AuthConfiguration{})
	latest_image, _ := c.Client.InspectImage(image)
	return latest_image.ID, err
}

func (c *Conductor) CreateAndStartContainer(cfg ConductorContainerConfig) {
	portBindings := map[docker.Port][]docker.PortBinding{}

	for k, v := range cfg.PortMap {
		portBindings[docker.Port(k)] = []docker.PortBinding{{HostIP: "0.0.0.0", HostPort: v}}
	}

	hostConfig := &docker.HostConfig{PortBindings: portBindings,
		Binds: cfg.Volumes, DNS: cfg.Dns,
		RestartPolicy: docker.AlwaysRestart()}

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
				real_container, _ := c.Client.InspectContainer(container.ID)
				return &ConductorContainer{Conductor: c, Container: real_container}
			}
		}
	}
	return &ConductorContainer{}
}
