from enum import Enum
from typing import Any, TypeAlias, Literal

RefKey = Literal["$ref"]
Ref: TypeAlias = dict[RefKey, str]

Property: TypeAlias = dict[str, Any]
Properties: TypeAlias = dict[str, Property]
Definitions: TypeAlias = dict[str, dict[str, Any]]

REF_KEY: RefKey = "$ref"
TYPE_KEY = "type"
ENUM_KEY = "enum"
PROPS_KEY = "properties"
ITEMS_KEY = "items"
SERIALIZABLE_OBJECT = "object-string"

# // Map which better fits the needs of types generation.
# // Output of the transformer.
# // Looks like:
# //
# //	{
# //		Resource: {
# //			F1: Struct1
# //			F2: string  # Xxxx
# //			F3: Array/Struct1 # A1
# //			F4: Array/integer
# //			F5: Enum1
# //		},
# //		Xxxx: string # type def
# //		Struct1: {
# //			F1: String
# //			F2: integer
# //			F3: Array/Enum2
# //		},
# //		Enum1: [E1, E2],
# //		Enum2: [E3, E4]
# //		A1: Array/Struct1 # type def
# //	}
TypesMapValue: TypeAlias = str, str | dict[str, str] | list[str]
TypesMap: TypeAlias = dict[TypesMapValue]


class JsonSchemaTypes(Enum):
    """
    An enumeration of JSON schema types.
    """

    STRING = "string"
    NUMBER = "number"
    INTEGER = "integer"
    BOOLEAN = "boolean"
    OBJECT = "object"
    ARRAY = "array"
    NULL = "null"  # not sure how to handle it


def to_types_map(
    resource_name: str, definitions: Definitions, properties: Properties
) -> TypesMap:
    """
    Creates a map of property names to types.
    :param resource_name: The name of the resource.
    :param definitions: The definitions of the resource.
    :param properties: The properties of the resource.
    :return: A map of property names to types.
    """
    definitions = _merge_resource_to_definitions(resource_name, definitions, properties)
    types_map = _init_types_map(resource_name, definitions)
    # now I can populate the structs
    for type_name, type_info in definitions.items():
        # only objects with properties, everything else is in TypesMap
        if (
            type_info.get(TYPE_KEY) == JsonSchemaTypes.OBJECT.value
            and PROPS_KEY in type_info
        ):
            for prop_name, prop_data in type_info[PROPS_KEY].items():
                discovered_type = traverse_properties(prop_name, prop_data, types_map)
                if _is_enum(discovered_type, types_map):
                    if isinstance(discovered_type, list):
                        discovered_type = _post_enum(
                            discovered_type, prop_name, types_map
                        )
                    discovered_type = _enum_name(discovered_type)
                types_map[type_name][prop_name] = discovered_type

    return types_map


def traverse_properties(
    key: str, data: dict[str, Any], types_map: TypesMap
) -> str | dict[str, Any]:
    if TYPE_KEY in data:
        if data[TYPE_KEY] in [
            JsonSchemaTypes.INTEGER.value,
            JsonSchemaTypes.NUMBER.value,
            JsonSchemaTypes.BOOLEAN.value,
            JsonSchemaTypes.NULL.value,
        ]:
            return data[TYPE_KEY]
        elif data[TYPE_KEY] == JsonSchemaTypes.STRING.value:
            if ENUM_KEY in data:
                return data[ENUM_KEY]
            return JsonSchemaTypes.STRING.value
        elif data[TYPE_KEY] == JsonSchemaTypes.ARRAY.value:
            # traverse
            discovered_type: str | list[str] = traverse_properties(
                ITEMS_KEY, data[ITEMS_KEY], types_map
            )
            discovered_type = _process_discovered_array_type(
                discovered_type, key, types_map
            )
            return discovered_type
        elif data[TYPE_KEY] == JsonSchemaTypes.OBJECT.value:
            if PROPS_KEY not in data:
                return SERIALIZABLE_OBJECT
            else:
                # need to create another struct
                types_map[key] = {}
                for prop_name, prop_data in data[PROPS_KEY].items():
                    discovered_type = traverse_properties(
                        prop_name, prop_data, types_map
                    )
                    if isinstance(discovered_type, list):
                        discovered_type = _post_enum(
                            discovered_type, prop_name, types_map
                        )
                    types_map[key][prop_name] = discovered_type
                return key
        elif isinstance(data[TYPE_KEY], list):
            return SERIALIZABLE_OBJECT
        else:
            raise ValueError(f"Invalid type: {data[TYPE_KEY]} in property {key}")
    elif REF_KEY in data:
        ref_type_name = _get_ref_type(data[REF_KEY])
        type_info = _get_from_types_map(ref_type_name, types_map)
        if type_info is None:
            raise ValueError(
                f"All references must have been discovered, but this one is not in TypesMap: {data[REF_KEY]}"
            )

        # only objects and enums have their type names
        if isinstance(type_info, dict) or isinstance(type_info, list):
            return ref_type_name
        # return underlying type
        return type_info
    elif (
        "oneOf" in data or "anyOf" in data
    ):  # only 2 files with it, not worth supporting type merge
        return SERIALIZABLE_OBJECT
    else:
        raise ValueError(f"Invalid data: {data}, no '{TYPE_KEY}' nor '{REF_KEY}'")


def traverse_arrays(
    key: str, data: dict[str, Any], types_map: TypesMap
) -> str | dict[str, list[str]]:
    if TYPE_KEY in data:
        if data[TYPE_KEY] in [
            JsonSchemaTypes.INTEGER.value,
            JsonSchemaTypes.NUMBER.value,
            JsonSchemaTypes.BOOLEAN.value,
            JsonSchemaTypes.NULL.value,
        ]:
            return data[TYPE_KEY]
        elif data[TYPE_KEY] == JsonSchemaTypes.STRING.value:
            if ENUM_KEY in data:
                return data[ENUM_KEY]
            return JsonSchemaTypes.STRING.value
        elif data[TYPE_KEY] == JsonSchemaTypes.ARRAY.value:
            # traverse
            discovered_type: str | list[str] = traverse_arrays(
                ITEMS_KEY, data[ITEMS_KEY], types_map
            )
            discovered_type = _process_discovered_array_type(
                discovered_type, key, types_map
            )
            return discovered_type
        elif data[TYPE_KEY] == JsonSchemaTypes.OBJECT.value:
            if PROPS_KEY not in data:
                return SERIALIZABLE_OBJECT
            else:
                return key
        elif isinstance(data[TYPE_KEY], list):
            return SERIALIZABLE_OBJECT
        else:
            raise ValueError(f"Invalid type: {data[TYPE_KEY]} in definition {key}")
    elif REF_KEY in data:
        ref_type_name = _get_ref_type(data[REF_KEY])
        type_info = _get_from_types_map(ref_type_name, types_map)
        if type_info is not None:
            # only objects and enums have their type names
            if isinstance(type_info, dict) or isinstance(type_info, list):
                return ref_type_name
            # return underlying type
            return type_info
        else:
            discovered_type: str | list[str] = traverse_arrays(
                ITEMS_KEY, data[ITEMS_KEY], types_map
            )
            discovered_type = _process_discovered_array_type(
                discovered_type, key, types_map
            )
            return discovered_type
    elif (
        "oneOf" in data or "anyOf" in data
    ):  # only 37 files with it, not worth supporting now
        return SERIALIZABLE_OBJECT
    else:
        raise ValueError(f"Invalid data: {data}, no '{TYPE_KEY}' nor '{REF_KEY}'")


def _merge_resource_to_definitions(
    resource_name: str, definitions: Definitions, properties: Properties
) -> Definitions:
    """
    Merges the resource name and properties to the definitions since resource properties are the same as definitions
    :param resource_name: The name of the resource.
    :param definitions: The definitions of the resource.
    :param properties: The properties of the resource.
    :return: The merged definitions.
    """
    definitions[resource_name] = {
        "type": JsonSchemaTypes.OBJECT.value,
        "properties": properties,
    }
    return definitions


def _init_types_map(resource_name: str, definitions: Definitions) -> TypesMap:
    """
    Initializes the types map with the resource name and the definitions, so we can dereference $ref.
    :param resource_name: The name of the resource.
    :param definitions: The definitions of the resource.
    :return: The initialized types map.
    """
    types_map: TypesMap = {"typeName": resource_name}
    array_defs: dict[str, Any] = {}

    for type_name, type_info in definitions.items():
        if TYPE_KEY in type_info:
            if type_info[TYPE_KEY] in [
                JsonSchemaTypes.INTEGER.value,
                JsonSchemaTypes.NUMBER.value,
                JsonSchemaTypes.BOOLEAN.value,
                JsonSchemaTypes.NULL.value,
            ]:
                types_map[type_name] = type_info[TYPE_KEY]
            elif type_info[TYPE_KEY] == JsonSchemaTypes.STRING.value:
                if ENUM_KEY in type_info:
                    _post_enum(type_info[ENUM_KEY], type_name, types_map)
                else:
                    types_map[type_name] = type_info[TYPE_KEY]
            elif type_info[TYPE_KEY] == JsonSchemaTypes.ARRAY.value:
                # to describe array I need to check if it's element type is scalar or not
                # for that reason do arrays only after all other types,
                # when the $ref to the element types can be dereferenced
                array_defs[type_name] = type_info
            elif type_info[TYPE_KEY] == JsonSchemaTypes.OBJECT.value:
                if PROPS_KEY not in type_info:
                    types_map[type_name] = SERIALIZABLE_OBJECT
                else:
                    types_map[type_name] = {}
            elif isinstance(type_info[TYPE_KEY], list):
                types_map[type_name] = SERIALIZABLE_OBJECT
            else:
                raise ValueError(f"Invalid type: {type_info} in definition {type_name}")
        elif "oneOf" in type_info or "anyOf" in type_info:
            types_map[type_name] = SERIALIZABLE_OBJECT
        else:
            raise ValueError(
                f"Invalid data, no key or One/AnyOf: {type_info} in definition {type_name}"
            )

    for type_name, type_info in array_defs.items():
        discovered_type: str | dict[str, list[str]] = traverse_arrays(
            type_name, type_info, types_map
        )
        types_map[type_name] = discovered_type

    return types_map


def _process_discovered_array_type(
    discovered_type: str | list[str], key: str, types_map: TypesMap
):
    if isinstance(discovered_type, list):  # enum
        discovered_type = _post_enum(discovered_type, key, types_map)
    return f"array/{discovered_type}"


def _post_enum(enum: list[str], key: str, types_map: TypesMap):
    discovered_type = _enum_name(key)
    test = _get_from_types_map(
        discovered_type, types_map
    )  # TODO: remove it, should work without it
    if not test:
        types_map[discovered_type] = enum
    return discovered_type


def _enum_name(type_or_prop_name: str) -> str:
    """
    Since enum can be found in Array items, it should have its own name
    :param type_or_prop_name:
    :return:
    """
    return (
        f"{type_or_prop_name}Enum"
        if not type_or_prop_name.endswith("Enum")
        else type_or_prop_name
    )


def _get_ref_type(ref: str) -> str:
    """
    Extracts the type name from "#/definitions/<TYPE_NAME>".
    :param ref: The $ref.
    :return: The type name.
    """
    return ref.split("/")[-1]


def _get_from_types_map(type_name: str, types_map: TypesMap) -> TypesMapValue:
    if type_name.endswith("Enum"):
        return types_map.get(type_name)
    type_info = types_map.get(type_name)
    if type_info is not None:
        return type_info
    return types_map.get(_enum_name(type_name))


def _is_enum(type_: str | list[str], types_map: TypesMap):
    if isinstance(type_, list):
        return True
    if type_.endswith("Enum"):
        return True
    type_info = types_map.get(type_)
    if type_info is not None:
        return False
    type_ = _enum_name(type_)
    type_info = types_map.get(type_)
    if type_info is not None:
        return True
    return False


def find_serialized_objects(types_map: TypesMap, resource_name: str) -> list[str]:
    serialized_objects: list[str] = []
    visited: set[str] = set()

    visited.add(resource_name)
    stack = []
    for prop_name, prop_type in types_map[resource_name].items():
        stack.append((prop_type, prop_name))

    while stack:
        prop_type_name, path = stack.pop()
        if prop_type_name == SERIALIZABLE_OBJECT:
            serialized_objects.append(path)

        elif prop_type_name in [
            JsonSchemaTypes.STRING.value,
            JsonSchemaTypes.INTEGER.value,
            JsonSchemaTypes.NUMBER.value,
            JsonSchemaTypes.BOOLEAN.value,
            JsonSchemaTypes.NULL.value,
        ]:
            continue  # nothing found

        elif prop_type_name.startswith("array/"):
            path = f"{path}[*]"
            # unwrap one array
            prop_type_name = "/".join(prop_type_name.split("/")[1:])
            stack.append((prop_type_name, path))

        elif prop_type_name in visited:
            continue  # already checked that type

        elif isinstance(types_map[prop_type_name], list):
            continue  # nothing to check in an Enum

        else:
            visited.add(prop_type_name)
            for pname, ptype in types_map[prop_type_name].items():
                new_path = f"{path}.{pname}"
                stack.append((ptype, new_path))

    return serialized_objects
