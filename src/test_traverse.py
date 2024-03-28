import json
from pathlib import Path
from typing import cast

import pytest

from traverse import follow_path, TypeBits


@pytest.fixture
def resource_types():
    return [
        "AWS::AutoScaling::AutoScalingGroup",
    ]


@pytest.fixture
def load_from_cfn_schema():
    def _load(file_path):
        path = Path(file_path)
        with open(path, "r") as file:
            return json.load(file)

    return _load


@pytest.fixture
def resources_schema(resource_types, load_from_cfn_schema):
    cfn_schema_test_dataset = {}
    for resource_type in resource_types:
        resource_type = cast(str, resource_type).replace("::", "-").lower()
        data = load_from_cfn_schema(f"cfn/{resource_type}.json")
        resource_name = data["typeName"]
        cfn_schema_test_dataset[resource_name] = data

    return cfn_schema_test_dataset


def test_traverse(resources_schema):
    set_of_types = follow_path(
        "NotificationConfigurations/TopicARN",
        resources_schema["AWS::AutoScaling::AutoScalingGroup"],
    )
    assert set_of_types == {TypeBits.STRING | TypeBits.OBJECT | TypeBits.ARRAY}
