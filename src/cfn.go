package schematic

import (
	"encoding/json"
	"fmt"
	"slices"
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

type Dict = map[string]any

func OneStepAtATime(propertyName string, fragment Dict, schema Dict) ([]Dict, error) {
	if data, ok := fragment[PROPS]; ok { // type == object with properties
		props := data.(Dict)
		if result, ok := props[propertyName]; ok {
			nextDestination := result.(Dict)
			if _, ok := nextDestination[TYPE]; ok {
				return []Dict{nextDestination}, nil
			}
			if _, ok := nextDestination[REF]; ok {
				return ResolveRef(nextDestination[REF].(string), schema)
			}
			if anyOneOf, ok := extractAnyOneOf(nextDestination); ok {
				return ResolveAnyOneOf(anyOneOf, schema)
			}
			return nil, CreateError(nextDestination, "no idea how to handle")
		}
		// it makes sense because of anyOf/oneOf, we might hit a branch which is not on the path
		return make([]Dict, 0), nil
	}

	// in cases below, propertyName might be found deeper, in array or in $ref

	if type_, ok := fragment[TYPE]; ok {
		// we get here if previous step was an array and next one is supposed to be in an item of object type
		stringTypes, isString := type_.(string)
		// I found 1 case when "type" == ["array", "string"] which looks like a mistake to me,
		// and it should probably be "oneOf/anyOf" or just "array" instead but
		// for the time being, it's faster to add this condition
		// than reach out to the owner and ask for a change
		arrayTypes, isArray := type_.([]string)
		if isArray && slices.Contains(arrayTypes, ARRAY) || isString && stringTypes == ARRAY {
			// keep looking for the same step in the items
			return OneStepAtATime(propertyName, fragment[ITEMS].(Dict), schema)
		}
		// no way to find step, type is neither typed object nor array
		return make([]Dict, 0), nil
	}

	// we get here if the previous step led us to array which 'item' is '$ref'
	if ref, ok := fragment[REF]; ok {
		// first resolve ref
		types, err := ResolveRef(ref.(string), schema)
		if err != nil {
			return nil, err
		}
		// then keep looking for the step in the resolved type(s)
		results := make([]Dict, 1)
		for _, t := range types {
			res, err := OneStepAtATime(propertyName, t, schema)
			if err != nil {
				return nil, err
			}
			results = append(results, res...)
		}
		return results, nil
	}

	// we get here if the previous step led us to array which 'item' is oneOf/anyOf
	if anyOneOf, ok := extractAnyOneOf(fragment[ITEMS].(Dict)); ok {
		// resolve to an actual type(s)
		types, err := ResolveAnyOneOf(anyOneOf, schema)
		if err != nil {
			return nil, err
		}
		// look for our property name
		results := make([]Dict, 1)
		for _, t := range types {
			res, err := OneStepAtATime(propertyName, t, schema)
			if err != nil {
				return nil, err
			}
			results = append(results, res...)
		}
		return results, nil
	}

	return nil, CreateError(fragment, "no idea how to handle")
}

// Resolve ref and return underlying type.
func ResolveRef(ref string, schema Dict) ([]Dict, error) {
	refTypeName, err := extractRefTypeName(ref)
	if err != nil {
		return nil, err
	}
	if _, ok := schema[DEFINITIONS]; !ok {
		return nil, CreateError(schema, "no definitions found in the schema")
	}
	if _, ok := schema[DEFINITIONS].(Dict)[refTypeName]; !ok {
		return nil, CreateError(
			schema[DEFINITIONS].(Dict),
			fmt.Sprintf("no type found for %s in", refTypeName),
		)
	}

	nextDestination, ok := schema[DEFINITIONS].(Dict)[refTypeName].(Dict)
	if !ok {
		return nil, fmt.Errorf("definitions[%s] is not a map[string]any", refTypeName)
	}

	if _, ok := nextDestination[TYPE]; ok {
		return []Dict{nextDestination}, nil
	}
	if _, ok := nextDestination[REF]; ok {
		return ResolveRef(nextDestination[REF].(string), schema)
	}
	if anyOneOf, ok := extractAnyOneOf(nextDestination); ok {
		return ResolveAnyOneOf(anyOneOf, schema)
	}
	return nil, CreateError(nextDestination, "no idea how to handle")
}

// Resolve anyOf and oneOf keys.
func ResolveAnyOneOf(anyOneOf []Dict, schema Dict) ([]Dict, error) {
	var resolved []Dict

	for _, branch := range anyOneOf {
		if _, ok := branch[TYPE]; ok {
			resolved = append(resolved, branch)
		} else if _, ok := branch[REF]; ok {
			resolvedRef, err := ResolveRef(branch[REF].(string), schema)
			if err != nil {
				return nil, err
			}
			resolved = append(resolved, resolvedRef...)
		} else if nestedAnyOne, ok := extractAnyOneOf(branch); ok {
			resolvedAnyOf, err := ResolveAnyOneOf(nestedAnyOne, schema)
			if err != nil {
				return nil, err
			}
			resolved = append(resolved, resolvedAnyOf...)
		} else {
			return nil, CreateError(branch, "no idea how to handle")
		}
	}
	return resolved, nil
}

// Extracts the type name from "#/definitions/<TypeName>".
func extractRefTypeName(ref string) (string, error) {
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
func CreateError(fragment Dict, message string) error {
	binary, err := json.Marshal(fragment)
	if err != nil {
		return fmt.Errorf("damn, I can't marshal the fragment after error %v", err)
	}
	return fmt.Errorf("%s %v", message, string(binary))
}

// Extracts anyOf or oneOf from a fragment.
func extractAnyOneOf(fragment Dict) ([]Dict, bool) {
	if anyOneOf, ok := fragment[ANY_OF]; ok {
		return anyOneOf.([]Dict), true
	}
	if anyOneOf, ok := fragment[ONE_OF]; ok {
		return anyOneOf.([]Dict), true
	}
	return nil, false
}
