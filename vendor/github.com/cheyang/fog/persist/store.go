package persist

import (
	"github.com/Sirupsen/logrus"
	"github.com/cheyang/fog/types"
)

type Store interface {
	// Exists returns whether a machine exists or not
	Exists(name string) (bool, error)

	// List returns a list of all hosts in the store
	List() ([]string, error)

	// Load loads a host by name
	Load(name string) (*types.Host, error)

	// Remove removes a machine from the store
	Remove(name string) error

	// Save persists a machine in the store
	Save(host *types.Host) error

	GetRoot() string

	GetMachinesDir() string

	// create the store path
	CreateStorePath(name string) error

	CreateDeploymentDir() error

	GetDeploymentDir() string

	SaveSpec(specs *types.Spec) error

	LoadSpec() (*types.Spec, error)
}

func LoadHosts(s Store, hostNames []string) ([]*types.Host, map[string]error) {
	loadedHosts := []*types.Host{}
	errors := map[string]error{}

	for _, hostName := range hostNames {
		h, err := s.Load(hostName)
		if err != nil {
			errors[hostName] = err
		} else {
			loadedHosts = append(loadedHosts, h)
		}
	}

	return loadedHosts, errors
}

func LoadAllHosts(s Store) ([]*types.Host, map[string]error, error) {
	hostNames, err := s.List()
	logrus.Infof("hostname :%v", hostNames)
	if err != nil {
		return nil, nil, err
	}
	loadedHosts, hostInError := LoadHosts(s, hostNames)
	return loadedHosts, hostInError, nil
}
