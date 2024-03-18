import abc
import re
from enum import Enum
from typing import Iterator, TypedDict

from transform import TypesMap

TYPE_NAME_PATTERN = re.compile(r"(?<!^)(?=[A-Z])")


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


def pascal_to_snake(pascal_string: str) -> str:
    snake_string = TYPE_NAME_PATTERN.sub("_", pascal_string).lower()
    return snake_string


class Translatable(abc.ABC):
    """Interface for any type and property to translate to Rust and GraphQL"""

    @property
    @abc.abstractmethod
    def rust(self) -> str:
        """Rust type"""

    @property
    @abc.abstractmethod
    def gql(self) -> str:
        """GraphQL type"""


class GqlFilterable(abc.ABC):
    @property
    @abc.abstractmethod
    def gql_filter(self) -> str:
        """GraphQL input"""


class TranslatableAndFilterable(Translatable, GqlFilterable, abc.ABC):
    """Uniting the two interfaces."""


class PropertyName(Translatable):
    def __init__(self, name: str):
        self._name = name

    @property
    def rust(self) -> str:
        return pascal_to_snake(self._name)

    @property
    def gql(self) -> str:
        return self._name

    def __repr__(self):
        return self._name


class EnumName(Translatable):
    def __init__(self, name: str, prefix: str):
        self._name = name
        self._prefix = prefix
        self._prefixed_name: str = f"{self._prefix}{self._name}"
        # don't add Enum suffix if it's already there

    @property
    def rust(self) -> str:
        return self._prefixed_name

    @property
    def gql(self) -> str:
        return self._prefixed_name

    def __repr__(self):
        return self._prefixed_name


class StructName(Translatable, GqlFilterable):
    # an option if we want to use it
    TYPE_SUFFIXES_TO_REMOVE = ("Config", "Configuration", "Options")

    def __init__(self, name: str, prefix: str):
        self._name = name
        # no prefix for resource
        self._prefix = "" if prefix == name else prefix
        self._prefixed_struct_name = f"{self._prefix}{self._name}"

    @property
    def rust(self) -> str:
        return self._prefixed_struct_name

    @property
    def gql(self) -> str:
        return self._prefixed_struct_name

    @property
    def gql_filter(self) -> str:
        return f"{self._prefixed_struct_name}Filter"

    @property
    def prefix(self) -> str:
        return self._prefix

    def __repr__(self) -> str:
        return self._prefixed_struct_name


class ScalarType(Translatable, GqlFilterable):

    def __init__(self, name: Types):
        self._name = name

    @property
    def rust(self) -> str:
        return TYPES[self._name]["rust"]

    @property
    def gql(self) -> str:
        return TYPES[self._name]["gql"]

    @property
    def gql_filter(self) -> str:
        return TYPES[self._name]["gql"]

    def __repr__(self):
        return self._name


def is_scalar(type_name: str) -> bool:
    return type_name in [
        Types.STRING.value,
        Types.BOOLEAN.value,
        Types.INTEGER.value,
        Types.NUMBER.value,
        Types.SERIALIZABLE_OBJECT.value,
        Types.NULL.value,
    ]


class Array(Translatable, GqlFilterable):

    def __init__(self, item_type: TranslatableAndFilterable):
        self._name = Types.ARRAY
        self._item_type = item_type

    @property
    def rust(self) -> str:
        return f"{TYPES[self._name]['rust'].format(self._item_type.rust)},"

    @property
    def gql(self) -> str:
        return f"{TYPES[self._name]['gql'].format(self._item_type.gql)}"

    @property
    def gql_filter(self) -> str:
        return f"{TYPES[self._name]['gql'].format(self._item_type.gql_filter)}"

    def __repr__(self):
        return f"array<{repr(self._item_type)}>"


class Property(Translatable, GqlFilterable):

    def __init__(
        self, type_: TranslatableAndFilterable | Translatable, name: PropertyName
    ):
        self.name = name
        self.property_type = type_

    def __repr__(self):
        return str({self.name: self.property_type})

    @property
    def rust(self) -> str:
        return f"pub {self.name.rust}: Option<{self.property_type.rust}>,"

    @property
    def gql(self) -> str:
        return f"{self.name.gql}: {self.property_type.gql}"

    @property
    def gql_filter(self) -> str:
        type_ = (
            self.property_type.gql_filter
            if isinstance(self.property_type, GqlFilterable)
            else self.property_type.gql
        )
        return f"{self.name.gql}: {type_}"


def property_type(
    prop_type_name: str, prefix: str
) -> TranslatableAndFilterable | Translatable:
    if is_scalar(prop_type_name):
        return ScalarType(Types(prop_type_name))
    elif prop_type_name.startswith(Types.ARRAY.value):
        return unwrap_arrays(prop_type_name, prefix)
    elif prop_type_name.endswith("Enum"):
        return EnumName(prop_type_name, prefix)
    else:
        return StructName(prop_type_name, prefix)


def unwrap_arrays(prop_type_name: str, prefix: str) -> Array:
    parts = prop_type_name.split("/")
    if len(parts) < 2 or parts[0] != Types.ARRAY.value:
        raise ValueError(f"Invalid array type name: {prop_type_name}")
    if len(parts) == 2:
        if is_scalar(parts[1]):
            return Array(ScalarType(Types(parts[1])))
        else:
            return Array(StructName(parts[1], prefix))
    nested_array = unwrap_arrays("/".join(parts[1:]), prefix)
    return Array(nested_array)


class Enum_(Translatable):  # pylint: disable=invalid-name
    def __init__(
        self,
        enum_name: EnumName,
        values: list[str],
    ):
        self.name = enum_name
        self.values = values

    @property
    def rust(self) -> str:
        enum_str = line(f"enum {self.name.rust} {{")
        for value in self.values:
            enum_str += line(tab(f"{value},"))
        enum_str += line("}")
        return enum_str

    @property
    def gql(self) -> str:
        enum_str = line(f"enum {self.name.gql} {{")
        for value in self.values:
            enum_str += line(tab(f"{value}"))
        enum_str += line("}")
        return enum_str

    def __repr__(self):
        return str({self.name: self.values})


class Struct(Translatable, GqlFilterable):
    def __init__(
        self,
        struct_name: StructName,
        is_resource: bool = False,
    ):
        self.name = struct_name
        self.properties: list[Property] = []
        self.is_resource = is_resource

    def append(self, prop: Property):
        self.properties.append(prop)

    @property
    def rust(self) -> str:
        struct_str = line("#[derive(Serialize, Deserialize, Debug)]")
        struct_str += line('#[serde(rename_all = "PascalCase")]')
        struct_str += line(f"pub struct {self.name.rust} {{")
        for prop in self.properties:
            struct_str += tab(line(prop.rust))
        struct_str += line("}")
        return struct_str

    @property
    def gql(self) -> str:
        if self.is_resource:
            struct_str = f"""
type {self.name.gql} implements Resource {{
    CcapiId: String!
    CcapiTypeName: String!
"""
        else:
            struct_str = line(f"type {self.name.gql} {{")
        for prop in self.properties:
            struct_str += line(tab(prop.gql))
        struct_str += line("}")
        return struct_str

    @property
    def gql_filter(self) -> str:
        struct_str = line(f"input {self.name.gql_filter} {{")
        for prop in self.properties:
            struct_str += line(tab(prop.gql_filter))
        struct_str += line("}")
        return struct_str

    def __repr__(self):
        return str({self.name: self.properties})


def tab(string: str, number: int = 1, size: int = 4) -> str:
    tab_str = " " * size
    many_tabs = tab_str * number
    return f"{many_tabs}{string}"


def line(string: str) -> str:
    return f"{string}\n" if not string.endswith("\n") else string


def translate(
    types_map: TypesMap,
) -> Iterator[TranslatableAndFilterable | Translatable]:
    resource_name = types_map["resourceTypeName"]
    for type_name, type_ in types_map.items():
        if isinstance(type_, list):
            enum_name = EnumName(type_name, prefix=resource_name)
            yield Enum_(enum_name, type_)
        elif isinstance(type_, dict):
            struct_name = StructName(type_name, prefix=resource_name)
            is_resource = type_name == resource_name
            struct = Struct(struct_name, is_resource=is_resource)
            for prop_name, prop_type_name in type_.items():
                prop_type = property_type(prop_type_name, resource_name)
                prop = Property(prop_type, PropertyName(prop_name))
                struct.append(prop)
            yield struct


def translate_resource(
    types_map: TypesMap,
) -> tuple[str, str, StructName]:
    rust_str = """
use serde::{Serialize, Deserialize};
use serde_json::Value;

"""

    gql_str = ""
    resource_name = None
    for type_ in translate(types_map):
        rust_str += type_.rust
        rust_str += line("")

        gql_str += type_.gql
        if isinstance(type_, Struct):
            gql_str += type_.gql_filter
            if type_.is_resource:
                resource_name = type_.name
        gql_str += line("")

    return rust_str, gql_str, resource_name


def translate_gql_common(resources: list[StructName]) -> str:
    gql_str = """
interface Resource {
    CcapiId: String!
    CcapiTypeName: String!
}

type Adjacent {
  Vertex: Resource!
  Adjacent: [Resource!]
}

union ReturnResourceType =
    """

    for i, resource in enumerate(resources):
        if i == 0:
            gql_str += line(tab(tab(resource.gql, 2)))
            continue
        gql_str += line(tab(f"| {resource.gql}"))

    gql_str += line("")

    gql_str += line("input ResourceFilter {")
    for resource in resources:
        gql_str += line(tab(f"{resource.gql}: {resource.gql_filter}"))
    gql_str += line("}")

    gql_str += """
type Query {
  listResources(type: String!): [ReturnResourceType!] 
  describeResource(id: String!, type: String!): ReturnResourceType
  graphFrom(id: String!, type: String!): [Adjacent]!
}

"""
    return gql_str
