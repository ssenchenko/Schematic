package translator

import "ssenchenko/schematic/tools"

const (
	TEMPLATE_DIR string = "./templates"
)

type Reference struct {
	TypeName  string `json:"typeName"`
	Attribute string `json:"attribute"`
}

type Relationship map[string][]Reference

type Resource struct {
	PrimaryKeys   []string       `json:"primaryIdentifier"`
	Relationships []Relationship `json:"relationships"`
}

type AllRelationships map[string]Resource

type RustModel []ResourceType

type ResourceType struct {
	CfnResourceName     string
	RustResourceName    string
	GraphQlResourceName string
	Properties          []ResourceProperty
	Relationships       []ResourceRelationship
}

type ResourceProperty struct {
	RustPropertyName string
	RustPropertyType string
}

type ResourceRelationship struct {
	RustSourcePropertyName string
	RustReturnType         string
	RustGenericType        string
	TargetUnion            *ResourceUnion
}

type ResourceUnion struct {
	RustUnionName string
	Resources     []ResourceType
}

func (allRelationships *AllRelationships) ApplyOverrides(overrides map[string]map[string]string) {
	for resourceName, overrides := range overrides {
		for oldName, newName := range overrides {
			for _, relationship := range (*allRelationships)[resourceName].Relationships {
				relationship[newName] = relationship[oldName]
				delete(relationship, oldName)
			}
		}
	}
}

func translate(allRelationships AllRelationships, filter map[string]bool) (RustModel, []error) {
	errors := make([]error, 0)
	rustModel := make(RustModel, 0, len(allRelationships))

	for resourceName, resource := range allRelationships {
		if _, ok := filter[resourceName]; !ok {
			continue
		}

		awsResourceName, err := tools.NewAwsResourceName(resourceName)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		properties := []ResourceProperty{
			{
				RustPropertyName: "id",
				RustPropertyType: "String",
			},
			{
				RustPropertyName: "all_properties",
				RustPropertyType: "String",
			},
		}

		relationships := []ResourceRelationship{}
		for _, relationship := range resource.Relationships {
			for relationshipName, references := range relationship {
				
				for _, reference := range references {
					rustPropertyName := PascalCaseToSnakeCase(reference.Attribute)
					rustPropertyType := MixedPascalToPascalCase(reference.TypeName)
					relationships = append(relationships, ResourceRelationship{
						RustSourcePropertyName: rustPropertyName,
						RustReturnType:         rustPropertyType,
						RustGenericType:        "",
						TargetUnion:            nil,
					})
				}
			}
		}

		rustModel = append(rustModel, ResourceType{
			CfnResourceName:     awsResourceName.AsCfn(),
			RustResourceName:    awsResourceName.AsRust(),
			GraphQlResourceName: awsResourceName.AsGraphQl(),
			Properties:          properties,
			Relationships:       relationships,
		})
	}

	return rustModel, errors
}

// Helper to deref possibly nil pointers in template
func Deref[T any](pointer *T) T {
	if pointer == nil {
		var zero T
		return zero
	}
	return *pointer
}
