package cluster

import (
	"github.com/Sirupsen/logrus"
	"github.com/cheyang/fog/persist"
	"github.com/cheyang/fog/types"
	"github.com/cheyang/fog/util"
)

func Scaleout(s persist.Store, spec types.Spec, requiredRoleMap map[string]bool) error {
	spec.Update = true
	runningHosts, _, err := persist.LoadAllHosts(s)
	if err != nil {
		return err
	}

	// key is the vmspec name, value is the host name list
	runningHostMap, err := util.BuildRunningMap(runningHosts)
	if err != nil {
		return err
	}
	for i, vmSpec := range spec.VMSpecs {
		spec.VMSpecs[i].Start, err = nextNumber(runningHostMap, vmSpec.Name)
		if err != nil {
			return err
		}
	}

	newHosts, err := provisionVMs(spec, false)
	if err != nil {
		return err
	}

	// pick up the hosts for deployment
	hosts := make([]*types.Host, 0)
	for _, host := range runningHosts {
	role_loop:
		for _, role := range host.Roles {
			if _, found := requiredRoleMap[role]; found {
				hosts = append(hosts, host)
				break role_loop
			}
		}
	}
	hosts = append(hosts, newHosts...)

	err = configureIaaS(hosts, spec)
	if err != nil {
		return err
	}
	return runDeploy(hosts, spec)
}

// next number of the specified vmspec name
func nextNumber(runningHostMap map[string][]string, name string) (uint, error) {
	if orderedHostnames, found := runningHostMap[name]; found {
		maxIndex := len(orderedHostnames) - 1
		// s := strings.Split(orderedHostnames[maxIndex], "-")
		// max, err := strconv.Atoi(s[len(s)-1])
		hostname := orderedHostnames[maxIndex]
		_, max, err := util.ParseHostname(hostname)
		if err != nil {
			return 0, err
		}
		logrus.Infof("The max of %s is %d", name, max)
		return uint(max + 1), nil
	}
	return 0, nil
}
