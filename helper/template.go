package helper

import (
	"fmt"
	"os"
	"text/template"
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
