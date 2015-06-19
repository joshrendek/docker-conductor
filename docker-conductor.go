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
	var force_deploy *bool = flag.BoolP("force", "f", false, "Force a redeploy of everything in the conductor.yml file")
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
			host_log := log.New("\t\t\t\t[host]", host, "[container]", instr.Container.Name)

			host_log.Info("[ ] pulling image")
			pulled_image, err := docker_ctrl.PullImage(instr.Container.Image + ":latest")
			if err != nil {
				host_log.Error("Error pulling image: " + err.Error())
			}
			host_log.Info("[x] finished pulling image")

			if container.Container != nil {
				if pulled_image == container.Container.Image && *force_deploy == false {
					host_log.Info("skipping, container running latest image : " + pulled_image)
					continue
				}

				if container.ID() != "" {
					if err := docker_ctrl.RemoveContainer(container.ID()); err != nil {
						host_log.Error(err.Error())
					}
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
