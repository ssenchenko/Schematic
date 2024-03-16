import abc
import re
from enum import Enum
from typing import TypedDict

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
    Types.SERIALIZABLE_OBJECT: {"rust": "String", "gql": "String"},
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


class Name:
    # an option if we want to use it
    TYPE_SUFFIXES_TO_REMOVE = ("Config", "Configuration", "Options")

    def __init__(self, name: str, prefix: str):
        self._name = name
        self._prefix = prefix
        self._prefixed_struct_name = f"{self._prefix}{self._name}"

    def __str__(self) -> str:
        return self._prefixed_struct_name

    def __repr__(self) -> str:
        return self._prefixed_struct_name

    @property
    def prefix(self) -> str:
        return self._prefix

    @property
    def gql_field(self) -> str:
        return self._name

    @property
    def rust_field(self) -> str:
        return pascal_to_snake(self._name)

    @property
    def rust_type(self) -> str:
        return self._prefixed_struct_name

    @property
    def gql_type(self) -> str:
        return self._prefixed_struct_name

    @property
    def gql_filter(self) -> str:
        return f"{self._prefixed_struct_name}Filter"


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


class PropertyName(Translatable, GqlFilterable):
    def __init__(self, name: str):
        self._name = name

    @property
    def rust(self) -> str:
        return self._name

    @property
    def gql(self) -> str:
        return self._name

    @property
    def gql_filter(self) -> str:
        return self._name

    def __repr__(self):
        return self._name
        # return f"PropertyName({self._name})"


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


class Array(Translatable, GqlFilterable):
    def __init__(self, item_type: ScalarType):
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
    def __init__(self, type_: TypeName, name: str):
        self.name = Name(name, "")
        self.property_type = type_

    def __repr__(self):
        return str({self.name: self.property_type})

    @property
    def rust(self) -> str:
        return f"{self.name.rust_field}: {self.property_type.rust}"

    @property
    def gql(self) -> str:
        return f"{self.name.gql_field}: {self.property_type.gql}"

    @property
    def gql_filter(self) -> str:
        return f"{self.name.gql_field}: {self.property_type.gql_filter}"
