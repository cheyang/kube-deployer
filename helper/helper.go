package helper

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"text/template"

	"github.com/cheyang/fog/persist"
	"github.com/cheyang/fog/types"
)

const Root string = "kube_aliyun"

func RenderTemplateToFile(content string, file *os.File, data interface{}) error {
	tmpl := template.New(file.Name())
	tmpl.Delims("[[[", "]]]")
	tmpl, err := tmpl.Parse(content)
	if err != nil {
		return err
	}
	return tmpl.Execute(file, data)
}

func GetRootDir() string {
	return fmt.Sprintf("/etc/.%s", Root)
}

func ParseScaleFlag(s string) (num int, err error) {
	re := regexp.MustCompile("^[\\+\\-](\\d+)$")
	if re.MatchString(s) {
		num, err = strconv.Atoi(s)
	} else {
		err = fmt.Errorf("scale parameter %s is illegal.", s)
	}

	return num, err
}

func GetCurrentHosts(s persist.Store) ([]types.Host, error) {
	hostList, hostInErrs, err := persist.LoadAllHosts(s)
	if err != nil {
		return hostList, err
	}

	for _, e := range hostInErrs {
		if e != nil {
			return hostList, e
		}
	}

	return hostList, err
}
