package cluster

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/cheyang/fog/persist"
	"github.com/cheyang/fog/types"
)

type desiredConfig struct {
	instances int
	vmSpec    types.VMSpec
}

func Scale(s persist.Store, desiredMap map[string]int) error {

	hostList, _, err := persist.LoadAllHosts(s)
	if err != nil {
		return err
	}

	var spec *types.Spec
	spec, err = s.LoadSpec()
	if err != nil {
		return nil
	}
	// defer s.SaveSpec(spec)
	// vm spec name: index
	runningMap := make(map[string]types.VMSpec)
	for _, vmSpec := range spec.VMSpecs {
		runningMap[vmSpec.Name] = vmSpec
	}
	// check if the vm spec name is validate
	desiredHostMap := make(map[string]desiredConfig)
	for name, _ := range desiredMap {
		if _, found := runningMap[name]; !found {
			return fmt.Errorf("The name %s doesn't exist in cluster %s, can't be scaled", name, spec.Name)
		}

		desiredHostMap[name] = desiredConfig{
			vmSpec:    runningMap[name],
			instances: desiredMap[name],
		}
	}

	currentHostMap := buildcurrentHostMap(hostList, runningMap)

	toRemoveHosts, toCreateHostSpecs, err := buildScaleList(currentHostMap, desiredHostMap)
	if err != nil {
		return err
	}
	// toRemoveHosts, toCreateHostSpecs, err := buildScaleList(currentHostMap, desiredHostMap)
	logrus.Infof("toRemoveHost: %v\n", toRemoveHosts)
	logrus.Infof("toCreateHostSpecs: %v\n", toCreateHostSpecs)

	return nil
}

func buildScaleList(currentHostMap map[string][]string, desiredHostMap map[string]desiredConfig) ([]string, []types.VMSpec, error) {
	var (
		toRemoveHosts     = make([]string, 0)
		toCreateHostSpecs = make([]types.VMSpec, 0)
	)

	for name, desired := range desiredHostMap {
		desiredNum := desired.instances
		runningNum := len(currentHostMap[name])
		if runningNum > desiredNum {
			//scale in
			// toRemoveHosts = append(toRemoveHosts,)
			for i := desiredNum; i < runningNum; i++ {
				toRemoveHosts = append(toRemoveHosts, currentHostMap[name][i])
			}
		} else if runningNum < desiredNum {
			//scale out
			s := strings.Split(currentHostMap[name][len(currentHostMap[name])-1], "-")
			maxNum, err := strconv.Atoi(s[len(s)-1])
			if err != nil {
				return toRemoveHosts, toCreateHostSpecs, err
			}
			for i := maxNum + 1; i < desiredNum; i++ {
				vm := desired.vmSpec
				vm.Name = fmt.Sprintf("%s-%d", vm.Name, i)
				toCreateHostSpecs = append(toCreateHostSpecs, vm)
			}
		}
	}

	return toRemoveHosts, toCreateHostSpecs, nil
}

// Get the name list
func buildcurrentHostMap(hosts []types.Host, runningMap map[string]types.VMSpec) (currentHostMap map[string][]string) {
	currentHostMap = make(map[string][]string)
	for k, _ := range runningMap {
		currentHostMap[k] = make([]string, 0)
	}

	for _, host := range hosts {
		for name, _ := range runningMap {
			if strings.HasPrefix(host.Name, name) {
				currentHostMap[name] = append(currentHostMap[name], host.Name)
				break
			}
		}
	}

	for _, v := range currentHostMap {
		// logrus.Infof("key: %s value: %+v", k, currentHostMap[k])
		sort.Sort(ByName(v))
		// logrus.Infof("key: %s value: %+v", k, currentHostMap[k])
	}

	return
}

type ByName []string

func (s ByName) Len() int {
	return len(s)
}
func (s ByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByName) Less(i, j int) bool {
	si, err := strconv.Atoi(strings.Split(s[i], "-")[len(strings.Split(s[i], "-"))-1])
	if err != nil {
		logrus.Infof("err: %v", err)
	}
	sj, err := strconv.Atoi(strings.Split(s[j], "-")[len(strings.Split(s[j], "-"))-1])
	if err != nil {
		logrus.Infof("err: %v", err)
	}
	return si < sj
}
