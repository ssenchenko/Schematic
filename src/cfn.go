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

type Dict = map[string]any

func OneStepAtATime(step string, fragment Dict, schema Dict) ([]Dict, error) {
	if data, ok := fragment[PROPS]; ok { // type == object with properties
		props := data.(Dict)
		if result, ok := props[step]; ok {
			nextDestination := result.(Dict)
			if _, ok := nextDestination[TYPE]; ok {
				return []Dict{nextDestination}, nil
			}
			if _, ok := nextDestination[REF]; ok {
				return ResolveRef(nextDestination[REF].(string), schema)
			}
			// if it's not TYPE nor REF, it should be ONE_OF or ANY_OF,
			// if not, ResolveAnyOneOf will take care of the exception
			return ResolveAnyOneOf(nextDestination, schema)
		} 
		// it makes sense because of anyOf/oneOf, we might hit a branch which is not on the path
		return make([]Dict, 0), nil

	// in cases below step might be found deeper, in array or in $ref
	} else if data, ok := fragment[TYPE]; ok {

	}
	
}

// Resolve ref and return underlying type.
func ResolveRef(ref string, schema Dict) ([]Dict, error) {
	refTypeName, err := ExtractRefTypeName(ref)
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
	// if it's not TYPE nor REF, it should be ONE_OF or ANY_OF,
	// if not, ResolveAnyOneOf will take care of the exception
	return ResolveAnyOneOf(nextDestination, schema)
}

// Resolve anyOf and oneOf keys.
func ResolveAnyOneOf(withAnyOneOf Dict, schema Dict) ([]Dict, error) {
	var resolved []Dict

	anyOf, anyOk := withAnyOneOf[ANY_OF]
	oneOf, oneOk := withAnyOneOf[ONE_OF]
	if !anyOk && !oneOk {
		return nil, CreateError(withAnyOneOf, fmt.Sprintf("no %s or %s found", ANY_OF, ONE_OF))
	}
	var anyOneOf []Dict
	if anyOk {
		anyOneOf = anyOf.([]Dict)
	} else {
		anyOneOf = oneOf.([]Dict)
	}

	for _, branch := range anyOneOf {
		if _, ok := branch[TYPE]; ok {
			resolved = append(resolved, branch)
		} else if _, ok := branch[REF]; ok {
			resolvedRef, err := ResolveRef(branch[REF].(string), schema)
			if err != nil {
				return nil, err
			}
			resolved = append(resolved, resolvedRef...)
		} else {
			// if it's not TYPE nor REF, it should be ONE_OF or ANY_OF,
			// if not, ResolveAnyOneOf will take care of the exception
			resolvedAnyOf, err := ResolveAnyOneOf(branch, schema)
			if err != nil {
				return nil, err
			}
			resolved = append(resolved, resolvedAnyOf...)
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
func CreateError(fragment Dict, message string) error {
	binary, err := json.Marshal(fragment)
	if err != nil {
		return fmt.Errorf("damn, I can't marshal the fragment after error %v", err)
	}
	return fmt.Errorf("%s %v", message, string(binary))
}
