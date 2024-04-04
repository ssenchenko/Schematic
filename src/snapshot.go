package schematic

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

const (
	SNAPSHOT_DIR string = "./snapshots"
	TEMPLATE_DIR string = "./templates"
)

func hydrateResourceStruct(resource ResourceType, templateFileName string) (string, error) {
	tmpl, err := template.
		New(templateFileName).
		ParseFiles(fmt.Sprintf("%s/%s", TEMPLATE_DIR, templateFileName))
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	err = tmpl.Execute(&out, resource)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func loadSnapshot(snapshotFileName string) (string, error){
	content, err := os.ReadFile(
		fmt.Sprintf("%s/%s", SNAPSHOT_DIR, snapshotFileName))
	if err != nil {
		return "", err
	}
	return string(content), nil
}
