package cluster

import (
	"github.com/Sirupsen/logrus"
	"github.com/cheyang/fog/persist"
)

func List(s persist.Store) error {
	hostList, hostInErrs, err := persist.LoadAllHosts(s)
	if err != nil {
		return err
	}

	for name, e := range hostInErrs {
		logrus.Infof("%s:%v", name, e)
	}

	for _, host := range hostList {
		logrus.Infof("name: %s", host.Name)
		logrus.Infof("%+v", host)
	}

	return nil
}
