import abc
import json
import re
from typing import Any, Optional, TypeAlias, TypedDict

DEBUG = True


def debug(message: str):
    if DEBUG:
        print(message)


PATTERN = re.compile(r"(?<!^)(?=[A-Z])")


class LanguageType(TypedDict):
    rust: str
    gql: str | tuple[str, str]


TYPES: dict[str, LanguageType] = {
    "string": {
        "rust": "String",
        "gql": "String",
    },
    "boolean": {
        "rust": "bool",
        "gql": "Boolean",
    },
    "integer": {
        "rust": "i32",
        "gql": "Int",
    },
    "array": {
        "rust": "Vec<{}>",
        "gql": "[{}]",
    },
}

KEYWORDS: dict[str, LanguageType] = {
    "enum": {
        "rust": "enum",
        "gql": "enum",
    },
    "struct": {
        "rust": "struct",
        "gql": ("type", "input"),
    },
}


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
        return self._name

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

    @property
    def rust_enum(self) -> str:
        return f"{self.rust_type}Enum"

    @property
    def gql_enum(self) -> str:
        return f"{self.gql_type}Enum"


class Translatable(abc.ABC):
    """Interface for any type and property to translate to Rust and GraphQL"""

    @abc.abstractmethod
    def to_rust(self) -> str:
        """Rust type"""

    @abc.abstractmethod
    def to_gql(self) -> str:
        """GraphQL type"""


class GqlFilterable(abc.ABC):
    @abc.abstractmethod
    def to_gql_filter(self) -> str:
        """GraphQL input"""


class TypeName(abc.ABC):
    @property
    @abc.abstractmethod
    def rust(self) -> str:
        """Rust type"""

    @property
    @abc.abstractmethod
    def gql(self) -> str:
        """GraphQL type"""

    @property
    @abc.abstractmethod
    def gql_filter(self) -> str:
        """GraphQL input"""


class DeclarableType(Translatable, TypeName, abc.ABC):
    """Parent type for Structs (types, inputs) and Enums"""

    def __init__(
        self, declarations: dict[str, "DeclarableType"], name: str, prefix: str
    ):
        self._name = Name(name, prefix)
        # register itself
        declarations[str(self._name)] = self


TypeDeclarations: TypeAlias = dict[str, DeclarableType]


class Boolean(TypeName):
    @property
    def rust(self) -> str:
        return TYPES["boolean"]["rust"]

    @property
    def gql(self) -> str:
        return TYPES["boolean"]["gql"]

    @property
    def gql_filter(self) -> str:
        return TYPES["boolean"]["gql"]

    def __repr__(self):
        return "boolean"


class Integer(TypeName):
    @property
    def rust(self) -> str:
        return TYPES["integer"]["rust"]

    @property
    def gql(self) -> str:
        return TYPES["integer"]["gql"]

    @property
    def gql_filter(self) -> str:
        return TYPES["integer"]["gql"]

    def __repr__(self):
        return "integer"


class String(TypeName):
    @property
    def rust(self) -> str:
        return TYPES["string"]["rust"]

    @property
    def gql(self) -> str:
        return TYPES["string"]["gql"]

    @property
    def gql_filter(self) -> str:
        return TYPES["string"]["gql"]

    def __repr__(self):
        return "string"


class Array(TypeName):
    def __init__(self, item_type: TypeName):
        self.item_type = item_type

    @property
    def rust(self) -> str:
        return f"{TYPES['array']['rust'].format(self.item_type.rust)},"

    @property
    def gql(self) -> str:
        return f"{TYPES['array']['gql'].format(self.item_type.gql)}"

    @property
    def gql_filter(self) -> str:
        return f"{TYPES['array']['gql'].format(self.item_type.gql_filter)}"

    def __repr__(self):
        return f"array<{self.item_type}>"


class Property(Translatable, GqlFilterable):
    def __init__(self, property_type: TypeName, property_name: str):
        self.name = Name(property_name, "")
        self.property_type = property_type

    def __repr__(self):
        return str({self.name: self.property_type})

    def to_rust(self) -> str:
        return f"{self.name.rust_field}: {self.property_type.rust}"

    def to_gql(self) -> str:
        return f"{self.name.gql_field}: {self.property_type.gql}"

    def to_gql_filter(self) -> str:
        return f"{self.name.gql_field}: {self.property_type.gql_filter}"


class Enum_(DeclarableType):  # not GqlFilterable
    def __init__(
        self,
        declarations: TypeDeclarations,
        values: list[str],
        property_name: str,
        prefix: str = "",
    ):
        super().__init__(declarations, property_name, prefix)
        self.values = values

    @property
    def rust(self) -> str:
        return self._name.rust_enum

    @property
    def gql(self) -> str:
        return self._name.gql_enum

    @property
    def gql_filter(self) -> str:
        return self._name.gql_enum  # doesn't change in filter

    def to_rust(self) -> str:
        enum_str = line(f"{KEYWORDS['enum']['rust']} {self.rust} {{")
        for value in self.values:
            # string enums to make it compatible with GraphQL
            # since they are serialized to JSON as strings
            enum_str += tab(line(f"{value}(String),"))
        enum_str += line("}")
        return enum_str

    def to_gql(self) -> str:
        enum_str = line(f"{KEYWORDS['enum']['gql']} {self.gql} {{")
        for value in self.values:
            enum_str += tab(line(f"{value}"))
        enum_str += line("}")
        return enum_str

    def __repr__(self):
        return str({self._name: self.values})


class Struct(DeclarableType, GqlFilterable):
    def __init__(self, declarations: TypeDeclarations, name: str, prefix: str = ""):
        super().__init__(declarations, name, prefix)
        self.properties: list[Property] = []

    def __repr__(self):
        return str({self._name: self.properties})

    def append(self, prop: Property):
        self.properties.append(prop)

    @property
    def rust(self) -> str:
        return self._name.rust_type

    @property
    def gql(self) -> str:
        return self._name.gql_type

    @property
    def gql_filter(self) -> str:
        return self._name.gql_filter

    def to_rust(self) -> str:
        struct_str = line("#[derive(Serialize, Deserialize, Debug)]")
        struct_str += line('#[serde(rename_all = "PascalCase")]')
        struct_str += line(f"{KEYWORDS['struct']['rust']} {self.rust} {{")
        for prop in self.properties:
            struct_str += tab(line(f"{prop.to_rust()},"))
        struct_str += line("}")
        return struct_str

    def to_gql(self) -> str:
        type_keyword = KEYWORDS["struct"]["gql"][0]
        struct_str = line(f"{type_keyword} {self.gql} {{")
        for prop in self.properties:
            struct_str += tab(line(f"{prop.to_gql()}"))
        struct_str += line("}")
        return struct_str

    def to_gql_filter(self) -> str:
        input_keyword = KEYWORDS["struct"]["gql"][1]
        struct_str = line(f"{input_keyword} {self.gql_filter} {{")
        for prop in self.properties:
            struct_str += tab(line(f"{prop.to_gql_filter()}"))
        struct_str += line("}")
        return struct_str


def main():
    data = read_json_file("cfn/lambda-function.json")
    map_file(data)


def map_file(data: dict[str, Any]):
    definitions: dict[str, Any] = data["definitions"]
    properties: dict[str, Any] = data["properties"]
    declarations: TypeDeclarations = {}
    resource_name = "LambdaFunction"
    for name, props in definitions.items():
        debug(f"Definition: {name}")
        dfs(declarations, props, name, resource_name)
    struct = Struct(declarations, resource_name)  # resource type struct
    for name, props in properties.items():
        dfs(declarations, props, name, resource_name, struct)
    with open("out/lambda-function.rs", "w") as rust_file:
        with open("out/lambda-function.gql", "w") as gql_file:
            for _, type_ in declarations.items():
                debug(f"About to print {type_}")
                rust_file.write(type_.to_rust())
                rust_file.write("\n")
                gql_file.write(type_.to_gql())
                if isinstance(type_, GqlFilterable):
                    gql_file.write(type_.to_gql_filter())
                gql_file.write("\n")


def dfs(
    declarations: TypeDeclarations,
    data: dict[str, Any],
    name: str,
    prefix,
    struct: Optional[Struct] = None,
):
    stack: list[tuple[dict[str, Any], str, str, Optional[Struct]]] = [
        (data, name, prefix, struct)
    ]

    while stack:
        data, name, prefix, struct = stack.pop()
        if "$ref" in data:
            property_type, type_name = follow_reference(declarations, data, prefix)
            if not property_type:
                print(
                    f"[ERROR]: Type: {type_name} in {name} should have been mapped but it's not"
                )
            struct.append(
                Property(property_type, name)
            )  # $ref is only in the properties

        elif data["type"] == "object":
            if "properties" in data:
                if struct is None:
                    struct = Struct(declarations, name, prefix)
                prefix = f"{prefix}{name}"
                for prop_name, prop_data in data["properties"].items():
                    stack.append((prop_data, prop_name, prefix, struct))
            else:
                # map to string bc in GraphQL arbitrary named keys are not possible,
                # so it will be serialized json
                struct.append(Property(String(), name))

        elif data["type"] == "array":
            if "enum" in data["items"]:
                property_type = Enum_(declarations, data["items"]["enum"], name, prefix)
            elif "$ref" in data["items"]:
                property_type, type_name = follow_reference(
                    declarations, data["items"], prefix
                )
                if not property_type:
                    print(
                        f"[ERROR]: Type: {type_name} in {name} should have been mapped but it's not"
                    )
            else:
                property_type = map_scalar_type(data["items"]["type"])
            struct.append(Property(Array(property_type), name))

        elif data["type"] == "string":
            if "enum" in data:
                property_type = Enum_(declarations, data["enum"], name, prefix)
            else:
                property_type = String()
            struct.append(Property(property_type, name))

        elif data["type"] == "integer":
            struct.append(Property(Integer(), name))

        elif data["type"] == "boolean":
            struct.append(Property(Boolean(), name))

        else:
            print(f"[ERROR]: Unknown type: {data['type']} in {name}")


def follow_reference(
    declarations: TypeDeclarations, property_info: dict[str, Any], prefix: str
) -> tuple[DeclarableType, str]:
    reference: str = property_info["$ref"]
    type_name = reference.split("/")[-1]
    type_ = declarations[f"{prefix}{type_name}"]
    return type_, type_name


def map_scalar_type(type_name: str) -> Optional[TypeName]:
    if type_name == "string":
        return String()
    if type_name == "integer":
        return Integer()
    if type_name == "boolean":
        return Boolean()
    print(f"[ERROR]: Unknown type: {type_name}")


def tab(string: str, number: int = 1, size: int = 4) -> str:
    tab_str = " " * size
    many_tabs = tab_str * number
    return f"{many_tabs}{string}"


def line(string: str) -> str:
    return f"{string}\n"


def pascal_to_snake(pascal_string: str) -> str:
    snake_string = PATTERN.sub("_", pascal_string).lower()
    return snake_string


def read_json_file(file_name: str) -> dict[str, Any]:
    with open(file_name, "r") as file:
        data: dict[str, Any] = json.load(file)
    return data


if __name__ == "__main__":
    main()
