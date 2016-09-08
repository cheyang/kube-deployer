package types

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/cheyang/fog/util/yaml"
	docker "github.com/docker/engine-api/types"
	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/state"
)

type Spec struct {
	Name            string                        `json:"ClusterName"` // The name of the cluster
	ClusterType     string                        `json:"ClusterType"` // The type of cluster
	VMSpecs         []VMSpec                      `json:"Vmspecs"`
	Properties      map[string]interface{}        `json:"Properties,omitempty"`
	DockerRun       *docker.ContainerCreateConfig `json:"DockerRun"`
	Run             []string                      `json:"Run"`
	CloudDriverName string                        `json:"Driver"`
	Update          bool                          `json:"-"` // Update an exist cluster, by default it's false.
}

type VMSpec struct {
	Name            string                 `json:"Name"`
	Roles           []string               `json:"Roles"`
	Properties      map[string]interface{} `json:"Properties,omitempty"`
	CloudDriverName string                 `json:"Driver"`
	Instances       int                    `json:"Instances,omitempty"`
	Start           int                    `json:"Start,omitempty"`
}

type Host struct {
	Err              error
	ErrMessage       string
	Name             string
	SSHUserName      string
	SSHPort          int
	SSHHostname      string
	PublicIPAddress  string
	PrivateIPAddress string // for most IAAS provider, it provides both public and private ip address
	SSHKeyPath       string
	Roles            []string
	State            state.State
	DriverName       string
	VMSpec
	Driver       drivers.Driver
	TemplateName string
}

func LoadSpec(configFile string) (spec Spec, err error) {
	spec = Spec{}
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return spec, err
	}
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return spec, err
	}
	decoder := yaml.NewYAMLToJSONDecoder(bytes.NewReader(data))
	err = decoder.Decode(&spec)
	return spec, err
}
