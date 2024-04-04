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
