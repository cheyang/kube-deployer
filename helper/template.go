package helper

import (
	"os"
	"text/template"
)

const RootDir string = "/etc/.kube_aliyun/"

func RenderTemplateToFile(content string, file *os.File, data interface{}) error {
	tmpl := template.New(file.Name())
	tmpl.Delims("[[[", "]]]")
	tmpl, err := tmpl.Parse(content)
	if err != nil {
		return err
	}
	return tmpl.Execute(file, data)
}
