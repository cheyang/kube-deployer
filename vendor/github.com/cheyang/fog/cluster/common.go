package cluster

import (
	provider_registry "github.com/cheyang/fog/cloudprovider/registry"
	"github.com/cheyang/fog/cluster/ansible"
	"github.com/cheyang/fog/cluster/deploy"
	host_utils "github.com/cheyang/fog/host"
	"github.com/cheyang/fog/types"
	"github.com/cheyang/fog/util"
)

func provisionVMs(spec types.Spec, save bool) (hosts []*types.Host, err error) {
	bus := make(chan *types.Host)
	defer close(bus)
	vmSpecs, err := host_utils.BuildHostConfigs(spec, save)
	if err != nil {
		return hosts, err
	}

	hostCount := len(vmSpecs)
	err = host_utils.CreateInBatch(vmSpecs, bus)
	if err != nil {
		return hosts, err
	}

	hosts = make([]*types.Host, hostCount)
	for i := 0; i < hostCount; i++ {
		hosts[i] = <-bus
	}

	for _, host := range hosts {
		if host.Err != nil {
			return hosts, host.Err
		}
	}

	return hosts, nil
}

func configureIaaS(hosts []*types.Host, spec types.Spec) (err error) {
	storage, err := util.GetStorage(spec.Name)
	if err != nil {
		return err
	}
	cp := provider_registry.GetProvider(spec.CloudDriverName, spec.ClusterType, storage)
	if cp != nil {
		cp.SetHosts(hosts)
		err = cp.Configure() // configure IaaS
		if err != nil {
			return err
		}
	}
	return nil
}

func runDeploy(hosts []*types.Host, spec types.Spec) (err error) {
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
