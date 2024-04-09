package translator

import (
	"strings"
	"testing"
)

func TestSnapshotResourceStruct(t *testing.T) {
	templateFileName := "resource_struct.go.tmpl"
	snapshotFileName := "resource_struct.rs.snap"
	testData := ResourceType{
		CfnResourceName:     "AWS::S3::Bucket",
		RustResourceName:    "AwsS3Bucket",
		GraphQlResourceName: "Aws_S3_Bucket",
		Properties: []ResourceProperty{
			{
				RustPropertyName: "id",
				RustPropertyType: "String",
			},
			{
				RustPropertyName: "all_properties",
				RustPropertyType: "String",
			},
		},
		Relationships: nil,
	}

	runSnapshotTest(t, snapshotFileName, testData, templateFileName)
}

func TestSnapshotResourceEnum(t *testing.T) {
	templateFileName := "interface_enum.go.tmpl"
	snapshotFileName := "interface_enum.rs.snap"
	testData := []ResourceType{
		{
			CfnResourceName:     "AWS::CloudWatch::Alarm",
			RustResourceName:    "AwsCloudWatchAlarm",
			GraphQlResourceName: "Aws_CloudWatch_Alarm",
			Properties:          nil,
			Relationships:       nil,
		},
		{
			CfnResourceName:     "AWS::EC2::Instance",
			RustResourceName:    "AwsEc2Instance",
			GraphQlResourceName: "Aws_Ec2_Instance",
			Properties:          nil,
			Relationships:       nil,
		},
	}

	runSnapshotTest(t, snapshotFileName, testData, templateFileName)
}

func TestSnapshotResourceUnion(t *testing.T) {
	templateFileName := "union_enum.go.tmpl"
	snapshotFileName := "union_enum.rs.snap"
	testData := ResourceUnion{
		RustUnionName: "AwsEc2InstanceConnections_SecurityGroupIds",
		Resources: []ResourceType{
			{
				CfnResourceName:     "AWS::EC2::SecurityGroup",
				RustResourceName:    "AwsEc2SecurityGroup",
				GraphQlResourceName: "AwsEc2SecurityGroup",
				Properties:          nil,
				Relationships:       nil,
			},
			{
				CfnResourceName:     "AWS::EC2::Subnet",
				RustResourceName:    "Node",
				GraphQlResourceName: "Node",
				Properties:          nil,
				Relationships:       nil,
			},
		},
	}

	runSnapshotTest(t, snapshotFileName, testData, templateFileName)
}

func TestSnapshotRelationship(t *testing.T) {
	templateFileName := "relationship.go.tmpl"
	snapshotFileName := "relationship.rs.snap"
	testData := ResourceType{
		CfnResourceName:     "AWS::EC2::SecurityGroup",
		RustResourceName:    "AwsEc2SecurityGroup",
		GraphQlResourceName: "AwsEc2SecurityGroup",
		Properties:          nil,
		Relationships: []ResourceRelationship{
			{
				RustSourcePropertyName: "security_group_ingress_source_security_group_name",
				RustReturnType:         "Vec<AwsEc2SecurityGroup>",
				RustGenericType:        "Vec<AwsEc2SecurityGroup>",
				TargetUnion:            nil,
			},
			{
				RustSourcePropertyName: "security_group_ingress_source_security_group_id",
				RustReturnType:         "Vec<AwsEc2SecurityGroupConnections_SecurityGroupIngressSourceSecurityGroupId>",
				RustGenericType:        "Vec<AwsEc2SecurityGroupConnections_SecurityGroupIngressSourceSecurityGroupId>",
				TargetUnion: &ResourceUnion{
					RustUnionName: "AwsEc2SecurityGroupConnections_SecurityGroupIngressSourceSecurityGroupId",
					Resources: []ResourceType{
						{
							CfnResourceName:     "AWS::EC2::SecurityGroup",
							RustResourceName:    "AwsEc2SecurityGroup",
							GraphQlResourceName: "AwsEc2SecurityGroup",
							Properties:          nil,
						},
						{
							CfnResourceName:     "AWS::EC2::Subnet",
							RustResourceName:    "Node",
							GraphQlResourceName: "Node",
							Properties:          nil,
						},
					},
				},
			},
		},
	}

	runSnapshotTest(t, snapshotFileName, testData, templateFileName, "union_enum.go.tmpl")
}

func TestSnapshotAll(t *testing.T) {
	templateFileName := "all.go.tmpl"
	snapshotFileName := "all.rs.snap"
	testData := RustModel{
		{
			CfnResourceName:     "AWS::CloudWatch::Alarm",
			RustResourceName:    "AwsCloudWatchAlarm",
			GraphQlResourceName: "Aws_CloudWatch_Alarm",
			Properties: []ResourceProperty{
				{
					RustPropertyName: "id",
					RustPropertyType: "String",
				},
				{
					RustPropertyName: "all_properties",
					RustPropertyType: "String",
				},
			},
			Relationships: []ResourceRelationship{
				{
					RustSourcePropertyName: "dimensions_value",
					RustReturnType:         "Vec<AwsCloudWatchAlarmConnections_DimensionsValue>",
					RustGenericType:        "Vec<AwsCloudWatchAlarmConnections_DimensionsValue>",
					TargetUnion: &ResourceUnion{
						RustUnionName: "AwsCloudWatchAlarmConnections_DimensionsValue",
						Resources: []ResourceType{
							{
								CfnResourceName:     "AWS::EC2::Instance",
								RustResourceName:    "AwsEc2Instance",
								GraphQlResourceName: "Aws_Ec2_Instance",
								Properties:          nil,
							},
							{
								CfnResourceName:     "AWS::S3::Bucket",
								RustResourceName:    "AwsS3Bucket",
								GraphQlResourceName: "Aws_S3_Bucket",
								Properties:          nil,
							},
						},
					},
				},
				{
					RustSourcePropertyName: "metrics_metric_stat_metric_dimensions_value",
					RustReturnType:         "Vec<AwsCloudWatchAlarmConnections_MetricsMetricStatMetricDimensionsValue>",
					RustGenericType:        "Vec<AwsCloudWatchAlarmConnections_MetricsMetricStatMetricDimensionsValue>",
					TargetUnion: &ResourceUnion{
						RustUnionName: "AwsCloudWatchAlarmConnections_MetricsMetricStatMetricDimensionsValue",
						Resources: []ResourceType{
							{
								CfnResourceName:     "AWS::EC2::Instance",
								RustResourceName:    "AwsEc2Instance",
								GraphQlResourceName: "Aws_Ec2_Instance",
								Properties:          nil,
							},
							{
								CfnResourceName:     "AWS::S3::Bucket",
								RustResourceName:    "AwsS3Bucket",
								GraphQlResourceName: "Aws_S3_Bucket",
								Properties:          nil,
							},
						},
					},
				},
			},
		},
		{
			CfnResourceName:     "AWS::EC2::Instance",
			RustResourceName:    "AwsEc2Instance",
			GraphQlResourceName: "Aws_Ec2_Instance",
			Properties: []ResourceProperty{
				{
					RustPropertyName: "id",
					RustPropertyType: "String",
				},
				{
					RustPropertyName: "all_properties",
					RustPropertyType: "String",
				},
			},
			Relationships: nil,
		},
	}

	runSnapshotTest(t, snapshotFileName, testData,
		templateFileName,
		"interface_enum.go.tmpl",
		"resource_struct.go.tmpl",
		"relationship.go.tmpl",
		"union_enum.go.tmpl",
	)
}

func runSnapshotTest[TestData any](
	t *testing.T,
	snapshotFileName string,
	testData TestData,
	templateNames ...string,
) {
	if len(templateNames) == 0 {
		t.Errorf("template name is required")
	}

	expected, err := loadSnapshot(snapshotFileName)
	if err != nil {
		t.Errorf("cannot load snapshot %s %v", snapshotFileName, err)
	}
	expected = strings.TrimSpace(expected)

	actual, err := hydrate(testData, templateNames...)
	actual = strings.TrimSpace(actual)
	if err != nil {
		t.Errorf("For %s, expected no error but got\n%v", templateNames[0], err)
	} else if expected != actual {
		t.Errorf("For %s, expected %s, but got\n%v", templateNames[0], snapshotFileName, actual)
	}
}
