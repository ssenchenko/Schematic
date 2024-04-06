package translator

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

const (
	SNAPSHOT_DIR string = "./snapshots"
)

// Hydrate template for snapshot testing.
func hydrate[T any](templateData T, templateNames ...string) (string, error) {
	if len(templateNames) == 0 {
		return "", fmt.Errorf("template name is required")
	}

	templateFileNames := make([]string, len(templateNames))
	for i, name := range templateNames {
		templateFileNames[i] = fmt.Sprintf("%s/%s", TEMPLATE_DIR, name)
	}

	tmpl, err := template.New(templateNames[0]).
		Funcs(template.FuncMap{"DerefResourceUnion": Deref[ResourceUnion]}).
		ParseFiles(templateFileNames...)
	if err != nil {
		return "", err
	}

	var out bytes.Buffer
	err = tmpl.Execute(&out, templateData)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// Load snapshot for snapshot testing.
func loadSnapshot(snapshotFileName string) (string, error) {
	content, err := os.ReadFile(
		fmt.Sprintf("%s/%s", SNAPSHOT_DIR, snapshotFileName))
	if err != nil {
		return "", err
	}
	return string(content), nil
}
