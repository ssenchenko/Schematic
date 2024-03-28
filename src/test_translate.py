import json
from pathlib import Path
from typing import cast

import pytest

from translate import translate_resource_relationships


@pytest.fixture
def resource_types():
    return [
        "AWS::Connect::UserHierarchyGroup",
        "AWS::EC2::TransitGatewayRouteTablePropagation",
        "AWS::RAM::ResourceShare",
        "AWS::ECR::ReplicationConfiguration",
    ]


@pytest.fixture
def load_relationship_schema():
    def _load():
        path = Path("relationship/all-schema-combined.json")
        with open(path, "r") as file:
            return json.load(file)

    return _load


@pytest.fixture
def load_from_cfn_schema():
    def _load(file_path):
        path = Path(file_path)
        with open(path, "r") as file:
            return json.load(file)

    return _load


@pytest.fixture
def relationships(resource_types, load_relationship_schema):
    all_schema_combined = load_relationship_schema()
    relationships_test_dataset = {}
    related_types = set()
    for resource_type in resource_types:
        relationships_test_dataset[resource_type] = all_schema_combined[resource_type]
        for r in relationships_test_dataset[resource_type]["relationships"]:
            for k, v in r.items():
                for item in v:
                    related_types.add(item["typeName"])
    for resource_type in related_types:
        relationships_test_dataset[resource_type] = all_schema_combined[resource_type]

    return relationships_test_dataset


@pytest.fixture
def resources_schema(resource_types, load_from_cfn_schema):
    cfn_schema_test_dataset = {}
    for resource_type in resource_types:
        resource_type = cast(str, resource_type).replace("::", "-").lower()
        data = load_from_cfn_schema(f"cfn/{resource_type}.json")
        resource_name = data["typeName"]
        cfn_schema_test_dataset[resource_name] = data

    return cfn_schema_test_dataset


def test_translate_relationships_with_single_relationship(
    relationships, resources_schema
):
    translated = translate_resource_relationships(
        "AWS::Connect::UserHierarchyGroup", relationships, resources_schema
    )
    assert (
        translated
        == """
#[derive(Serialize)]
pub struct AwsConnectUserHierarchyGroupRelationships<'a> {
    properties: &'a String,
}

#[Object(rename_fields = "PascalCase")]
impl AwsConnectUserHierarchyGroupRelationships<'_> {

    pub async fn instance_arn(&self) -> Node {
        Node{
            id: "".to_string(),
            type_name: "AWS::Connect::Instance".to_string(),
            all_properties: "".to_string(),
        }
    }

}
"""
    )


def test_translate_relationship_with_nested_property(relationships, resources_schema):
    translated = translate_resource_relationships(
        "AWS::ECR::ReplicationConfiguration", relationships, resources_schema
    )
    assert (
        translated
        == """
#[derive(Serialize)]
pub struct AwsEcrReplicationConfigurationRelationships<'a> {
    properties: &'a String,
}

#[Object(rename_fields = "PascalCase")]
impl AwsEcrReplicationConfigurationRelationships<'_> {

    pub async fn replication_configuration_rules_repository_filters_filter(&self) -> Vec<AwsEcrRepository> {
        vec![]
    }

}
"""
    )


def test_translate_relationships_with_multiple_fields(relationships, resources_schema):
    translated = translate_resource_relationships(
        "AWS::EC2::TransitGatewayRouteTablePropagation", relationships, resources_schema
    )
    assert (
        translated
        == """
#[derive(Serialize)]
pub struct AwsEc2TransitGatewayRouteTablePropagationRelationships<'a> {
    properties: &'a String,
}

#[Object(rename_fields = "PascalCase")]
impl AwsEc2TransitGatewayRouteTablePropagationRelationships<'_> {

    pub async fn transit_gateway_route_table_id(&self) -> AwsEc2TransitGatewayRouteTable {
        AwsEc2TransitGatewayRouteTable{
            id: "".to_string(),
            all_properties: "".to_string(),
        }
    }

    pub async fn transit_gateway_attachment_id(&self) -> AwsEc2TransitGatewayAttachment {
        AwsEc2TransitGatewayAttachment{
            id: "".to_string(),
            all_properties: "".to_string(),
        }
    }

}
"""
    )


def test_translate_relationships_with_multiple_relationships(
    relationships, resources_schema
):
    translated = translate_resource_relationships(
        "AWS::RAM::ResourceShare", relationships, resources_schema
    )
    assert (
        translated
        == """
#[derive(Serialize)]
pub struct AwsRamResourceShareRelationships<'a> {
    properties: &'a String,
}

#[derive(Union, Serialize)]
enum AwsRamResourceShareResourceArnsConnections {
    AwsEc2PrefixList(AwsEc2PrefixList),
    AwsRoute53ResolverFirewallRuleGroup(AwsRoute53ResolverFirewallRuleGroup),
    Node(Node),
}

#[Object(rename_fields = "PascalCase")]
impl AwsRamResourceShareRelationships<'_> {

    pub async fn resource_arns(&self) -> Vec<AwsRamResourceShareResourceArnsConnections> {
        vec![]
    }

    pub async fn tags_value(&self) -> Vec<Node> {
        vec![]
    }

}
"""
    )
