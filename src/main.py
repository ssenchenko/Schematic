import abc
import json
import logging
import re
from datetime import datetime
from pathlib import Path
from typing import Any, Optional, TypeAlias, TypedDict

from transform import to_types_map, TypesMap, find_serialized_objects

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
    "number": {"rust": "f64", "gql": "Float"},
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


class ResourceFileName:
    DIR = "out"
    RUST_DIR = "model"
    GQL_DIR = "graphql"
    MAP_DIR = "map"
    MAP_EXTENSION = "json"
    RUST_EXTENSION = "rs"
    GQL_EXTENSION = "gql"

    def __init__(self, file_name: str):
        self.file_name = file_name
        self.name_base = file_name.split(".")[0]

    @property
    def resource_type_name(self):
        parts = self.name_base.split("-")[1:]
        parts = [x.capitalize() for x in parts]
        return "".join(parts)

    @property
    def rust_file(self):
        return f"{self.DIR}/{self.RUST_DIR}/{self.name_base}.{self.RUST_EXTENSION}"

    @property
    def gql_file(self):
        return f"{self.DIR}/{self.GQL_DIR}/{self.name_base}.{self.GQL_EXTENSION}"

    @property
    def map_file(self):
        return f"{self.DIR}/{self.MAP_DIR}/{self.name_base}.{self.MAP_EXTENSION}"


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


class Float(TypeName):
    @property
    def rust(self) -> str:
        return TYPES["number"]["rust"]

    @property
    def gql(self) -> str:
        return TYPES["number"]["gql"]

    @property
    def gql_filter(self) -> str:
        return TYPES["number"]["gql"]

    def __repr__(self):
        return "number"


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
    now = datetime.now()
    now_str = now.isoformat()
    logging.basicConfig(
        filename=f"logs/{now_str}.log", encoding="utf-8", level=logging.INFO
    )
    all_schema = "cfn/"
    files = Path(all_schema).glob("*.json")

    for file in files:
        try:
            data = read_json_file(file)
            resource_file_name = ResourceFileName(file.name)
            resource_name = resource_type_to_name(data["typeName"])
            types_map = create_types_map(data, resource_name)

            with open(resource_file_name.map_file, "w") as map_file:
                json.dump(types_map, map_file, indent=4)

            # map_file(data, file.name)
        except Exception as e:
            logging.error(f"File {file} failed. {repr(e)}")

    logging.info("=======[Maps have been created]=======")

    all_maps = "out/map/"
    maps = Path(all_maps).glob("*.json")
    serialized_objects: dict[str, list[str]] = {}
    for file in maps:
        try:
            types_map = read_json_file(file)
            resource_name = types_map["typeName"]
            del types_map["typeName"]

            so_entry = find_serialized_objects(types_map, resource_name)
            if so_entry:
                serialized_objects[resource_name] = so_entry
        except Exception as e:
            logging.error(f"File {file} failed. {repr(e)}")
    with open(Path("out/serialized-objects.json"), "w") as so:
        json.dump(serialized_objects, so, indent=4)


def create_types_map(data: dict[str, Any], resource_name: str):
    definitions: dict[str, Any] = data.get("definitions", {})
    properties: dict[str, Any] = data["properties"]
    types_map = to_types_map(resource_name, definitions, properties)
    return types_map


def map_file(data: dict[str, Any], file_name: str):
    resource_file_name = ResourceFileName(file_name)
    definitions: dict[str, Any] = data.get("definitions", {})
    properties: dict[str, Any] = data["properties"]
    declarations: TypeDeclarations = {}
    resource_name = resource_file_name.resource_type_name

    visited: set[str] = set()
    for name, props in definitions.items():
        logging.debug(f"Definition: {name}")
        if name in visited:
            continue
        dfs(
            declarations,
            props,
            name,
            resource_name,
            definitions=definitions,
            visited=visited,
        )

    struct = Struct(declarations, resource_name)  # resource type struct
    for name, props in properties.items():
        dfs(declarations, props, name, resource_name, struct)

    with open(resource_file_name.rust_file, "w") as rust_file:
        with open(resource_file_name.gql_file, "w") as gql_file:
            for _, type_ in declarations.items():
                logging.debug(f"Printing {type_} to file")
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
    definitions: Optional[dict[str, Any]] = None,
    visited: Optional[set[str]] = None,
):
    stack: list[tuple[dict[str, Any], str, str, Optional[Struct]]] = [
        (data, name, prefix, struct)
    ]

    while stack:
        data, name, prefix, struct = stack.pop()
        if visited is not None:
            visited.add(name)

        if "$ref" in data:
            reference: str = data["$ref"]
            type_name = reference.split("/")[-1]
            property_type = declarations.get(f"{prefix}{type_name}")
            if property_type is None:
                # we are still processing definitions and the reference to another definitions
                if type_name not in definitions:
                    raise Exception(
                        f"Type: {type_name} in {name} not found in definitions"
                    )
                if definitions[type_name]["type"] == "object":
                    # stop processing current node till type is ready, send them back to stack
                    stack.append((data, name, prefix, struct))
                    # visit referred type first and start new struct
                    stack.append((definitions[type_name], type_name, prefix, None))
                else:  # keep working with the same struct
                    stack.append((definitions[type_name], type_name, prefix, struct))
            else:
                struct.append(
                    Property(property_type, name)
                )  # $ref is only in the properties

        elif data["type"] == "object":
            if "properties" in data:
                if struct is None:
                    struct = Struct(declarations, name, prefix)
                prefix = f"{prefix}{name}"
                for prop_name, prop_data in data["properties"].items():
                    if prop_name not in visited:
                        stack.append((prop_data, prop_name, prefix, struct))
            else:
                # map to string bc in GraphQL arbitrary named keys are not possible,
                # so it will be serialized json
                struct.append(Property(String(), name))

        elif data["type"] == "array":
            if "enum" in data["items"]:
                property_type = Enum_(declarations, data["items"]["enum"], name, prefix)
            elif "$ref" in data["items"]:
                stack.append((data["items"], name, prefix, struct))
                continue
            else:
                property_type = map_scalar_type(data["items"]["type"])
            struct.append(Property(Array(property_type), name))

        elif data["type"] == "string":
            if "enum" in data:
                property_type = Enum_(declarations, data["enum"], name, prefix)
            else:
                property_type = String()
            struct.append(Property(property_type, name))

        elif data["type"] == [
            "object",
            "string",
        ]:  # maybe a typo, from aws-iam-managedpolicy.json
            struct.append(Property(String(), name))  # it'll be serialized json

        elif data["type"] == "integer":
            struct.append(Property(Integer(), name))

        elif data["type"] == "number":
            struct.append(Property(Float(), name))

        elif data["type"] == "boolean":
            struct.append(Property(Boolean(), name))

        else:
            raise Exception(f" Unknown type: {data['type']} in {name}")


def map_scalar_type(type_name: str) -> Optional[TypeName]:
    if type_name == "string":
        return String()
    if type_name == "integer":
        return Integer()
    if type_name == "boolean":
        return Boolean()
    if type_name == "number":
        return Float()
    raise Exception(f"[ERROR]: Unknown type: {type_name}")


def tab(string: str, number: int = 1, size: int = 4) -> str:
    tab_str = " " * size
    many_tabs = tab_str * number
    return f"{many_tabs}{string}"


def line(string: str) -> str:
    return f"{string}\n"


def pascal_to_snake(pascal_string: str) -> str:
    snake_string = PATTERN.sub("_", pascal_string).lower()
    return snake_string


def resource_type_to_name(type_: str) -> str:
    _, service, resource = type_.split("::")
    return f"{service}{resource}"


def read_json_file(file_name: Path) -> dict[str, Any]:
    with open(file_name, "r") as file:
        data: dict[str, Any] = json.load(file)
    return data


if __name__ == "__main__":
    main()
