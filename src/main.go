package schematic

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func main() {
	schemaCombined := "data/all-schema-combined.json"
	if len(os.Args) > 1 {
		schemaCombined = os.Args[1]
	}

	if !filepath.IsAbs(schemaCombined) {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		schemaCombined = filepath.Join(wd, schemaCombined)
	}

	content, err := os.ReadFile(schemaCombined)
	if err != nil {
		panic(err)
	}

	var input SchemaCombined
	err = json.Unmarshal(content, &input)
	if err != nil {
		panic(err)
	}

	log.Println(input)
}
