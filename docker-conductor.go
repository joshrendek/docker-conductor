package main

import (
	"fmt"
	"github.com/joshrendek/docker-conductor/conductor"
	flag "github.com/ogier/pflag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ConductorDirections struct {
	Name      string
	Project   string
	Hosts     []string
	Container ConductorDirectionsContainer
}

type ConductorDirectionsContainer struct {
	Name        string
	Image       string
	Ports       map[string]string
	Environment []string
	Volumes     []string
	Dns         []string
}

func main() {

	var name *string = flag.StringP("name", "n", "", "Only run the instruction with this name")
	var project *string = flag.StringP("project", "p", "", "Only run the instruction that are apart of this project")
	flag.Parse()

	cd := []ConductorDirections{}

	data, _ := ioutil.ReadFile("conductor.yml")
	err := yaml.Unmarshal(data, &cd)
	if err != nil {
		panic(err)
	}

	for _, instr := range cd {
		if *name != "" {
			if instr.Name != *name {
				continue
			}
		}
		if *project != "" {
			if instr.Project != *project {
				continue
			}
		}
		// fmt.Printf("--- m:\n%v\n\n", instr)
		for _, host := range instr.Hosts {
			docker_ctrl := conductor.New(host)
			docker_ctrl.PullImage(instr.Container.Image + ":latest")
			container := docker_ctrl.FindContainer(instr.Container.Name)
			fmt.Println("Container ID: " + container.ID())
			if container.ID() != "" {
				if err := docker_ctrl.RemoveContainer(container.ID()); err != nil {
					panic(err)
				}
			}
			docker_ctrl.CreateAndStartContainer(conductor.ConductorContainerConfig{
				Name:        instr.Container.Name,
				Image:       instr.Container.Image,
				PortMap:     instr.Container.Ports,
				Environment: instr.Container.Environment,
				Volumes:     instr.Container.Volumes,
				Dns:         instr.Container.Dns,
			})
		}

	}

}
