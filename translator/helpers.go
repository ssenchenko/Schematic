package translator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	SnapshotDir string = "./snapshots"
)

// hydrateTemplate populates templates with provided data.
func hydrateTemplate[T any](
	templateData T,
	templateName string,
	templateDir string,
	funcs template.FuncMap,
	nestedTemplates ...string,
) (bytes.Buffer, error) {
	tmpl := template.New(templateName)

	if len(funcs) != 0 {
		tmpl.Funcs(funcs)
	}

	tmpl, err := tmpl.ParseFiles(filepath.Join(templateDir, templateName))
	if err != nil {
		return bytes.Buffer{}, err
	}

	if len(nestedTemplates) != 0 {
		mockedTemplates := make([]string, 0, len(nestedTemplates))
		fullTemplates := make([]string, 0, len(nestedTemplates))
		for _, nestedTemplate := range nestedTemplates {
			if strings.HasPrefix(nestedTemplate, "{{") {
				mockedTemplates = append(mockedTemplates, nestedTemplate)
			} else {
				fullTemplates = append(fullTemplates, nestedTemplate)
			}
		}

		for _, template := range mockedTemplates {
			tmpl, err = tmpl.Parse(template)
			if err != nil {
				return bytes.Buffer{}, err
			}
		}

		if len(fullTemplates) != 0 {
			for i := 0; i < len(fullTemplates); i++ {
				fullTemplates[i] = filepath.Join(templateDir, fullTemplates[i])
			}

			tmpl, err = tmpl.ParseFiles(fullTemplates...)
			if err != nil {
				return bytes.Buffer{}, err
			}
		}
	}

	var out bytes.Buffer
	err = tmpl.Execute(&out, templateData)
	if err != nil {
		return bytes.Buffer{}, err
	}
	return out, nil
}

// LoadSnapshot loads snapshot for snapshot testing.
func LoadSnapshot(snapshotFileName string) (string, error) {
	content, err := os.ReadFile(
		fmt.Sprintf("%s/%s", SnapshotDir, snapshotFileName))
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// LoadAllRelationships loads all relationships from a file.
func LoadAllRelationships(relationshipFileRelativePath string) (AllRelationships, error) {
	relationshipsFile, err := GetFullPath(relationshipFileRelativePath)
	if err != nil {
		return AllRelationships{}, err
	}
	content, err := os.ReadFile(relationshipsFile)
	if err != nil {
		return AllRelationships{}, err
	}

	var allRelationships AllRelationships
	err = json.Unmarshal(content, &allRelationships)
	if err != nil {
		return AllRelationships{}, err
	}
	return allRelationships, nil
}

func LoadCfnSchemaCombined(cfnSchemaDir string, filter []string) (map[string]map[string]any, error) {
	cfnSchemaCombined := make(map[string]map[string]any)

	cfnSchemaDir, err := GetFullPath(cfnSchemaDir)
	if err != nil {
		return cfnSchemaCombined, err
	}
	cfnFiles, err := os.ReadDir(cfnSchemaDir)
	if err != nil {
		return cfnSchemaCombined, err
	}
	var cfnFileNames []string
	if filter != nil {
		cfnFileNames = filter
	} else {
		for _, file := range cfnFiles {
			cfnFileNames = append(cfnFileNames, file.Name())
		}
	}

	for _, fileName := range cfnFileNames {
		content, err := os.ReadFile(filepath.Join(cfnSchemaDir, fileName))
		if err != nil {
			return cfnSchemaCombined, err
		}
		cfnJsonSchema := make(map[string]any)
		err = json.Unmarshal(content, &cfnJsonSchema)
		if err != nil {
			return cfnSchemaCombined, err
		}
		resourceTypeName, _ := cfnJsonSchema["typeName"].(string)
		cfnSchemaCombined[resourceTypeName] = cfnJsonSchema
	}
	return cfnSchemaCombined, nil
}

// GetFullPath transforms relative path to an absolute one
// using current working directory as a base
func GetFullPath(path string) (string, error) {
	if !filepath.IsAbs(path) {
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return filepath.Join(wd, path), nil
	}
	return path, nil
}

// Deref helper to dereference possibly pointers in template
func Deref[T any](pointer *T) T {
	if pointer == nil {
		var zero T
		return zero
	}
	return *pointer
}
