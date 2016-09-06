package config

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/cheyang/kube-deployer/config/templates"
)

var parentDir string = "/etc/.kube_aliyun/"

var files = []struct {
	Filename string
	Content  string
}{
	{"aliyun.yaml", templates.AliyunTemplate},
	{"ansible.yaml", templates.AnsibleTemplate},
}

// stage can be create/scale
func GenerateConfigFile(clusterName, stage string, data interface{}) (configFiles []string, err error) {
	// f, err := os.OpenFile("templates.go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer f.Close()
	// data := Data{tmpls, time.Now().UTC()}
	// if err := tmpl.Execute(f, data); err != nil {
	// 	log.Fatal("Failed to render template:", err)
	// }

	configFiles = []string{}
	inputDir := filepath.Join(parentDir, clusterName, "input", stage)

	err = os.MkdirAll(inputDir, 0700)
	if err != nil {
		return configFiles, err
	}

	for _, file := range files {
		absPath := filepath.Join(inputDir, file.Filename)
		_, err = os.Stat(absPath)
		if os.IsNotExist(err) {
			err = nil
		} else {
			return configFiles, fmt.Errorf(`config file %s is created before, 
				please check again to make sure you want to overwrite the cluster %s. 
				If you are sure`, absPath, clusterName)
		}

		f, err := os.OpenFile(absPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return configureFiles, err
		}
		tmpl := template.New(file.Filename)
		tmpl.Delims("[[[", "]]]")
		tmpl, _ = tmpl.Parse(file.Content)
		tmpl.Execute(f, data)
		configFiles = append(configFiles, absPath)
		err = f.Close()
		if err != nil {
			return configFiles, err
		}
	}

	return configFiles, nil
}

func GetClusterInputPath(name, stage string) string {
	return filepath.Join(parentDir, name, "input", stage)
}
