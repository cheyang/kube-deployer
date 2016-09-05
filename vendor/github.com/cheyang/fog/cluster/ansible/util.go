package ansible

import (
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/cheyang/fog/types"
)

type byHostName []types.Host

func (s byHostName) Len() int {
	return len(s)
}
func (s byHostName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byHostName) Less(i, j int) bool {
	sai := strings.Split(s[i].Name, "-")
	si, err := strconv.Atoi(sai[len(sai)-1])
	if err != nil {
		logrus.Infof("err: %v", err)
	}
	saj := strings.Split(s[j].Name, "-")
	sj, err := strconv.Atoi(saj[len(saj)-1])
	if err != nil {
		logrus.Infof("err: %v", err)
	}
	return si < sj
}
