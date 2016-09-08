package cluster

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/cheyang/fog/persist"
	"github.com/cheyang/fog/util"
)

func Scalein(s persist.Store, scaleInMap map[string]int) error {
	hostList, _, err := persist.LoadAllHosts(s)
	if err != nil {
		return err
	}

	runningHostMap, err := util.BuildRunningMap(hostList)
	if err != nil {
		return err
	}

	for k, v := range scaleInMap {
		if list, found := runningHostMap[k]; found {
			start := len(list) - v
			if start < 0 {
				start = 0
			}
			for i := start; i < len(list)-1; i++ {
				name := list[i]
				logrus.Infof("To Remove %s from %s", name, k)
				host, err := s.Load(name)
				if err != nil {
					return err
				}
				logrus.Infof("%s: %+v", name, host)
				// to remove it later
			}

		} else {
			return fmt.Errorf("failed to find %s in storage, so can't scale it in %d", k, v)
		}
	}

	return nil
}
