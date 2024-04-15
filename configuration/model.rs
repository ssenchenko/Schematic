// =======================================================
// This file is generated. Do not edit manually!
// =======================================================
use async_graphql::{Context, Enum, Error, Interface, OutputType, SimpleObject, Result, Object, ComplexObject, Union};
use serde::Serialize;

pub const AWS_RESOURCE_TYPE: &str = "AWS::Resource";

#[derive(Enum, Copy, Clone, Eq, PartialEq)]
pub enum Region {
    ApSouth1,
    CaCentral1,
    EuCentral1,
    UsWest1,
    UsWest2,
    EuNorth1,
    EuWest3,
    EuWest2,
    EuWest1,
    ApNortheast3,
    ApNortheast2,
    ApNortheast1,
    SaEast1,
    ApSoutheast1,
    ApSoutheast2,
    UsEast1,
    UsEast2,
}

#[derive(InputObject)]
#[graphql(name = "Topology_Config_Input")]
pub struct TopologyConfigInput {
    region: Option<Region>,
}

#[derive(InputObject)]
#[graphql(name = "Resource_Config_Input")]
pub struct ResourceConfigInput {
    region: Option<Region>,
}

#[derive(InputObject)]
#[graphql(name = "Resources_Config_Input")]
pub struct ResourcesConfigInput {
    region: Option<Region>,
    pub filter: Option<String>,
}

#[derive(SimpleObject, Serialize)]
#[graphql(name = "Edge", rename_fields = "camelCase")]
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
#[graphql(name = "Topology", rename_fields = "camelCase")]
pub struct Topology {
    pub nodes: Vec<NodeInterface>,
    pub edges: Vec<Edge>,
}

#[derive(SimpleObject, Serialize, Clone)]
#[graphql(name = "AWS_Resource", rename_fields = "camelCase")]
pub struct AWSResource {
    pub r#type: String,
    pub resource_type: String,
    pub identifier: String,
    pub all_properties: String,
}

#[derive(Interface, Serialize)]
#[graphql(
    name = "Node_Interface",
    rename_fields = "camelCase",
    field(name = "type", ty = "String"),
    field(name = "identifier", ty = "String")
)]
pub enum NodeInterface {
    Node(Node),
    AWSResource(AWSResource),
}

#[derive(SimpleObject, Serialize)]
#[graphql(name = "Node", rename_fields = "camelCase")]
pub struct Node {
    pub r#type: String,
    pub identifier: String,
}

#[derive(SimpleObject, Serialize)]
#[graphql(name = "Aws_CloudWatch_Alarm", rename_fields = "camelCase", complex)]
pub struct AwsCloudWatchAlarm {
    pub identifier: String,
    pub all_properties: String,
}

#[ComplexObject(rename_fields = "camelCase")]
impl AwsCloudWatchAlarm {
    pub async fn r#type(&self) -> String {
        AWS_RESOURCE_TYPE.to_string()
    }
    pub async fn resource_type(&self) -> String {
        "AWS::CloudWatch::Alarm".to_string()
    }
}
#[derive(SimpleObject, Serialize)]
#[graphql(name = "Aws_Ec2_SecurityGroup", rename_fields = "camelCase", complex)]
pub struct AwsEc2SecurityGroup {
    pub identifier: String,
    pub all_properties: String,
}

#[ComplexObject(rename_fields = "camelCase")]
impl AwsEc2SecurityGroup {
    pub async fn r#type(&self) -> String {
        AWS_RESOURCE_TYPE.to_string()
    }
    pub async fn resource_type(&self) -> String {
        "AWS::EC2::SecurityGroup".to_string()
    }
}
#[derive(SimpleObject, Serialize)]
#[graphql(name = "Aws_Iam_Policy", rename_fields = "camelCase", complex)]
pub struct AwsIamPolicy {
    pub identifier: String,
    pub all_properties: String,
}

#[ComplexObject(rename_fields = "camelCase")]
impl AwsIamPolicy {
    pub async fn r#type(&self) -> String {
        AWS_RESOURCE_TYPE.to_string()
    }
    pub async fn resource_type(&self) -> String {
        "AWS::IAM::Policy".to_string()
    }
}
#[derive(SimpleObject, Serialize)]
#[graphql(name = "Aws_S3_Bucket", rename_fields = "camelCase", complex)]
pub struct AwsS3Bucket {
    pub identifier: String,
    pub all_properties: String,
}

#[ComplexObject(rename_fields = "camelCase")]
impl AwsS3Bucket {
    pub async fn r#type(&self) -> String {
        AWS_RESOURCE_TYPE.to_string()
    }
    pub async fn resource_type(&self) -> String {
        "AWS::S3::Bucket".to_string()
    }
}
#[derive(SimpleObject, Serialize)]
#[graphql(name = "Aws_Ec2_VolumeAttachment", rename_fields = "camelCase", complex)]
pub struct AwsEc2VolumeAttachment {
    pub identifier: String,
    pub all_properties: String,
}

#[ComplexObject(rename_fields = "camelCase")]
impl AwsEc2VolumeAttachment {
    pub async fn r#type(&self) -> String {
        AWS_RESOURCE_TYPE.to_string()
    }
    pub async fn resource_type(&self) -> String {
        "AWS::EC2::VolumeAttachment".to_string()
    }
}
#[derive(SimpleObject, Serialize)]
#[graphql(name = "Aws_Iam_InstanceProfile", rename_fields = "camelCase", complex)]
pub struct AwsIamInstanceProfile {
    pub identifier: String,
    pub all_properties: String,
}

#[ComplexObject(rename_fields = "camelCase")]
impl AwsIamInstanceProfile {
    pub async fn r#type(&self) -> String {
        AWS_RESOURCE_TYPE.to_string()
    }
    pub async fn resource_type(&self) -> String {
        "AWS::IAM::InstanceProfile".to_string()
    }
}
#[derive(SimpleObject, Serialize)]
#[graphql(name = "Aws_Ec2_Instance", rename_fields = "camelCase", complex)]
pub struct AwsEc2Instance {
    pub identifier: String,
    pub all_properties: String,
}

#[ComplexObject(rename_fields = "camelCase")]
impl AwsEc2Instance {
    pub async fn r#type(&self) -> String {
        AWS_RESOURCE_TYPE.to_string()
    }
    pub async fn resource_type(&self) -> String {
        "AWS::EC2::Instance".to_string()
    }
}
#[derive(SimpleObject, Serialize)]
#[graphql(name = "Aws_Iam_Role", rename_fields = "camelCase", complex)]
pub struct AwsIamRole {
    pub identifier: String,
    pub all_properties: String,
}

#[ComplexObject(rename_fields = "camelCase")]
impl AwsIamRole {
    pub async fn r#type(&self) -> String {
        AWS_RESOURCE_TYPE.to_string()
    }
    pub async fn resource_type(&self) -> String {
        "AWS::IAM::Role".to_string()
    }
}

#[derive(Interface, Serialize)]
#[graphql(
    name = "AWS_Resource_Interface",
    rename_fields = "camelCase",
    field(name = "type", ty = "String"),
    field(name = "resource_type", ty = "String"),
    field(name = "identifier", ty = "String"),
    field(name = "all_properties", ty = "String")
)]
pub enum Resource {
    AwsCloudWatchAlarm(AwsCloudWatchAlarm),
    AwsEc2SecurityGroup(AwsEc2SecurityGroup),
    AwsIamPolicy(AwsIamPolicy),
    AwsS3Bucket(AwsS3Bucket),
    AwsEc2VolumeAttachment(AwsEc2VolumeAttachment),
    AwsIamInstanceProfile(AwsIamInstanceProfile),
    AwsEc2Instance(AwsEc2Instance),
    AwsIamRole(AwsIamRole),
}

impl AwsResource {
    pub fn get_resource(&self) -> AwsResourceInterface {
        match self.resource_type.as_str() {
            "AWS::CloudWatch::Alarm" => {
                AwsResourceInterface::AwsCloudWatchAlarm(AwsCloudWatchAlarm {
                    identifier: self.identifier.clone(),
                    all_properties: self.all_properties.clone(),
                })
            }
            "AWS::EC2::SecurityGroup" => {
                AwsResourceInterface::AwsEc2SecurityGroup(AwsEc2SecurityGroup {
                    identifier: self.identifier.clone(),
                    all_properties: self.all_properties.clone(),
                })
            }
            "AWS::IAM::Policy" => {
                AwsResourceInterface::AwsIamPolicy(AwsIamPolicy {
                    identifier: self.identifier.clone(),
                    all_properties: self.all_properties.clone(),
                })
            }
            "AWS::S3::Bucket" => {
                AwsResourceInterface::AwsS3Bucket(AwsS3Bucket {
                    identifier: self.identifier.clone(),
                    all_properties: self.all_properties.clone(),
                })
            }
            "AWS::EC2::VolumeAttachment" => {
                AwsResourceInterface::AwsEc2VolumeAttachment(AwsEc2VolumeAttachment {
                    identifier: self.identifier.clone(),
                    all_properties: self.all_properties.clone(),
                })
            }
            "AWS::IAM::InstanceProfile" => {
                AwsResourceInterface::AwsIamInstanceProfile(AwsIamInstanceProfile {
                    identifier: self.identifier.clone(),
                    all_properties: self.all_properties.clone(),
                })
            }
            "AWS::EC2::Instance" => {
                AwsResourceInterface::AwsEc2Instance(AwsEc2Instance {
                    identifier: self.identifier.clone(),
                    all_properties: self.all_properties.clone(),
                })
            }
            "AWS::IAM::Role" => {
                AwsResourceInterface::AwsIamRole(AwsIamRole {
                    identifier: self.identifier.clone(),
                    all_properties: self.all_properties.clone(),
                })
            }
            _ => AwsResourceInterface::AwsResource(self.clone()),
        }
    }
}

// =========== Relationships ===========

#[derive(Serialize)]
pub struct AwsCloudWatchAlarmRelationships<'a> {
    properties: &'a String,
}

impl AwsCloudWatchAlarm {
    pub async fn relationships(&self) -> AwsCloudWatchAlarmRelationships {
        AwsCloudWatchAlarmRelationships {
            properties: &self.all_properties,
        }
    }
}

#[derive(Union, Serialize)]
pub enum AwsCloudWatchAlarmConnections_DimensionsValue {
    AwsEc2Instance(AwsEc2Instance),
    AwsS3Bucket(AwsS3Bucket),
}#[derive(Union, Serialize)]
pub enum AwsCloudWatchAlarmConnections_MetricsMetricStatMetricDimensionsValue {
    AwsEc2Instance(AwsEc2Instance),
    AwsS3Bucket(AwsS3Bucket),
}

#[Object(
    name = "Aws_CloudWatch_Alarm_Relationships",
    rename_fields = "camelCase"
)]
impl AwsCloudWatchAlarmRelationships<'_> {
    pub async fn dimensions_value(&self, ctx: &Context<'_>) -> Vec<AwsCloudWatchAlarmConnections_DimensionsValue> {
        let atx_context = ctx.data::<AtsContext>()?;
        get_related_resource::<Vec<AwsCloudWatchAlarmConnections_DimensionsValue>>(
            &atx_context,
            &self.identifier,
            &self.resource_type,
            &self.properties,
        ).await?
    }
    pub async fn metrics_label(&self, ctx: &Context<'_>) -> Vec<AwsS3Bucket> {
        let atx_context = ctx.data::<AtsContext>()?;
        get_related_resource::<Vec<AwsS3Bucket>>(
            &atx_context,
            &self.identifier,
            &self.resource_type,
            &self.properties,
        ).await?
    }
    pub async fn metrics_metric_stat_metric_dimensions_value(&self, ctx: &Context<'_>) -> Vec<AwsCloudWatchAlarmConnections_MetricsMetricStatMetricDimensionsValue> {
        let atx_context = ctx.data::<AtsContext>()?;
        get_related_resource::<Vec<AwsCloudWatchAlarmConnections_MetricsMetricStatMetricDimensionsValue>>(
            &atx_context,
            &self.identifier,
            &self.resource_type,
            &self.properties,
        ).await?
    }
}


#[derive(Serialize)]
pub struct AwsEc2SecurityGroupRelationships<'a> {
    properties: &'a String,
}

impl AwsEc2SecurityGroup {
    pub async fn relationships(&self) -> AwsEc2SecurityGroupRelationships {
        AwsEc2SecurityGroupRelationships {
            properties: &self.all_properties,
        }
    }
}

#[derive(Union, Serialize)]
pub enum AwsEc2SecurityGroupConnections_SecurityGroupIngressSourceSecurityGroupId {
    AwsEc2SecurityGroup(AwsEc2SecurityGroup),
    Node(Node),
}

#[Object(
    name = "Aws_Ec2_SecurityGroup_Relationships",
    rename_fields = "camelCase"
)]
impl AwsEc2SecurityGroupRelationships<'_> {
    pub async fn vpc_id(&self, ctx: &Context<'_>) -> Option<Node> {
        let atx_context = ctx.data::<AtsContext>()?;
        get_related_resource::<Node>(
            &atx_context,
            &self.identifier,
            &self.resource_type,
            &self.properties,
        ).await?
    }
    pub async fn security_group_ingress_cidr_ip(&self, ctx: &Context<'_>) -> Vec<Node> {
        let atx_context = ctx.data::<AtsContext>()?;
        get_related_resource::<Vec<Node>>(
            &atx_context,
            &self.identifier,
            &self.resource_type,
            &self.properties,
        ).await?
    }
    pub async fn security_group_ingress_source_security_group_name(&self, ctx: &Context<'_>) -> Vec<AwsEc2SecurityGroup> {
        let atx_context = ctx.data::<AtsContext>()?;
        get_related_resource::<Vec<AwsEc2SecurityGroup>>(
            &atx_context,
            &self.identifier,
            &self.resource_type,
            &self.properties,
        ).await?
    }
    pub async fn security_group_ingress_source_security_group_id(&self, ctx: &Context<'_>) -> Vec<AwsEc2SecurityGroupConnections_SecurityGroupIngressSourceSecurityGroupId> {
        let atx_context = ctx.data::<AtsContext>()?;
        get_related_resource::<Vec<AwsEc2SecurityGroupConnections_SecurityGroupIngressSourceSecurityGroupId>>(
            &atx_context,
            &self.identifier,
            &self.resource_type,
            &self.properties,
        ).await?
    }
    pub async fn security_group_egress_cidr_ip(&self, ctx: &Context<'_>) -> Vec<Node> {
        let atx_context = ctx.data::<AtsContext>()?;
        get_related_resource::<Vec<Node>>(
            &atx_context,
            &self.identifier,
            &self.resource_type,
            &self.properties,
        ).await?
    }
    pub async fn security_group_egress_destination_security_group_id(&self, ctx: &Context<'_>) -> Vec<AwsEc2SecurityGroup> {
        let atx_context = ctx.data::<AtsContext>()?;
        get_related_resource::<Vec<AwsEc2SecurityGroup>>(
            &atx_context,
            &self.identifier,
            &self.resource_type,
            &self.properties,
        ).await?
    }
}


#[derive(Serialize)]
pub struct AwsIamPolicyRelationships<'a> {
    properties: &'a String,
}

impl AwsIamPolicy {
    pub async fn relationships(&self) -> AwsIamPolicyRelationships {
        AwsIamPolicyRelationships {
            properties: &self.all_properties,
        }
    }
}



#[Object(
    name = "Aws_Iam_Policy_Relationships",
    rename_fields = "camelCase"
)]
impl AwsIamPolicyRelationships<'_> {
    pub async fn roles(&self, ctx: &Context<'_>) -> Vec<Node> {
        let atx_context = ctx.data::<AtsContext>()?;
        get_related_resource::<Vec<Node>>(
            &atx_context,
            &self.identifier,
            &self.resource_type,
            &self.properties,
        ).await?
    }
}


#[derive(Serialize)]
pub struct AwsS3BucketRelationships<'a> {
    properties: &'a String,
}

impl AwsS3Bucket {
    pub async fn relationships(&self) -> AwsS3BucketRelationships {
        AwsS3BucketRelationships {
            properties: &self.all_properties,
        }
    }
}



#[Object(
    name = "Aws_S3_Bucket_Relationships",
    rename_fields = "camelCase"
)]
impl AwsS3BucketRelationships<'_> {
    pub async fn analytics_configurations_storage_class_analysis_data_export_destination_bucket_arn(&self, ctx: &Context<'_>) -> Vec<AwsS3Bucket> {
        let atx_context = ctx.data::<AtsContext>()?;
        get_related_resource::<Vec<AwsS3Bucket>>(
            &atx_context,
            &self.identifier,
            &self.resource_type,
            &self.properties,
        ).await?
    }
    pub async fn inventory_configurations_destination_bucket_arn(&self, ctx: &Context<'_>) -> Vec<AwsS3Bucket> {
        let atx_context = ctx.data::<AtsContext>()?;
        get_related_resource::<Vec<AwsS3Bucket>>(
            &atx_context,
            &self.identifier,
            &self.resource_type,
            &self.properties,
        ).await?
    }
    pub async fn logging_configuration_destination_bucket_name(&self, ctx: &Context<'_>) -> Option<AwsS3Bucket> {
        let atx_context = ctx.data::<AtsContext>()?;
        get_related_resource::<AwsS3Bucket>(
            &atx_context,
            &self.identifier,
            &self.resource_type,
            &self.properties,
        ).await?
    }
    pub async fn replication_configuration_role(&self, ctx: &Context<'_>) -> Option<Node> {
        let atx_context = ctx.data::<AtsContext>()?;
        get_related_resource::<Node>(
            &atx_context,
            &self.identifier,
            &self.resource_type,
            &self.properties,
        ).await?
    }
    pub async fn replication_configuration_rules_destination_bucket(&self, ctx: &Context<'_>) -> Vec<AwsS3Bucket> {
        let atx_context = ctx.data::<AtsContext>()?;
        get_related_resource::<Vec<AwsS3Bucket>>(
            &atx_context,
            &self.identifier,
            &self.resource_type,
            &self.properties,
        ).await?
    }
}


#[derive(Serialize)]
pub struct AwsEc2VolumeAttachmentRelationships<'a> {
    properties: &'a String,
}

impl AwsEc2VolumeAttachment {
    pub async fn relationships(&self) -> AwsEc2VolumeAttachmentRelationships {
        AwsEc2VolumeAttachmentRelationships {
            properties: &self.all_properties,
        }
    }
}



#[Object(
    name = "Aws_Ec2_VolumeAttachment_Relationships",
    rename_fields = "camelCase"
)]
impl AwsEc2VolumeAttachmentRelationships<'_> {
    pub async fn volume_id(&self, ctx: &Context<'_>) -> Option<Node> {
        let atx_context = ctx.data::<AtsContext>()?;
        get_related_resource::<Node>(
            &atx_context,
            &self.identifier,
            &self.resource_type,
            &self.properties,
        ).await?
    }
    pub async fn instance_id(&self, ctx: &Context<'_>) -> Option<AwsEc2Instance> {
        let atx_context = ctx.data::<AtsContext>()?;
        get_related_resource::<AwsEc2Instance>(
            &atx_context,
            &self.identifier,
            &self.resource_type,
            &self.properties,
        ).await?
    }
}


#[derive(Serialize)]
pub struct AwsIamInstanceProfileRelationships<'a> {
    properties: &'a String,
}

impl AwsIamInstanceProfile {
    pub async fn relationships(&self) -> AwsIamInstanceProfileRelationships {
        AwsIamInstanceProfileRelationships {
            properties: &self.all_properties,
        }
    }
}



#[Object(
    name = "Aws_Iam_InstanceProfile_Relationships",
    rename_fields = "camelCase"
)]
impl AwsIamInstanceProfileRelationships<'_> {
    pub async fn roles(&self, ctx: &Context<'_>) -> Vec<Node> {
        let atx_context = ctx.data::<AtsContext>()?;
        get_related_resource::<Vec<Node>>(
            &atx_context,
            &self.identifier,
            &self.resource_type,
            &self.properties,
        ).await?
    }
    pub async fn instance_profile_name(&self, ctx: &Context<'_>) -> Option<Node> {
        let atx_context = ctx.data::<AtsContext>()?;
        get_related_resource::<Node>(
            &atx_context,
            &self.identifier,
            &self.resource_type,
            &self.properties,
        ).await?
    }
}


#[derive(Serialize)]
pub struct AwsEc2InstanceRelationships<'a> {
    properties: &'a String,
}

impl AwsEc2Instance {
    pub async fn relationships(&self) -> AwsEc2InstanceRelationships {
        AwsEc2InstanceRelationships {
            properties: &self.all_properties,
        }
    }
}

#[derive(Union, Serialize)]
pub enum AwsEc2InstanceConnections_SecurityGroupIds {
    AwsEc2SecurityGroup(AwsEc2SecurityGroup),
    Node(Node),
}

#[Object(
    name = "Aws_Ec2_Instance_Relationships",
    rename_fields = "camelCase"
)]
impl AwsEc2InstanceRelationships<'_> {
    pub async fn security_groups(&self, ctx: &Context<'_>) -> Vec<AwsEc2SecurityGroup> {
        let atx_context = ctx.data::<AtsContext>()?;
        get_related_resource::<Vec<AwsEc2SecurityGroup>>(
            &atx_context,
            &self.identifier,
            &self.resource_type,
            &self.properties,
        ).await?
    }
    pub async fn iam_instance_profile(&self, ctx: &Context<'_>) -> Option<AwsIamInstanceProfile> {
        let atx_context = ctx.data::<AtsContext>()?;
        get_related_resource::<AwsIamInstanceProfile>(
            &atx_context,
            &self.identifier,
            &self.resource_type,
            &self.properties,
        ).await?
    }
    pub async fn network_interfaces_group_set(&self, ctx: &Context<'_>) -> Vec<AwsEc2SecurityGroup> {
        let atx_context = ctx.data::<AtsContext>()?;
        get_related_resource::<Vec<AwsEc2SecurityGroup>>(
            &atx_context,
            &self.identifier,
            &self.resource_type,
            &self.properties,
        ).await?
    }
    pub async fn security_group_ids(&self, ctx: &Context<'_>) -> Vec<AwsEc2InstanceConnections_SecurityGroupIds> {
        let atx_context = ctx.data::<AtsContext>()?;
        get_related_resource::<Vec<AwsEc2InstanceConnections_SecurityGroupIds>>(
            &atx_context,
            &self.identifier,
            &self.resource_type,
            &self.properties,
        ).await?
    }
    pub async fn volumes(&self, ctx: &Context<'_>) -> Vec<AwsEc2VolumeAttachment> {
        let atx_context = ctx.data::<AtsContext>()?;
        get_related_resource::<Vec<AwsEc2VolumeAttachment>>(
            &atx_context,
            &self.identifier,
            &self.resource_type,
            &self.properties,
        ).await?
    }
}


