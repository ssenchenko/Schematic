package schematic

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	TYPE_NAME      string = "typeName"
	REF            string = "$ref"
	TYPE           string = "type"
	ENUM           string = "enum"
	PROPS          string = "properties"
	DEFINITIONS    string = "definitions"
	ITEMS          string = "items"
	ONE_OF         string = "oneOf"
	ANY_OF         string = "anyOf"
	ARRAY          string = "array"
	OBJECT         string = "object"
	STRING         string = "string"
	NUMBER         string = "number"
	BOOL           string = "boolean"
	INT            string = "integer"
	UNTYPED_OBJECT string = "untyped-object" // untyped object
)

// Resolve ref and return underlying type.
func ResolveRef(ref string, schema map[string]any) ([]map[string]any, error) {
	refTypeName, err := ExtractRefTypeName(ref)
	if err != nil {
		return nil, err
	}
	if _, ok := schema[DEFINITIONS]; !ok {
		return nil, CreateError(schema, "no definitions found in the schema")
	}
	if _, ok := schema[DEFINITIONS].(map[string]any)[refTypeName]; !ok {
		return nil, CreateError(schema, fmt.Sprintf("no type found for %s in", refTypeName))
	}

	nextDestination, ok := schema[DEFINITIONS].(map[string]any)[refTypeName].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("definitions[%s] is not a map[string]any", refTypeName)
	}
	if _, ok := nextDestination[TYPE]; ok {
		return []map[string]any{nextDestination}, nil
	}
	if _, ok := nextDestination[REF]; ok {
		return ResolveRef(nextDestination[REF].(string), schema)
	}
	// TODO: resolve anyOf and oneOf

	return nil, CreateError(schema, "oh no, another uncharted terrain")
}

// Resolve anyOf and oneOf keys.
func ResolveAnyOneOf(anyOneOf []map[string]any, schema map[string]any) ([]map[string]any, error) {
	var resolved []map[string]any

	for i, branch := range anyOneOf {
		if _, ok := branch[TYPE]; ok {
			resolved = append(resolved, branch)
		} else if _, ok := branch[REF]; ok {
			resolvedRef, err := ResolveRef(branch[REF].(string), schema)
			if err != nil {
				return nil, CreateError(anyOneOf[i], err.Error())
			}
			resolved = append(resolved, resolvedRef...)
		} else if _, ok := branch[ANY_OF]; ok {
			resolvedAnyOf, err := ResolveAnyOneOf(branch[ANY_OF].([]map[string]any), schema)
			if err != nil {
				return nil, CreateError(anyOneOf[i], err.Error())
			}
			resolved = append(resolved, resolvedAnyOf...)
		} else if _, ok := branch[ONE_OF]; ok {
			resolvedOneOf, err := ResolveAnyOneOf(branch[ONE_OF].([]map[string]any), schema)
			if err != nil {
				return nil, CreateError(anyOneOf[i], err.Error())
			}
			resolved = append(resolved, resolvedOneOf...)
		} else {
			return nil, CreateError(anyOneOf[i], "unexpected branch")
		}

	}
	return resolved, nil
}

// Extracts the type name from "#/definitions/<TypeName>".
func ExtractRefTypeName(ref string) (string, error) {
	err := fmt.Errorf("unexpected ref: %s", ref)

	if ref == "" {
		return "", err
	}

	refParts := strings.Split(ref, "/")
	if len(refParts) != 3 ||
		refParts[0] != "#" ||
		refParts[1] != DEFINITIONS ||
		refParts[2] == "" {
		return "", err
	}

	return refParts[2], nil
}

// Create error message with arbitrary json fragment output.
func CreateError(fragment map[string]any, message string) error {
	binary, err := json.MarshalIndent(fragment, "", "  ")
	if err != nil {
		return fmt.Errorf("damn, I can't marshal the fragment after error %v", err)
	}
	return fmt.Errorf("%s\n%v", message, string(binary))
}
