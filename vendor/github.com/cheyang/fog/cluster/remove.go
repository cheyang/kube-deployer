package cluster

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/cheyang/fog/persist"
)

func Remove(storePath string, s persist.Store) error {

	hostList, _, err := persist.LoadAllHosts(s)
	if err != nil {
		return err
	}

	for _, host := range hostList {
		logrus.Infof("To remove %s", host.Name)
		err := host.Driver.Remove()
		if err != nil {
			return err
		}
		err = s.Remove(host.Name)
		if err != nil {
			logrus.Infof("host err: %v", err)
		}
	}
	return os.RemoveAll(storePath)
}
