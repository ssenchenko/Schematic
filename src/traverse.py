import logging
from enum import Enum, Flag, auto
from typing import Any

LOG = logging.getLogger(__name__)

REF_KEY = "$ref"
TYPE_KEY = "type"
ENUM_KEY = "enum"
PROPS_KEY = "properties"
ITEMS_KEY = "items"
ONE_OF_KEY = "oneOf"
ANY_OF_KEY = "anyOf"


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


class TypeBits(Flag):
    STRING = auto()
    NUMBER = auto()
    INTEGER = auto()
    BOOLEAN = auto()
    OBJECT = auto()
    ARRAY = auto()
    NULL = auto()


def get_ref_type(ref: str) -> str:
    """
    Extracts the type name from "#/definitions/<TYPE_NAME>".
    :param ref: The $ref.
    :return: The type name.
    """
    return ref.split("/")[-1]


def json_type_to_bits(type_: str) -> TypeBits:
    type_map = {
        "string": TypeBits.STRING,
        "number": TypeBits.NUMBER,
        "integer": TypeBits.INTEGER,
        "boolean": TypeBits.BOOLEAN,
        "object": TypeBits.OBJECT,
        "array": TypeBits.ARRAY,
        "null": TypeBits.NULL,
    }
    return type_map[type_]


def is_array(
    path_to_property: str, type_name_cfn_format: str, cfn_schema: dict[str, Any]
) -> bool:
    """
    Find type of property.

    :param path_to_property: single property like "ApplicationId"
            or chain of properties joined with "/" ie "BackupPlan/BackupPlanRule/TargetBackupVault"
    :param type_name_cfn_format: like Alexa::ASK::Skill
    :param cfn_schema: full cfn schema, a map with types in cfn format
    :return:
    """
    types = follow_path(path_to_property, cfn_schema[type_name_cfn_format])
    if not types:
        raise Exception(
            f"No type found for {path_to_property}, in {type_name_cfn_format}"
        )

    array_test = [TypeBits.ARRAY in type_ for type_ in types]
    # if they are not the same, I don't know what should I do with that
    if not all(element == array_test[0] for element in array_test):
        raise Exception(
            f"WTF am I supposed to do with it??? {types} of {path_to_property} in {type_name_cfn_format}"
        )
    return array_test[0]


def follow_path(path: str, whereabouts: dict[str, Any]) -> set[TypeBits]:
    steps = path.split("/")
    where = whereabouts
    type_ = TypeBits.OBJECT
    step_number = 0
    # dfs bc we want to explore each branch of any/one fully first
    # and being able to return to the fork node if it wasn't the right branch
    journey: list[tuple[int, dict[str, Any], TypeBits]] = [(step_number, where, type_)]
    types_seen: set[TypeBits] = set()
    try:
        while journey:
            step_number, where, type_ = journey.pop()
            if step_number == len(steps):
                types_seen.add(type_)
                continue

            places = one_step_at_a_time(steps[step_number], where, whereabouts)
            if not places:
                LOG.debug(f"Step {steps[step_number]} not found in {where}")
                continue  # probably not our branch
            for place in places:
                # keeping track of types to be able to check if array happened on the "correct" path
                # sometimes "type" is an array like ["object", "string"] or ["array", "list"]
                # which are probably errors, but we need to handle those without breaking
                # as far as we concerned, we need to know if any of the types along the path
                # can be an array, so let's bit-OR all the types from the array
                if isinstance(place[TYPE_KEY], list):
                    for place_type in place[TYPE_KEY]:
                        type_ |= json_type_to_bits(place_type)
                else:
                    type_ |= json_type_to_bits(place[TYPE_KEY])
                journey.append((step_number + 1, place, type_))
    except Exception as e:
        raise Exception(f"Error in {path}  {repr(e)}") from e
    return types_seen


def one_step_at_a_time(
    step: str, where: dict[str, Any], whereabouts: dict[str, Any]
) -> list[dict[str, Any]]:
    if PROPS_KEY in where:  # type == object with properties
        if step in where[PROPS_KEY]:
            # step has been found, resolve now the type
            if TYPE_KEY in where[PROPS_KEY][step]:
                return [where[PROPS_KEY][step]]
            if REF_KEY in where[PROPS_KEY][step]:
                return resolve_ref(where[PROPS_KEY][step][REF_KEY], whereabouts)
            if (
                ONE_OF_KEY in where[PROPS_KEY][step]
                or ANY_OF_KEY in where[PROPS_KEY][step]
            ):
                key = ONE_OF_KEY if ONE_OF_KEY in where[PROPS_KEY][step] else ANY_OF_KEY
                return resolve_any(where[PROPS_KEY][step][key], whereabouts)
            raise Exception(f"Oh no, another uncharted terrain {where}")
        else:
            return []  # treat it carefully, it makes sense because of anyOf/oneOf
    # in cases below step might be found deeper, in array or in $ref
    elif TYPE_KEY in where:
        # we get here if previous step was an array and next one is supposed to be in an item of object type
        if (
            where[TYPE_KEY] == JsonSchemaTypes.ARRAY.value
            # I found 1 case when "type" == ["array", "string"] which looks like a mistake to me,
            # and it should probably be "oneOf/anyOf" or just "array" instead but
            # for the time being, it's faster to add this condition
            # than reach out to the owner and ask for a change
            or JsonSchemaTypes.ARRAY.value in where[TYPE_KEY]
        ):
            # keep looking for the same step
            return one_step_at_a_time(step, where[ITEMS_KEY], whereabouts)
        return []  # no way to find step, type is neither typed object nor array
    # we get here if the previous step led us to array which 'item' is '$ref' (see "elif TYPE_KEY in where" branch)
    elif REF_KEY in where:
        # first resolve ref
        places = resolve_ref(where[REF_KEY], whereabouts)
        # then keep looking for the step in the resolved type(s)
        results = []
        for place in places:
            results.extend(one_step_at_a_time(step, place, whereabouts))
        return results
    # we get here if the previous step led us to array which 'item' is oneOf/anyOf (see "elif TYPE_KEY in where" branch)
    elif ONE_OF_KEY in where or ANY_OF_KEY in where:
        # resolve to an actual type(s)
        key = ONE_OF_KEY if ONE_OF_KEY in where[PROPS_KEY][step] else ANY_OF_KEY
        places = resolve_any(where[key], whereabouts)
        # look for our step
        results = []
        for place in places:
            results.extend(one_step_at_a_time(step, place, whereabouts))
        return results
    else:
        raise Exception(f"Hmm... Lost in {where} on step {step}")


def resolve_ref(ref: str, whereabouts: dict[str, Any]) -> list[dict[str, Any]]:
    ref_type_name = get_ref_type(ref)
    next_destination = whereabouts["definitions"][ref_type_name]
    if TYPE_KEY in next_destination:
        return [next_destination]
    if REF_KEY in next_destination:
        return resolve_ref(next_destination[REF_KEY], whereabouts)
    if ONE_OF_KEY in next_destination or ANY_OF_KEY in next_destination:
        return resolve_any(next_destination, whereabouts)
    raise Exception(f"Oh no, another uncharted terrain {next_destination}")


def resolve_any(
    where: list[dict[str, Any]], whereabouts: dict[str, Any]
) -> list[dict[str, Any]]:
    results = []
    for branch in where:
        if TYPE_KEY in branch:
            results.append(branch)
        elif REF_KEY in branch:
            results.extend(resolve_ref(branch[REF_KEY], whereabouts))
        elif ONE_OF_KEY in branch or ANY_OF_KEY in branch:
            key = ONE_OF_KEY if ONE_OF_KEY in branch else ANY_OF_KEY
            results.extend(resolve_any(branch[key], whereabouts))
        else:
            raise Exception(f"Oh no, another uncharted terrain {branch}")
    return results
