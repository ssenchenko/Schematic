import json
from pathlib import Path

import pytest

from main import ResourceFileName
from transform import to_types_map, find_serialized_objects


@pytest.fixture()
def load_resource():
    def _load_resource(file_path):
        path = Path(file_path)
        resource_file_name = ResourceFileName(path.name)
        with open(path, "r", encoding="utf-8") as file:
            return json.load(file), resource_file_name.resource_type_name

    return _load_resource


def test_to_types_map_not_break(load_resource):
    inputs = [
        "cfn/aws-connect-userhierarchygroup.json",
        "cfn/aws-apigateway-model.json",
    ]
    for input_ in inputs:
        resource, resource_name = load_resource(file_path=input_)
        definitions = resource.get("definitions", {})
        properties = resource["properties"]
        types_map = to_types_map(resource_name, definitions, properties)
        assert types_map is not None


def test_to_types_map_not_double_arrays(load_resource):
    resource, resource_name = load_resource(file_path="cfn/aws-kendra-faq.json")
    definitions = resource.get("definitions", {})
    properties = resource["properties"]
    types_map = to_types_map(resource_name, definitions, properties)
    assert "TagList" in types_map
    assert types_map["TagList"] == "array/Tag"


def test_to_types_set_enum_name_in_array(load_resource):
    resource, resource_name = load_resource(file_path="cfn/aws-kendra-faq.json")
    definitions = resource.get("definitions", {})
    properties = resource["properties"]
    types_map = to_types_map(resource_name, definitions, properties)
    assert "TagList" in types_map
    assert types_map["TagList"] == "array/Tag"


def test_find_serialized_objects():
    inputs = [
        (
            {
                "resourceTypeName": "CEAnomalySubscription",
                "Arn": "string",
                "Subscriber": {
                    "Address": "string",
                    "Status": "StatusEnum",
                    "Type": "TypeEnum",
                },
                "ResourceTag": {"Key": "string", "Value": "string"},
                "CEAnomalySubscription": {
                    "SubscriptionArn": "string",
                    "SubscriptionName": "string",
                    "AccountId": "string",
                    "MonitorArnList": "array/string",
                    "Subscribers": "array/Subscriber",
                    "Threshold": "number",
                    "ThresholdExpression": "string",
                    "Frequency": "FrequencyEnum",
                    "ResourceTags": "array/ResourceTag",
                },
                "StatusEnum": ["CONFIRMED", "DECLINED"],
                "TypeEnum": ["EMAIL", "SNS"],
                "FrequencyEnum": ["DAILY", "IMMEDIATE", "WEEKLY"],
            },
            [],
        ),
        (
            {
                "resourceTypeName": "GluePartition",
                "SchemaReference": {
                    "SchemaId": "SchemaId",
                    "SchemaVersionId": "string",
                    "SchemaVersionNumber": "integer",
                },
                "Order": {"Column": "string", "SortOrder": "integer"},
                "SkewedInfo": {
                    "SkewedColumnValues": "array/string",
                    "SkewedColumnValueLocationMaps": "object-string",
                    "SkewedColumnNames": "array/string",
                },
                "Column": {"Comment": "string", "Type": "string", "Name": "string"},
                "StorageDescriptor": {
                    "StoredAsSubDirectories": "boolean",
                    "Parameters": "object-string",
                    "BucketColumns": "array/string",
                    "NumberOfBuckets": "integer",
                    "OutputFormat": "string",
                    "Columns": "array/Column",
                    "SerdeInfo": "SerdeInfo",
                    "SortColumns": "array/Order",
                    "Compressed": "boolean",
                    "SchemaReference": "SchemaReference",
                    "SkewedInfo": "SkewedInfo",
                    "InputFormat": "string",
                    "Location": "string",
                },
                "SchemaId": {
                    "RegistryName": "string",
                    "SchemaName": "string",
                    "SchemaArn": "string",
                },
                "SerdeInfo": {
                    "Parameters": "object-string",
                    "SerializationLibrary": "string",
                    "Name": "string",
                },
                "PartitionInput": {
                    "StorageDescriptor": "StorageDescriptor",
                    "Values": "array/string",
                    "Parameters": "object-string",
                },
                "GluePartition": {
                    "DatabaseName": "string",
                    "TableName": "string",
                    "Id": "string",
                    "CatalogId": "string",
                    "PartitionInput": "PartitionInput",
                },
            },
            [
                "PartitionInput.Parameters",
                "PartitionInput.StorageDescriptor.Parameters",
                "PartitionInput.StorageDescriptor.SerdeInfo.Parameters",
                "PartitionInput.StorageDescriptor.SkewedInfo.SkewedColumnValueLocationMaps",
            ],
        ),
    ]

    for input_, output in inputs:
        so = find_serialized_objects(input_, input_["resourceTypeName"])

        assert so.sort() == output.sort()
