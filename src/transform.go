package schematic

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Map which better fits the needs of types generation.
// Output of the transformer.
// Looks like:
//
//	{
//		Resource: {
//			F1: Struct1
//			F2: string  # Xxxx
//			F3: Array/Struct1 # A1
//			F4: Array/integer
//			F5: Enum1
//		},
//		Xxxx: string # type def
//		Struct1: {
//			F1: String
//			F2: integer
//			F3: Array/Enum2
//		},
//		Enum1: [E1, E2],
//		Enum2: [E3, E4]
//		A1: Array/Struct1 # type def
//	}
type TypesMap = map[string]any

// Definitions are the JSON schema definitions.
type Definitions = map[string]map[string]any

// Errors for the transformer
type DefinitionErrors struct {
	Errors []DefinitionError
}
type DefinitionError struct {
	DefinitionName string
	Message        *string
	PropertyErrors []PropertyError
}
type PropertyError struct {
	PropertyName string
	Messsage     string
}

// JSON schema types
type JsonSchemaType = string

const (
	String  JsonSchemaType = "string"
	Integer JsonSchemaType = "integer"
	Number  JsonSchemaType = "number"
	Object  JsonSchemaType = "object"
	Array   JsonSchemaType = "array"
	Boolean JsonSchemaType = "boolean"
	Null    JsonSchemaType = "null" // not sure how to react on it if it appears
)

const SerializedObject string = "string/object" // should it be a separate type?

// Transform Json schema properties and definitions to TypesMap.
// resourceName + properties -> just another type definition
//
// properties and definitions are JSON schema parts under same name keys
func ToTypesMap(
	resourceName string,
	definitions Definitions,
	properties map[string]any,
) (TypesMap, error) {
	var typesMap = make(TypesMap, len(definitions)+1)
	var allErrors = DefinitionErrors{Errors: make([]DefinitionError, 0, len(definitions))}

	// add resource to the definitions since it's not any different
	resource := map[string]any{"type": "object", "properties": properties}
	definitions[resourceName] = resource

	//initialized typesMap to be able to dereference all $ref's
	var initErrors DefinitionErrors
	initErrors = initializeTypesMap(definitions, typesMap)
	if len(initErrors.Errors) > 0 {
		allErrors.Errors = append(allErrors.Errors, initErrors.Errors...)
	}

	for definitionName, definition := range definitions {
		var populateErrors DefinitionErrors
		populateErrors = populateTypesMap(definition, typesMap, definitionName, "")
		if len(populateErrors.Errors) > 0 {
			allErrors.Errors = append(allErrors.Errors, populateErrors.Errors...)
		}
	}

	return typesMap, allErrors
}

// Do recursive depth-first search to populate typesMap with dereferenced types.
func populateTypesMap(json map[string]any, typesMap TypesMap, typeName string, propertyName string) DefinitionErrors {
	var allErrors = DefinitionErrors{Errors: make([]DefinitionError, 0, len(json))}

	if jsonType, found := json["type"]; found {
		switch jsonType.(type) {
		case string:
			switch jsonType.(string) {
			case Object:
				// new type definition found in the property
				if propertyName != "" {
					if _, found := json["properties"]; found {
						// property will have a type of the same name
						typesMap[typeName].(map[string]string)[propertyName] = propertyName
						// add that type to the typesMap
						typesMap[typeName] = make(map[string]string)
					} else {
						// no need  to create a new type
						typesMap[typeName].(map[string]string)[propertyName] = SerializedObject
					}
					// keep discovering that new type
					typeName = propertyName
					// its properties are not known yet
					propertyName = ""
				}
				if properties, found := json["properties"]; found {
					for propertyName, property := range properties.(map[string]any) {
						populateErrors :=
							populateTypesMap(property.(map[string]any), typesMap, typeName, propertyName)
						if len(populateErrors.Errors) > 0 {
							allErrors.Errors = append(allErrors.Errors, populateErrors.Errors...)
						}
					}
				} 
			case Array:
				// do nothing

			case String:
			case Integer, Number, Boolean, Null:
				// do nothing
			default:
				message := fmt.Sprintf("unknown type: %s", jsonType.(string))
				if propertyName != "" {
					allErrors.Errors = append(allErrors.Errors, DefinitionError{
						DefinitionName: typeName, 
						PropertyErrors: []PropertyError{{PropertyName: propertyName, Messsage: message}}})
				} else {
					allErrors.Errors = append(allErrors.Errors, DefinitionError{DefinitionName: typeName, Message: &message})
				}
			}
		case []string:
			if reflect.DeepEqual(jsonType, []JsonSchemaType{Object, String}) {
				// do nothing
			}
		default:
			message := fmt.Sprintf("unknown type: %s", jsonType)
			allErrors.Errors = append(allErrors.Errors, DefinitionError{DefinitionName: typeName, Message: &message})
		}
	} else if ref, found := json["$ref"]; found {
		// dereference ref
		refType := ref.(string)[14:] // remove "#/definitions/"
		if _, found := typesMap[refType]; found {
			if propertyName != "" {
				typesMap[typeName].(map[string]string)[propertyName] = refType
			} 
		} else {
			message := fmt.Sprintf("unknown ref: %s", ref.(string))
			if propertyName != "" {
				allErrors.Errors = append(allErrors.Errors, DefinitionError{
					DefinitionName: typeName, 
					PropertyErrors: []PropertyError{{PropertyName: propertyName, Messsage: message}}})
			} else {
				allErrors.Errors = append(allErrors.Errors, DefinitionError{DefinitionName: typeName, Message: &message})
			}
		}
	} else {
		// do nothing

	}

	return allErrors
}

func initObjectOnTypesMap(json map[string]any, typesMap map[string]any, typeName string) {
	if _, found := json["properties"]; found {
		typesMap[typeName] = make(map[string]string)
	} else {
		typesMap[typeName] = SerializedObject
	}
}

// Initialize typesMap with types from definitions.
func initializeTypesMap(definitions Definitions, typesMap TypesMap) DefinitionErrors {
	var allErrors = DefinitionErrors{Errors: make([]DefinitionError, 0, len(definitions))}

	for typeName, typeInfo := range definitions {
		jsonType, found := typeInfo["type"]
		if !found {
			message := "type is missing"
			allErrors.Errors = append(allErrors.Errors, DefinitionError{DefinitionName: typeName, Message: &message})
		}
		switch jsonType.(type) {
		case string:
			switch jsonType.(string) {
			case Object:
				if _, found := typeInfo["properties"]; found {
					typesMap[typeName] = make(map[string]string)
				} else {
					typesMap[typeName] = SerializedObject
				}
			case Array:
				typesMap[typeName] = "array/"
			case String, Integer, Number, Boolean, Null:
				typesMap[typeName] = jsonType.(string)
			default:
				message := fmt.Sprintf("unknown type: %s", jsonType.(string))
				allErrors.Errors = append(allErrors.Errors, DefinitionError{DefinitionName: typeName, Message: &message})
			}
		case []string:
			if reflect.DeepEqual(jsonType, []JsonSchemaType{Object, String}) {
				typesMap[typeName] = SerializedObject // meaning it is a serialized object
			}
		default:
			message := fmt.Sprintf("unknown type: %s", jsonType)
			allErrors.Errors = append(allErrors.Errors, DefinitionError{DefinitionName: typeName, Message: &message})
		}
	}
	return allErrors
}

// Error for PropertyError
// object is supposed to be small enough to bee passed by copy
func (e PropertyError) Error() string {
	// use marshalling to get a string representation of the value
	binary, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("Property %s: %s", e.PropertyName, err.Error())
	}
	return fmt.Sprintf("%s", string(binary))
}

// Error for DefinitionError
// object is supposed to be small enough to bee passed by copy
func (e DefinitionError) Error() string {
	binary, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("Definition %s: %s", e.DefinitionName, err.Error())
	}
	return fmt.Sprintf("%s", string(binary))
}

// Error for DefinitionErrors
// object is supposed to be small enough to bee passed by copy
func (e DefinitionErrors) Error() string {
	binary, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("DefinitionErrors: %s", err.Error())
	}
	return fmt.Sprintf("%s", string(binary))
}
