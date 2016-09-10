package cluster

import (
	"github.com/Sirupsen/logrus"
	"github.com/cheyang/fog/types"
	"github.com/cheyang/fog/util/dump"
)

func Bootstrap(spec types.Spec) (hosts []types.Host, err error) {

	err = types.Validate(spec)
	if err != nil {
		return hosts, err
	}

	logrus.Infof("spec: %+v", spec)
	//register core dump tool
	dump.InstallCoreDumpGenerator()

	// save spec
	hosts, err = provisionVMs(spec, true)
	if err != nil {
		return hosts, err
	}

	err = configureIaaS(hosts, spec)
	if err != nil {
		return hosts, err
	}

	return hosts, runDeploy(hosts, spec)
}
