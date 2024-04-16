package translator

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/iancoleman/strcase"

	"ssenchenko/schematic/cfn"
)

const (
	TemplateDir string = "./templates"
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
	TargetUnion            ResourceUnion
}

type ResourceUnion struct {
	RustUnionName     string
	RustResourceNames map[string]bool
}

func (rel *AllRelationships) ApplyOverrides(overrides map[string]map[string]string) {
	for resourceName, overrides := range overrides {
		for oldName, newName := range overrides {
			for _, relationship := range (*rel)[resourceName].Relationships {
				if _, ok := relationship[oldName]; ok {
					relationship[newName] = relationship[oldName]
					delete(relationship, oldName)
					break
				}
			}
		}
	}
}

func (rel *AllRelationships) HasRelationships(resourceName string, filter map[string]bool) bool {
	if resource, ok := (*rel)[resourceName]; ok {
		for _, relationship := range resource.Relationships {
			for _, references := range relationship {
				for _, ref := range references {
					if isInFilter(ref.TypeName, filter) {
						return true
					}
				}
			}
		}
	}
	return false
}

func Translate(
	allRelationships AllRelationships,
	cfnSchemaCombined map[string]cfn.Dict,
	filter map[string]bool,
)  (RustModel, []error) {
	return translateAllResources(allRelationships, cfnSchemaCombined, filter, translateResource)
}

func translateAllResources(
	allRelationships AllRelationships,
	cfnSchemaCombined map[string]cfn.Dict,
	filter map[string]bool,
	translatorFunc func( // to mock for testing
		string,
		AllRelationships,
		map[string]cfn.Dict,
		map[string]bool,
	) (ResourceType, []error),
) (RustModel, []error) {
	errors := make([]error, 0)
	rustModel := make(RustModel, 0, len(allRelationships))

	// to make generated file "stable" so the diff is not generated because of the
	var resources []string
	if len(filter) > 0 {
		resources = make([]string, 0, len(filter))
		for resourceName := range filter {
			resources = append(resources, resourceName)
		}
	} else {
		resources = make([]string, 0, len(allRelationships))
		for resourceName := range allRelationships {
			resources = append(resources, resourceName)
		}
	}
	sort.Slice(resources, func(i, j int) bool { return resources[i] < resources[j] })

	for _, resourceName := range resources {
		res, resErrors := translatorFunc(
			resourceName, allRelationships, cfnSchemaCombined, filter)
		if len(resErrors) > 0 {
			errors = append(errors, resErrors...)
		}
		if res.RustResourceName != "" {
			rustModel = append(rustModel, res)
		}
	}

	return rustModel, errors
}

func translateResource(
	resourceName string,
	allRelationships AllRelationships,
	cfnSchemaCombined map[string]cfn.Dict,
	filter map[string]bool,
) (ResourceType, []error) {
	resource := allRelationships[resourceName]
	errors := make([]error, 0)

	if len(resource.Relationships) == 0 {
		// no separate type for this resource, it's just Node
		return ResourceType{}, nil
	}

	awsResourceName, err := NewAwsResourceName(resourceName)
	if err != nil {
		errors = append(errors, err)
		return ResourceType{}, errors
	}
	rustResourceName := awsResourceName.AsRust()
	graphQlResourceName := awsResourceName.AsGraphQl()

	properties := []ResourceProperty{
		{
			RustPropertyName: "identifier",
			RustPropertyType: "String",
		},
		{
			RustPropertyName: "all_properties",
			RustPropertyType: "String",
		},
	}

	var relationships []ResourceRelationship

	for _, relationship := range resource.Relationships {
		for propertyName, references := range relationship {
			withoutSlashes := strings.ReplaceAll(propertyName, "/", "")
			rustPropertyName := strcase.ToSnake(withoutSlashes)
			resourceUnion, unionErrors := createResourceUnion(rustResourceName, withoutSlashes, references, allRelationships, filter)
			if len(unionErrors) > 0 {
				errors = append(errors, unionErrors...)
				continue
			}

			var rustUnderlyingType string
			switch len(resourceUnion.RustResourceNames) {
			case 0:
				continue // all connected types are not in filter
			case 1:
				for t := range resourceUnion.RustResourceNames {
					rustUnderlyingType = t
					// break -- no need for it
				}
				resourceUnion = ResourceUnion{}
			default:
				rustUnderlyingType = resourceUnion.RustUnionName
			}

			isArray, err := cfn.IsArray(propertyName, resourceName, cfnSchemaCombined)
			if err != nil {
				errors = append(errors, err)
				continue
			}

			var rustReturnType, rustGenericType string
			if isArray {
				rustReturnType = fmt.Sprintf("Vec<%s>", rustUnderlyingType)
				rustGenericType = rustReturnType
			} else {
				rustReturnType = fmt.Sprintf("Option<%s>", rustUnderlyingType)
				rustGenericType = rustUnderlyingType
			}

			relationships = append(relationships, ResourceRelationship{
				RustSourcePropertyName: rustPropertyName,
				RustReturnType:         rustReturnType,
				RustGenericType:        rustGenericType,
				TargetUnion:            resourceUnion,
			})
		}
	}

	return ResourceType{
		CfnResourceName:     resourceName,
		RustResourceName:    rustResourceName,
		GraphQlResourceName: graphQlResourceName,
		Properties:          properties,
		Relationships:       relationships,
	}, errors
}

func HydrateTemplates(templateData RustModel) (bytes.Buffer, error) {
	return hydrateTemplate(
		templateData,
		"all.go.tmpl",
		filepath.Join("translator", TemplateDir),
		nil,
		"interface_enum.go.tmpl",
		"resource_struct.go.tmpl",
		"relationship.go.tmpl",
		"union_enum.go.tmpl",
		"aws_resource_impl.go.tmpl",
	)
}

func getRustUnderlyingType(allRelationships AllRelationships, reference Reference, filter map[string]bool) (string, error) {
	var rustUnderlyingType string
	if allRelationships.HasRelationships(reference.TypeName, filter) {
		referenceType, err := NewAwsResourceName(reference.TypeName)
		if err != nil {
			return "", err
		}
		rustUnderlyingType = referenceType.AsRust()
	} else {
		rustUnderlyingType = "Node"
	}
	return rustUnderlyingType, nil
}

func createResourceUnion(
	rustResourceName string,
	propertyPathWithoutSlashes string,
	references []Reference,
	allRelationships AllRelationships,
	filter map[string]bool,
) (ResourceUnion, []error) {
	errors := make([]error, 0)

	rustUnionName := fmt.Sprintf(
		"%sConnections_%s", rustResourceName, strcase.ToCamel(propertyPathWithoutSlashes))

	uniqueTypes := make(map[string]bool)
	for _, reference := range references {
		if !isInFilter(reference.TypeName, filter) {
			continue
		}

		rustUnderlyingType, err := getRustUnderlyingType(allRelationships, reference, filter)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		uniqueTypes[rustUnderlyingType] = true
	}

	return ResourceUnion{
		RustUnionName:     rustUnionName,
		RustResourceNames: uniqueTypes,
	}, errors
}

func isInFilter(resourceName string, filter map[string]bool) bool {
	if len(filter) == 0 {
		return true
	}
	_, ok := filter[resourceName]
	return ok
}
