package main

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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
	Name    string
	Image   string
	PortMap map[string]string
}

func (c *ConductorContainer) ID() string {
	return c.Container.ID
}

func NewConductor(Host string) *Conductor {
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
		Config:     &docker.Config{Image: cfg.Image},
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
	containers, _ := c.Client.ListContainers(docker.ListContainersOptions{All: false})
	for _, container := range containers {
		for _, name := range container.Names {
			if name == "/"+needle {
				return &ConductorContainer{Conductor: c, Container: container}
			}
		}
	}
	return &ConductorContainer{}
}

type ConductorDirections struct {
	Name      string
	Hosts     []string
	Container ConductorDirectionsContainer
}

type ConductorDirectionsContainer struct {
	Name  string
	Image string
	Ports map[string]string
}

func main() {
	cd := []ConductorDirections{}

	data, _ := ioutil.ReadFile("conductor.yml")
	err := yaml.Unmarshal(data, &cd)
	if err != nil {
		panic(err)
	}

	for _, instr := range cd {
		// fmt.Printf("--- m:\n%v\n\n", instr)

		for _, host := range instr.Hosts {
			conductor := NewConductor(host)
			conductor.PullImage(instr.Container.Image + ":latest")
			container := conductor.FindContainer(instr.Container.Name)
			fmt.Println("Container ID: " + container.ID())
			conductor.RemoveContainer(container.ID())
			conductor.CreateAndStartContainer(ConductorContainerConfig{
				Name:    instr.Container.Name,
				Image:   instr.Container.Image,
				PortMap: instr.Container.Ports,
			})
		}

	}

}
