import json
import logging
from datetime import datetime
from pathlib import Path
from typing import Any

from transform import to_types_map, find_serialized_objects
from translate import translate_gql_common, translate_resource


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
            # keep original CFN name
            types_map["typeName"] = data["typeName"]

            write_json_file(Path(resource_file_name.map_file), types_map)

            # map_file(data, file.name)
        except Exception as e:
            logging.error("File %s failed. %s", file, repr(e))

    logging.info("=======[Maps have been created]=======")

    maps = Path("out/map/").glob("*.json")
    serialized_objects: dict[str, list[str]] = {}
    for file in maps:
        try:
            types_map = read_json_file(file)
            resource_name = types_map["resourceTypeName"]

            so_entry = find_serialized_objects(types_map, resource_name)
            if so_entry:
                serialized_objects[types_map["typeName"]] = so_entry

        except Exception as e:
            logging.error("File %s failed. %s", file, repr(e))
    write_json_file(Path("out/serialized-objects.json"), serialized_objects)

    logging.info("=======[Untyped objects file created]=======")

    maps = Path("out/map/").glob("*.json")
    resource_struct_names = []
    for file in maps:
        try:
            types_map = read_json_file(file)
            rust, gql, resource_struct_name = translate_resource(types_map)
            resource_struct_names.append(resource_struct_name)
            resource_file_name = ResourceFileName(file.name)
            with open(Path(resource_file_name.rust_file), "w", encoding="utf-8") as rf:
                rf.write(rust)
            with open(Path(resource_file_name.gql_file), "w", encoding="utf-8") as gf:
                gf.write(gql)
        except Exception as e:
            logging.error("File %s failed. %s", file, repr(e))

    try:
        common_gql = translate_gql_common(resource_struct_names)
        with open("out/graphql/common.gql", "w", encoding="utf-8") as gf:
            gf.write(common_gql)
    except Exception as e:
        logging.error("File common.gql failed. %s", repr(e))


def create_types_map(data: dict[str, Any], resource_name: str):
    definitions: dict[str, Any] = data.get("definitions", {})
    properties: dict[str, Any] = data["properties"]
    types_map = to_types_map(resource_name, definitions, properties)
    return types_map


def resource_type_to_name(type_: str) -> str:
    _, service, resource = type_.split("::")
    return f"{service}{resource}"


def read_json_file(file_name: Path) -> dict[str, Any]:
    with open(file_name, "r", encoding="utf-8") as file:
        data: dict[str, Any] = json.load(file)
    return data


def write_json_file(file_name: Path, data: dict[str, Any]):
    with open(file_name, "w", encoding="utf-8") as file:
        json.dump(data, file, indent=4)


if __name__ == "__main__":
    main()
