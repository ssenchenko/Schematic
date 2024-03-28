// =======================================================
// This file is generated.  Do not edit manually!
// =======================================================
use async_graphql::{Context, Enum, Error, Interface, OutputType, SimpleObject, Result, Object, ComplexObject, Union};
use serde::Serialize;

#[derive(SimpleObject, Serialize)]
#[graphql(name = "Aws_CloudWatch_Alarm", rename_fields = "PascalCase", complex)]
pub struct AwsCloudWatchAlarm {
    pub id: String,
    pub all_properties: String,
}

#[ComplexObject(rename_fields = "PascalCase")]
impl AwsCloudWatchAlarm {
    pub async fn type_name(&self) -> String {
        "AWS::CloudWatch::Alarm".to_string()
    }
}

#[derive(SimpleObject, Serialize)]
#[graphql(name = "Aws_Ec2_Instance", rename_fields = "PascalCase", complex)]
pub struct AwsEc2Instance {
    pub id: String,
    pub all_properties: String,
}

#[ComplexObject(rename_fields = "PascalCase")]
impl AwsEc2Instance {
    pub async fn type_name(&self) -> String {
        "AWS::EC2::Instance".to_string()
    }
}

#[derive(SimpleObject, Serialize)]
#[graphql(name = "Aws_Ec2_SecurityGroup", rename_fields = "PascalCase", complex)]
pub struct AwsEc2SecurityGroup {
    pub id: String,
    pub all_properties: String,
}

#[ComplexObject(rename_fields = "PascalCase")]
impl AwsEc2SecurityGroup {
    pub async fn type_name(&self) -> String {
        "AWS::EC2::SecurityGroup".to_string()
    }
}

#[derive(SimpleObject, Serialize)]
#[graphql(name = "Aws_Ec2_VolumeAttachment", rename_fields = "PascalCase", complex)]
pub struct AwsEc2VolumeAttachment {
    pub id: String,
    pub all_properties: String,
}

#[ComplexObject(rename_fields = "PascalCase")]
impl AwsEc2VolumeAttachment {
    pub async fn type_name(&self) -> String {
        "AWS::EC2::VolumeAttachment".to_string()
    }
}

#[derive(SimpleObject, Serialize)]
#[graphql(name = "Aws_Iam_InstanceProfile", rename_fields = "PascalCase", complex)]
pub struct AwsIamInstanceProfile {
    pub id: String,
    pub all_properties: String,
}

#[ComplexObject(rename_fields = "PascalCase")]
impl AwsIamInstanceProfile {
    pub async fn type_name(&self) -> String {
        "AWS::IAM::InstanceProfile".to_string()
    }
}

#[derive(SimpleObject, Serialize)]
#[graphql(name = "Aws_Iam_Policy", rename_fields = "PascalCase", complex)]
pub struct AwsIamPolicy {
    pub id: String,
    pub all_properties: String,
}

#[ComplexObject(rename_fields = "PascalCase")]
impl AwsIamPolicy {
    pub async fn type_name(&self) -> String {
        "AWS::IAM::Policy".to_string()
    }
}

#[derive(SimpleObject, Serialize)]
#[graphql(name = "Aws_Iam_Role", rename_fields = "PascalCase", complex)]
pub struct AwsIamRole {
    pub id: String,
    pub all_properties: String,
}

#[ComplexObject(rename_fields = "PascalCase")]
impl AwsIamRole {
    pub async fn type_name(&self) -> String {
        "AWS::IAM::Role".to_string()
    }
}

#[derive(SimpleObject, Serialize)]
#[graphql(name = "Aws_S3_Bucket", rename_fields = "PascalCase", complex)]
pub struct AwsS3Bucket {
    pub id: String,
    pub all_properties: String,
}

#[ComplexObject(rename_fields = "PascalCase")]
impl AwsS3Bucket {
    pub async fn type_name(&self) -> String {
        "AWS::S3::Bucket".to_string()
    }
}

#[derive(Interface, Serialize)]
#[graphql(
    name = "Resource",
    rename_fields = "PascalCase",
    field(name = "id", ty = "String"),
    field(name = "type_name", ty = "String"),
    field(name = "all_properties", ty = "String"),
)]
pub enum Resource {
    AwsCloudWatchAlarm(AwsCloudWatchAlarm),
    AwsEc2Instance(AwsEc2Instance),
    AwsEc2SecurityGroup(AwsEc2SecurityGroup),
    AwsEc2VolumeAttachment(AwsEc2VolumeAttachment),
    AwsIamInstanceProfile(AwsIamInstanceProfile),
    AwsIamPolicy(AwsIamPolicy),
    AwsIamRole(AwsIamRole),
    AwsS3Bucket(AwsS3Bucket),
    Node(Node)
}

#[derive(SimpleObject, Serialize, Clone)]
#[graphql(name = "Node", rename_fields = "PascalCase")]
pub struct Node {
    pub id: String,
    pub type_name: String,
    pub all_properties: String,
}

#[derive(SimpleObject, Serialize)]
#[graphql(name = "Edge", rename_fields = "PascalCase")]
pub struct Edge {
    source: String,
    target: String,
    relation: Relation,
}

#[derive(Enum, Copy, Clone, Eq, PartialEq, Serialize)]
pub enum Relation {
    IsRelatedTo,
}

#[derive(SimpleObject, Serialize)]
#[graphql(name = "Topology", rename_fields = "PascalCase")]
pub struct Topology {
    nodes: Vec<Node>,
    edges: Vec<Edge>,
}

// =========== Relationships ===========
