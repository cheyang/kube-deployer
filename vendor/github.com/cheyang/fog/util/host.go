package util

import (
	"regexp"
	"sort"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/cheyang/fog/types"
)

const (
	splitHostname = "(.+)-(\\d+)"
)

func BuildRunningMap(hosts []types.Host) (runningHostMap map[string][]string, err error) {
	runningHostMap = make(map[string][]string)

	for _, host := range hosts {
		// build running host map
		hostname := host.Name
		key, _, err := ParseHostname(hostname)
		if err != nil {
			return runningHostMap, err
		}

		if _, found := runningHostMap[key]; !found {
			runningHostMap[key] = []string{}
		}

		runningHostMap[key] = append(runningHostMap[key], hostname)
	}

	for _, v := range runningHostMap {
		sort.Sort(byHostname(v))
	}
	return runningHostMap, nil
}

func ParseHostname(hostname string) (specName string, id int, err error) {
	re := regexp.MustCompile(splitHostname)
	match := re.FindStringSubmatch(hostname)
	specName = match[1]
	id, err = strconv.Atoi(match[2])
	return specName, id, err
}

type byHostname []string

func (s byHostname) Len() int {
	return len(s)
}
func (s byHostname) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byHostname) Less(i, j int) bool {
	_, si, err := ParseHostname(s[i])
	if err != nil {
		logrus.Infof("err: %v", err)
	}
	_, sj, err := ParseHostname(s[j])
	if err != nil {
		logrus.Infof("err: %v", err)
	}
	return si < sj
}
