package translator

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

type RustModel struct {
	Resources []ResourceType
	Connections []ResourceConnections
}

type ResourceType struct {
	CfnResourceName     string
	RustResourceName    string
	GraphQlResourceName string
	UseComplex          string
	Properties          []ResourceProperty
}

type ResourceProperty struct {
	RustPropertyName string
	RustPropertyType string
}

type ResourceConnections struct {
	Source ResourceType
	Relationships []ResourceRelationship
}

type ResourceRelationship struct {
	RustSourcePropertyName string
	RustReturnType string
	RustGenericType string
	TargetUnion *ResourceUnion
}

type ResourceUnion struct {
	RustUnionName string
	Resources []ResourceType
}

func (allRelationships *AllRelationships) ApplyOverrides(overrides map[string]map[string]string) {
	for resourceName, overrides := range overrides {
		for oldName, newName := range overrides {
			for _, relationship := range(*allRelationships)[resourceName].Relationships {
				relationship[newName] = relationship[oldName]
				delete(relationship, oldName)
			}
		}
	}
}

// helper to deref possibly nil pointers in template
func Deref[T any](pointer *T) T {
	if pointer == nil {
		var zero T
		return zero
	}
	return *pointer
}
