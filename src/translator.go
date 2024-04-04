package schematic

type Reference struct {
	TypeName  string `json:"typeName"`
	Attribute string `json:"attribute"`
}

type Relationship = map[string][]Reference

type Resourse struct {
	PrimaryKeys   []string       `json:"primaryIdentifier"`
	Relationships []Relationship `json:"relationships"`
}

type SchemaCombined = map[string]Resourse

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

type ResourceUnion struct {
	RustUnionName string
	Resources []ResourceType
}

type ResourceConnections struct {
	Source ResourceType
	Relationships []ResourceRelationship
}

type ResourceRelationship struct {
	SourceProperty ResourceProperty
	RustReturnType string
	RustGenericType string
	TargetUnion *ResourceUnion
}
