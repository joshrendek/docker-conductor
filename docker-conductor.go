package main

import (
	"github.com/joshrendek/docker-conductor/conductor"
	flag "github.com/ogier/pflag"
	log "gopkg.in/inconshreveable/log15.v2"
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
			container := docker_ctrl.FindContainer(instr.Container.Name)
			host_log := log.New("\t\t\t\t[host]", host, "[container]", container.Container.Name)

			host_log.Info("[ ] pulling image")
			pulled_image := docker_ctrl.PullImage(instr.Container.Image + ":latest")
			host_log.Info("[x] finished pulling image")
			//log.Info("Container ID: " + container.ID())
			//log.Info("Container image: " + container.Container.Image)

			if pulled_image == container.Container.Image {
				host_log.Info("skipping, container running latest image")
				continue
			}

			if container.ID() != "" {
				if err := docker_ctrl.RemoveContainer(container.ID()); err != nil {
					host_log.Error(err.Error())
				}
			}
			host_log.Info("[ ] creating container")
			docker_ctrl.CreateAndStartContainer(conductor.ConductorContainerConfig{
				Name:        instr.Container.Name,
				Image:       instr.Container.Image,
				PortMap:     instr.Container.Ports,
				Environment: instr.Container.Environment,
				Volumes:     instr.Container.Volumes,
				Dns:         instr.Container.Dns,
			})
			host_log.Info("[x] finished creating container")
		}

	}

}
