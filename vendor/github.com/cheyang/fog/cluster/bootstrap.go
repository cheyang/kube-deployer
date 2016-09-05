package cluster

import (
	"github.com/Sirupsen/logrus"
	"github.com/cheyang/fog/cluster/ansible"
	"github.com/cheyang/fog/cluster/deploy"
	"github.com/cheyang/fog/host"
	"github.com/cheyang/fog/types"
	"github.com/cheyang/fog/util/dump"
)

func Bootstrap(spec types.Spec) error {

	err := types.Validate(spec)
	if err != nil {
		return err
	}

	logrus.Infof("spec: %+v", spec)

	//register core dump tool
	dump.InstallCoreDumpGenerator()

	bus := make(chan types.Host)
	defer close(bus)
	vmSpecs, err := host.BuildHostConfigs(spec)
	if err != nil {
		return err
	}

	hostCount := len(vmSpecs)
	err = host.CreateInBatch(vmSpecs, bus)
	if err != nil {
		return err
	}

	hosts := make([]types.Host, hostCount)
	for i := 0; i < hostCount; i++ {
		hosts[i] = <-bus
	}

	for _, host := range hosts {
		if host.Err != nil {
			return host.Err
		}
	}

	cp := initProivder(spec.CloudDriverName, spec.ClusterType)
	if cp != nil {
		cp.SetHosts(hosts)
		cp.Configure() // configure IaaS
	}

	var deployer deploy.Deployer
	deployer, err = ansible.NewDeployer(spec.Name)
	if err != nil {
		return err
	}
	deployer.SetHosts(hosts)
	if len(spec.Run) > 0 {
		deployer.SetCommander(spec.Run)
	} else {
		deployer.SetCommander(spec.DockerRun)
	}

	return deployer.Run()
}
