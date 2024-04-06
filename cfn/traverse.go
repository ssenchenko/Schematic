package cfn

import (
	"encoding/json"
	"fmt"
	"log"
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

type TypeBits uint8

const (
	BOOLB TypeBits = 1 << iota
	INTB
	NUMBERB
	STRINGB
	ARRAYB
	OBJECTB
)

type Dict map[string]any

func IsArray(pathToProperty string, typeNameCfnFormat string, schema map[string]Dict) (bool, error) {
	types, err := FollowPath(pathToProperty, schema[typeNameCfnFormat])
	if err != nil {
		return false, err
	}
	if len(types) == 0 {
		return false, fmt.Errorf("no type found for %s in %s", pathToProperty, typeNameCfnFormat)
	}
	arrayTest := make([]bool, len(types))
	counter := 0
	allAreSame := true
	for typeBit := range types {
		arrayTest[counter] = typeBit&ARRAYB == ARRAYB
		if arrayTest[counter] != arrayTest[0] {
			allAreSame = false
			break
		}
		counter++
	}
	if !allAreSame {
		return false, fmt.Errorf(
			"what should I do with it? %s in %s has branches; some of them end with array and some - don't", pathToProperty, typeNameCfnFormat)
	}
	return arrayTest[0], nil
}

func FollowPath(path string, schema Dict) (map[TypeBits]bool, error) {
	propertyNames := strings.Split(path, "/")
	var propertyNumberInPath int = 0
	type StackNode struct {
		propertyNumberInPath int
		fragment             Dict
		type_                TypeBits
	}
	stack := make([]StackNode, 1, 5) // its capacity is hard to predict, but it's better than 0 or 1
	stack[0] = StackNode{propertyNumberInPath: propertyNumberInPath, fragment: schema, type_: OBJECTB}
	types_seen := make(map[TypeBits]bool)
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if node.propertyNumberInPath == len(propertyNames) {
			types_seen[node.type_] = true
			continue
		}
		fragments, err := OneStepAtATime(propertyNames[node.propertyNumberInPath], node.fragment, schema)
		if err != nil {
			return nil, fmt.Errorf("error in path %s %v", path, err)
		}
		if len(fragments) == 0 {
			log.Printf("step %s not found in %s", propertyNames[node.propertyNumberInPath], node.fragment)
			continue
		}
		for _, fragment := range fragments {
			nodeType := node.type_
			// keeping track of types to be able to check if array happened on the "correct" path
			// sometimes "type" is an array like ["object", "string"] or ["array", "string"]
			// which are probably errors, but we need to handle those without breaking
			// as far as we concerned, we need to know if any of the types along the path
			// can be an array, so let's bit-OR all the types from the array
			switch fragment[TYPE].(type) {
			case []string: // type == ["object", "string"] or ["array", "string"]
				for _, t := range fragment[TYPE].([]string) {
					type_, err := typeBitsFromJsonType(t)
					if err != nil {
						return nil, fmt.Errorf("error in path %s %v", path, err)
					}
					nodeType |= type_
				}
			case string:
				type_, err := typeBitsFromJsonType(fragment[TYPE].(string))
				if err != nil {
					return nil, fmt.Errorf("error in path %s %v", path, err)
				}
				nodeType |= type_
			default:
				return nil, fmt.Errorf("error in path %s unknown type %v", path, fragment)
			}
			stack = append(stack, StackNode{propertyNumberInPath: node.propertyNumberInPath + 1, fragment: fragment, type_: nodeType})
		}
	}
	return types_seen, nil
}

func OneStepAtATime(propertyName string, fragment Dict, schema Dict) ([]Dict, error) {
	if data, ok := fragment[PROPS]; ok { // type == object with properties
		props := data.(Dict)
		if result, ok := props[propertyName]; ok {
			nextDestination := result.(Dict)
			if _, ok := nextDestination[TYPE]; ok {
				return []Dict{nextDestination}, nil
			}
			if refString, ok := nextDestination[REF]; ok {
				return ResolveRef(refString.(string), schema)
			}
			if anyOneOf, ok := extractAnyOneOf(nextDestination); ok {
				return ResolveAnyOneOf(anyOneOf, schema)
			}
			return nil, fmt.Errorf("no idea how to handle %s", result)
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
		results := make([]Dict, 0, 1)
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
	if anyOneOf, ok := extractAnyOneOf(fragment); ok {
		// resolve to an actual type(s)
		types, err := ResolveAnyOneOf(anyOneOf, schema)
		if err != nil {
			return nil, err
		}
		// look for our property name
		results := make([]Dict, 0, 1)
		for _, t := range types {
			res, err := OneStepAtATime(propertyName, t, schema)
			if err != nil {
				return nil, err
			}
			results = append(results, res...)
		}
		return results, nil
	}

	return nil, fmt.Errorf("no idea how to handle %s", fragment)
}

// Resolve ref and return underlying type.
func ResolveRef(ref string, schema Dict) ([]Dict, error) {
	refTypeName, err := extractRefTypeName(ref)
	if err != nil {
		return nil, err
	}
	if _, ok := schema[DEFINITIONS]; !ok {
		return nil, fmt.Errorf("no definitions found in the schema %s", schema)
	}
	if _, ok := schema[DEFINITIONS].(Dict)[refTypeName]; !ok {
		return nil, fmt.Errorf("no type found for %s in %s", refTypeName, schema[DEFINITIONS].(Dict))
	}

	nextDestination := schema[DEFINITIONS].(Dict)[refTypeName].(Dict)
	if _, ok := nextDestination[TYPE]; ok {
		return []Dict{nextDestination}, nil
	}
	if refString, ok := nextDestination[REF]; ok {
		return ResolveRef(refString.(string), schema)
	}
	if anyOneOf, ok := extractAnyOneOf(nextDestination); ok {
		return ResolveAnyOneOf(anyOneOf, schema)
	}
	return nil, fmt.Errorf("no idea how to handle %s", nextDestination)
}

// Resolve anyOf and oneOf keys.
func ResolveAnyOneOf(anyOneOf []Dict, schema Dict) ([]Dict, error) {
	var resolved []Dict

	for _, branch := range anyOneOf {
		if _, ok := branch[TYPE]; ok {
			resolved = append(resolved, branch)
		} else if refString, ok := branch[REF]; ok {
			resolvedRef, err := ResolveRef(refString.(string), schema)
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
			return nil, fmt.Errorf("no idea how to handle %s", branch)
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

func (fragment Dict) String() string {
	binary, err := json.Marshal(fragment)
	if err != nil {
		return fmt.Sprintf("damn, I can't marshal the fragment after error %v", err)
	}
	return string(binary)
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

func typeBitsFromJsonType(jsonType string) (TypeBits, error) {
	switch jsonType {
	case BOOL:
		return BOOLB, nil
	case INT:
		return INTB, nil
	case NUMBER:
		return NUMBERB, nil
	case STRING:
		return STRINGB, nil
	case ARRAY:
		return ARRAYB, nil
	case OBJECT:
		return OBJECTB, nil
	default:
		return 0, fmt.Errorf("unexpected json type: %s", jsonType)
	}
}
