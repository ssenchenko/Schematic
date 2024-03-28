import json
import logging
from datetime import datetime
from pathlib import Path
from typing import Any

from translate import translate_all


LOG = logging.getLogger(__name__)

IA_SCOPE_ONLY = True  # set to True if necessary to test only for IA resources

IA_RESOURCES = frozenset(
    {
        "AWS::EC2::Instance",
        "AWS::EC2::VolumeAttachment",
        "AWS::EC2::Volume",
        "AWS::EC2::VPC",
        "AWS::EC2::SecurityGroup",
        "AWS::IAM::InstanceProfile",
        "AWS::IAM::Role",
        "AWS::IAM::Policy",
        "AWS::CloudWatch::Alarm",
        "AWS::S3::Bucket",
    }
)

# some property names seems to be wrong in relationship schema file
ALL_SCHEMA_COMBINED_OVERRIDES = {"AWS::EC2::Instance": {"VolumeAttachments": "Volumes"}}


def main():
    cfn_schema_files = Path("cfn/").glob("*.json")
    cfn_schema = {}
    filter = IA_RESOURCES if IA_SCOPE_ONLY else frozenset()
    for file in cfn_schema_files:
        try:
            data = read_json_file(file)
            resource_name = data["typeName"]
            if filter and resource_name not in filter:
                continue
            cfn_schema[resource_name] = data
        except Exception as e:
            LOG.error("File %s failed. %s", file.name, repr(e))

    try:
        all_schema_combined_file = Path("relationship/all-schema-combined.json")
        all_schema_combined = read_json_file(all_schema_combined_file)
        all_schema_combined = apply_all_schema_combined_overrides(
            all_schema_combined, ALL_SCHEMA_COMBINED_OVERRIDES
        )
        output = translate_all(all_schema_combined, cfn_schema, filter=filter)
        with open(Path("out/model.rs"), "w", encoding="utf-8") as gf:
            gf.write(output)
    except Exception as e:
        LOG.error("Translation failed. %s", repr(e))


def apply_all_schema_combined_overrides(
    all_schema_combined: dict[str, Any], overrides: dict[str, dict[str, str]]
):
    for override_resource in overrides:
        for item in all_schema_combined[override_resource]["relationships"]:
            for old_name, new_name in overrides[override_resource].items():
                if old_name in item:
                    item[new_name] = item.pop(old_name)

    return all_schema_combined


def read_json_file(file_name: Path) -> dict[str, Any]:
    with open(file_name, "r", encoding="utf-8") as file:
        data: dict[str, Any] = json.load(file)
    return data


if __name__ == "__main__":
    now = datetime.now()
    now_str = now.isoformat()
    logging.basicConfig(
        filename=f"logs/{now_str}.log", encoding="utf-8", level=logging.INFO
    )
    main()
