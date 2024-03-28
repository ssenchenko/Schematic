import abc
import logging
from dataclasses import dataclass, field
from enum import Enum
from typing import Iterator, TypedDict, Any

from traverse import is_array
from util import resource_type_to_name, pascal_to_snake

LOG = logging.getLogger(__name__)


class LanguageType(TypedDict):
    rust: str
    gql: str


class Types(Enum):
    STRING = "string"
    NUMBER = "number"
    INTEGER = "integer"
    BOOLEAN = "boolean"
    OBJECT = "object"
    ARRAY = "array"
    NULL = "null"
    SERIALIZABLE_OBJECT = "object-string"


TYPES: dict[Types, LanguageType] = {
    Types.STRING: {
        "rust": "String",
        "gql": "String",
    },
    Types.SERIALIZABLE_OBJECT: {"rust": "Value", "gql": "String"},
    Types.BOOLEAN: {
        "rust": "bool",
        "gql": "Boolean",
    },
    Types.INTEGER: {
        "rust": "i32",
        "gql": "Int",
    },
    Types.NUMBER: {"rust": "f64", "gql": "Float"},
    Types.ARRAY: {
        "rust": "Vec<{}>",
        "gql": "[{}]",
    },
}


class Translatable(abc.ABC):
    """Interface for any type and property to translate to Rust and GraphQL"""

    @property
    @abc.abstractmethod
    def rust(self) -> str:
        """Rust type"""

    @property
    def gql(self) -> str:
        """GraphQL type"""
        return ""


class PropertyName(Translatable):

    def __init__(self, name: str, keep_case: bool = False):
        self.original_name = name
        self.keep_case = keep_case

    @property
    def rust(self) -> str:
        if self.keep_case:
            return self.original_name
        return pascal_to_snake(self.original_name)

    def __repr__(self):
        return self.original_name


class StructName(Translatable):
    # an option if we want to use it
    TYPE_SUFFIXES_TO_REMOVE = ("Config", "Configuration", "Options")

    def __init__(self, name: str, prefix: str = ""):
        self.original_name = name
        # no prefix for resource
        self._prefix = resource_type_to_name(
            prefix, preserve_case=False, keep_partition=True
        )

    def _prefixed_struct_name(self, name: str):
        return f"{self._prefix}{name}"

    @property
    def rust(self) -> str:
        name = resource_type_to_name(
            self.original_name, preserve_case=False, keep_partition=True
        )
        return self._prefixed_struct_name(name)

    @property
    def gql(self) -> str:
        name = resource_type_to_name(
            self.original_name, separator="_", preserve_case=False, keep_partition=True
        )
        return self._prefixed_struct_name(name)

    def __repr__(self) -> str:
        return self._prefixed_struct_name(self.original_name)


class ScalarType(Translatable):

    def __init__(self, name: Types):
        self._name = name

    @property
    def rust(self) -> str:
        return TYPES[self._name]["rust"]

    def __repr__(self):
        return self._name


class Property(Translatable):

    def __init__(
        self,
        type_: Translatable,
        name: PropertyName,
        required: bool = False,
    ):
        self.name = name
        self.required = required
        self.property_type = type_

    def __repr__(self):
        return str({self.name: self.property_type})

    @property
    def rust(self) -> str:
        if self.required:
            return f"pub {self.name.rust}: {self.property_type.rust},"
        return f"{self.name.rust}: Option<{self.property_type.rust}>,"


class Struct(Translatable):
    def __init__(
        self,
        struct_name: StructName,
        is_resource: bool = False,
        in_relationships: bool = False,
    ):
        self.name = struct_name
        self.properties: list[Property] = []
        self.is_resource = is_resource
        self.has_relationships = in_relationships

        if self.is_resource:
            self.properties.append(
                Property(ScalarType(Types.STRING), PropertyName("Id"), required=True)
            )
            if self.has_relationships:
                self.properties.append(
                    Property(
                        ScalarType(Types.STRING),
                        PropertyName("AllProperties"),
                        required=True,
                    )
                )

    def append(self, prop: Property):
        self.properties.append(prop)

    @property
    def rust(self) -> str:
        use_complex = ", complex" if self.has_relationships else ""

        struct_str = f"""
#[derive(SimpleObject, Serialize)]
#[graphql(name = "{self.name.gql}", rename_fields = "PascalCase"{use_complex})]
pub struct {self.name.rust} {{
"""

        for prop in self.properties:
            struct_str += tab(line(prop.rust))
        struct_str += line("}")

        # relationship struct declaration in translate_relationship

        # type_name goes as implementation now
        struct_str += f"""
#[ComplexObject(rename_fields = "PascalCase")]
impl {self.name.rust} {{
    pub async fn type_name(&self) -> String {{
        "{self.name.original_name}".to_string()
    }}
}}
"""
        return struct_str

    def __repr__(self):
        return str({self.name: self.properties})


class RustEnumInterface(Translatable):
    """
    Rust enum which is translated to GraphQL interface and its implementations.

    Only for the purpose of generating GraphQL from Rust types.
    Looks like this:

    #[derive(Interface, Serialize)]
    #[graphql(
        name = "Resource",
        rename_fields = "PascalCase",
        field(name = "id", ty = "String"),
        field(name = "type_name", ty = "String"),
        field(name = "all_properties", ty = "String"),
    )]
    pub enum Resource {
        AwsLambdaFunction(AwsLambdaFunction),
        AwsIamRole(AwsIamRole),
        Node(Node),
    }

    In GraphQL, it will generate interface

    interface Resource {
        Id: String!
        TypeName: String!
        AllProperties: String!
    }

    And adds 'implements Resource' to all types in the enum.
    """

    def __init__(self, name: str = "Resource"):
        self.name = name
        self.resources: list[StructName] = []

    def append(self, resource: StructName):
        self.resources.append(resource)

    @property
    def rust(self) -> str:
        enum_str = f"""
#[derive(Interface, Serialize)]
#[graphql(
    name = \"{self.name}\",
    rename_fields = "PascalCase",
    field(name = "id", ty = "String"),
    field(name = "type_name", ty = "String"),
    field(name = "all_properties", ty = "String"),
)]
pub enum Resource {{
"""
        for resource in self.resources:
            enum_str += line(tab(f"{resource.rust}({resource.rust}),"))
        enum_str += line(tab("Node(Node)"))
        enum_str += line("}")

        return enum_str

    @property
    def gql(self) -> str:
        return ""


class RustEnumUnion(Translatable):
    """
    Rust enum which is translated to GraphQL union and its implementations.

    Only for the purpose of generating GraphQL from Rust types.
    Looks like this:

    #[derive(Union, Serialize)]
    pub enum MyUnion {
        AwsLambdaFunction(AwsLambdaFunction),
        AwsIamRole(AwsIamRole),
        Node(Node),
    }

    In GraphQL, it will generate union

    union Resource = AwsLambdaFunction | AwsIamRole | Node
    """

    def __init__(self, name: StructName, rust_types: list[str]):
        self._property_name = name
        self.name = f"{self._property_name.rust}Connections"
        self.types: list[str] = rust_types

    @property
    def rust(self) -> str:
        enum_str = f"""
#[derive(Union, Serialize)]
pub enum {self.name} {{
"""
        for type_ in self.types:
            enum_str += line(tab(f"{type_}({type_}),"))
        enum_str += line("}")
        return enum_str


def tab(string: str, number: int = 1, size: int = 4) -> str:
    tab_str = " " * size
    many_tabs = tab_str * number
    return f"{many_tabs}{string}"


def line(string: str) -> str:
    return f"{string}\n" if not string.endswith("\n") else string


@dataclass
class Relationship:
    source_property_name: PropertyName
    rust_to_cfn_type_map: dict[str, list[str]] = field(default_factory=dict)
    is_source_array: bool = False


def translate_relationship(
    source_type_cfn_notation: str, relationships: list[Relationship]
) -> str:
    source_struct_name = StructName(source_type_cfn_notation)
    template = """
#[derive(Serialize)]
pub struct {source_struct_name_rust}Relationships<'a> {{
    properties: &'a String,
}}

impl {source_struct_name_rust} {{

    pub async fn relationships(&self) -> {source_struct_name_rust}Relationships {{
        {source_struct_name_rust}Relationships {{ properties: &self.all_properties }}
    }}

}}
{unions}
#[Object(rename_fields = "PascalCase")]
impl {source_struct_name_rust}Relationships<'_> {{
{relationship_functions}
}}
"""

    vector_fn_template = """
    pub async fn {source_property_name_rust}(&self) -> Vec<{return_type}> {{
        vec![]
    }}
"""

    single_fn_template = """
    pub async fn {source_property_name_rust}(&self) -> {return_type} {{
        {return_expression}
    }}
"""

    return_expression_template = """{instance_type}{{
            {instance_body}
        }}"""

    wrapped_return_expression_template = """model::{return_type}::{instance_type}({instance_type}{{
            {instance_body}
        }})"""

    instance_body_template = """id: "".to_string(),
            all_properties: "".to_string(),"""

    node_body_template = """id: "".to_string(),
            type_name: "{type_name_cfn_format}".to_string(),
            all_properties: "".to_string(),"""

    def has_union_return(rel: Relationship):
        return len(rel.rust_to_cfn_type_map.keys()) > 1

    def hydrate_instance_body(
        instance_type_rust_notation: str,
        rust_to_cfn_type_map: dict[str, list[str]],
    ):
        # if there are > 1 cfn types, use any of mapped types as a placeholder
        type_cfn_notation = rust_to_cfn_type_map[instance_type_rust_notation][0]
        instance_body = (
            node_body_template.format(
                # use any of Node types as a placeholder
                type_name_cfn_format=type_cfn_notation
            )
            if instance_type_rust_notation == "Node"
            else instance_body_template
        )
        return instance_body

    unions: list[RustEnumUnion] = []
    relationship_functions: list[str] = []
    for relationship in relationships:
        # if not a union, there is a single type, otherwise, it will be a random type from union
        first_rust_type = list(relationship.rust_to_cfn_type_map.keys())[0]

        return_type = first_rust_type

        if has_union_return(relationship):  # need to create a union
            union = RustEnumUnion(
                StructName(
                    relationship.source_property_name.original_name,
                    prefix=source_type_cfn_notation,
                ),
                sorted(list(relationship.rust_to_cfn_type_map.keys())),
            )
            unions.append(union)
            return_type = union.name

        if relationship.is_source_array:
            fn_template = vector_fn_template.format(
                source_property_name_rust=relationship.source_property_name.rust,
                return_type=return_type,
            )
        else:
            # it's either a single item or we use a random type from union as a placeholder
            instance_type_rust_notation = first_rust_type
            instance_body = hydrate_instance_body(
                instance_type_rust_notation, relationship.rust_to_cfn_type_map
            )

            if has_union_return(relationship):
                return_expression = wrapped_return_expression_template.format(
                    return_type=return_type,
                    instance_type=instance_type_rust_notation,
                    instance_body=instance_body,
                )
            else:
                return_expression = return_expression_template.format(
                    instance_type=instance_type_rust_notation,
                    instance_body=instance_body,
                )

            fn_template = single_fn_template.format(
                source_property_name_rust=relationship.source_property_name.rust,
                return_type=return_type,
                return_expression=return_expression,
            )

        relationship_functions.append(fn_template)

    unions_output = "".join(union.rust for union in unions)
    relationship_functions_output = "".join(fn for fn in relationship_functions)
    output = template.format(
        source_struct_name_rust=source_struct_name.rust,
        unions=unions_output,
        relationship_functions=relationship_functions_output,
    )

    return output


def translate_types_with_relationships(
    schema_all_combined: dict[str, Any],
    filter: frozenset[str] = frozenset(),
) -> list[Struct]:
    structs = []
    for resource_name, resource_data in schema_all_combined.items():
        if filter and resource_name not in filter:
            continue

        if has_relationships(resource_name, schema_all_combined):
            structs.append(
                Struct(
                    StructName(resource_name),
                    is_resource=True,
                    in_relationships=True,
                )
            )
    # it's easier to find smth in sorted list
    structs.sort(key=lambda x: x.name.original_name)
    return structs


def translate_all_relationships(
    all_schema_combined: dict[str, Any],
    cfn_schema: dict[str, Any],
    filter: frozenset[str] = frozenset(),
) -> list[str]:
    relationships: list[str] = []
    for source_type_cfn_notation in all_schema_combined:
        if filter and source_type_cfn_notation not in filter:
            continue

        try:
            rel = translate_resource_relationships(
                source_type_cfn_notation, all_schema_combined, cfn_schema, filter=filter
            )
            if rel:
                relationships.append(rel)
        except Exception as e:
            LOG.error("Relationship translation failed. %s", repr(e))

    return relationships


def translate_resource_relationships(
    source_type_cfn_notation: str,
    all_schema_combined: dict[str, Any],
    cfn_schema: dict[str, Any],
    filter: frozenset[str] = frozenset(),
):
    relationships: list[Relationship] = []
    for relation in all_schema_combined[source_type_cfn_notation]["relationships"]:
        for property_name_with_slashes, references in relation.items():
            property_name_with_slashes: str
            references: list[dict[str, Any]]
            if not references:
                LOG.error(
                    "Relationship field doesn't have any references %s %s",
                    source_type_cfn_notation,
                    property_name_with_slashes,
                )
                continue

            target_type_cfn_notation = list(
                {
                    ref["typeName"]
                    for ref in references
                    if not filter or filter and ref["typeName"] in filter
                }
            )
            if not target_type_cfn_notation:
                continue

            is_source_array = is_array(
                property_name_with_slashes, source_type_cfn_notation, cfn_schema
            )

            type_name_map: dict[str, list[str]] = {}
            for type_ in target_type_cfn_notation:
                rust_name = (
                    StructName(type_).rust
                    if has_relationships(type_, all_schema_combined)
                    else "Node"
                )
                if rust_name not in type_name_map:
                    type_name_map[rust_name] = []
                type_name_map[rust_name].append(type_)

            property_name = PropertyName(property_name_with_slashes.replace("/", ""))

            relationships.append(
                Relationship(
                    source_property_name=property_name,
                    rust_to_cfn_type_map=type_name_map,
                    is_source_array=is_source_array,
                )
            )

    return (
        translate_relationship(source_type_cfn_notation, relationships)
        if relationships
        else ""
    )


def translate_all(
    all_schema_combined: dict[str, Any],
    cfn_schema: dict[str, Any],
    filter: frozenset[str] = frozenset(),
) -> str:
    output = """// =======================================================
// This file is generated.  Do not edit manually!
// =======================================================
use async_graphql::{Context, Enum, Error, Interface, OutputType, SimpleObject, Result, Object, ComplexObject, Union};
use serde::Serialize;
"""
    queries = """
pub struct Query;

#[Object]
impl Query {
    pub async fn resource(&self, id: String, type_name: String) -> Resource {
        Resource::AwsLambdaFunction(AwsLambdaFunction{
            id: id,
            type_name: type_name,
            all_properties: "{}".to_string(),
        })
    }

    pub async fn resources(&self, type_name: String) -> Vec<Resource> {
        let resources: Vec<Resource> = vec![];
        resources
    }

    pub async fn topology(&self, id: String, type_name: String) -> Topology {
        Topology {
            nodes: vec![],
            edges: vec![],
        }
    }
}
"""
    interface = RustEnumInterface()
    number_of_resources_translated = 0
    structs = translate_types_with_relationships(all_schema_combined, filter=filter)
    for struct in structs:
        output += struct.rust
        number_of_resources_translated += 1
        interface.append(struct.name)

    LOG.info("Translated %d resources", number_of_resources_translated)

    output += interface.rust

    output += """
#[derive(SimpleObject, Serialize, Clone)]
#[graphql(name = "Node", rename_fields = "PascalCase")]
pub struct Node {
    pub id: String,
    pub type_name: String,
    pub all_properties: String,
}

#[derive(SimpleObject, Serialize)]
#[graphql(name = "Edge", rename_fields = "PascalCase")]
pub struct Edge {
    source: String,
    target: String,
    relation: Relation,
}

#[derive(Enum, Copy, Clone, Eq, PartialEq, Serialize)]
pub enum Relation {
    IsRelatedTo,
}

#[derive(SimpleObject, Serialize)]
#[graphql(name = "Topology", rename_fields = "PascalCase")]
pub struct Topology {
    nodes: Vec<Node>,
    edges: Vec<Edge>,
}
"""
    output += """
// =========== Relationships ===========
"""

    number_of_relationships_translated = 0
    for relationship in translate_all_relationships(
        all_schema_combined, cfn_schema, filter=filter
    ):
        output += relationship
        number_of_relationships_translated += 1
        output += line("")

    LOG.info("Translated %d relationships", number_of_relationships_translated)

    # TODO: Add queries when they are ready
    # output += queries

    return output


def has_relationships(
    type_name_cfn_format: str, all_schema_combined: dict[str, Any]
) -> bool:
    return len(all_schema_combined[type_name_cfn_format]["relationships"]) > 0
