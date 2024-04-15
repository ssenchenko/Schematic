package translator

import (
	"strings"
	"testing"
	"text/template"
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
				RustPropertyName: "identifier",
				RustPropertyType: "String",
			},
			{
				RustPropertyName: "all_properties",
				RustPropertyType: "String",
			},
		},
		Relationships: nil,
	}

	runSnapshotTest(t, testData, snapshotFileName, templateFileName, TemplateDir, nil)
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

	runSnapshotTest(t, testData, snapshotFileName, templateFileName, TemplateDir, nil)
}

func TestSnapshotResourceUnion(t *testing.T) {
	templateFileName := "union_enum.go.tmpl"
	snapshotFileName := "union_enum.rs.snap"
	testData := ResourceUnion{
		RustUnionName: "AwsEc2InstanceConnections_SecurityGroupIds",
		RustResourceNames: map[string]bool{
			"AwsEc2SecurityGroup": true,
			"Node":                true,
		},
	}

	runSnapshotTest(t, testData, snapshotFileName, templateFileName, TemplateDir, nil)
}

func TestSnapshotAwsResourceImpl(t *testing.T) {
	templateFileName := "aws_resource_impl.go.tmpl"
	snapshotFileName := "aws_resource_impl.rs.snap"
	testData := []ResourceType{
		{
			CfnResourceName:     "AWS::CloudWatch::Alarm",
			RustResourceName:    "AwsCloudWatchAlarm",
			GraphQlResourceName: "Aws_CloudWatch_Alarm",
			Properties:          nil,
			Relationships:       nil,
		},
		{
			CfnResourceName:     "AWS::IAM::InstanceProfile",
			RustResourceName:    "AwsIamInstanceProfile",
			GraphQlResourceName: "Aws_Iam_InstanceProfile",
			Properties:          nil,
			Relationships:       nil,
		},
	}

	runSnapshotTest(t, testData, snapshotFileName, templateFileName, TemplateDir, nil)
}

func TestSnapshotRelationship(t *testing.T) {
	templateFileName := "relationship.go.tmpl"
	snapshotFileName := "relationship.rs.snap"
	testData := ResourceType{
		CfnResourceName:     "AWS::EC2::SecurityGroup",
		RustResourceName:    "AwsEc2SecurityGroup",
		GraphQlResourceName: "Aws_Ec2_SecurityGroup",
		Properties:          nil,
		Relationships: []ResourceRelationship{
			{
				RustSourcePropertyName: "security_group_ingress_source_security_group_name",
				RustReturnType:         "Vec<AwsEc2SecurityGroup>",
				RustGenericType:        "Vec<AwsEc2SecurityGroup>",
				TargetUnion:            ResourceUnion{},
			},
			{
				RustSourcePropertyName: "security_group_ingress_source_security_group_id",
				RustReturnType:         "Vec<AwsEc2SecurityGroupConnections_SecurityGroupIngressSourceSecurityGroupId>",
				RustGenericType:        "Vec<AwsEc2SecurityGroupConnections_SecurityGroupIngressSourceSecurityGroupId>",
				TargetUnion: ResourceUnion{
					RustUnionName: "AwsEc2SecurityGroupConnections_SecurityGroupIngressSourceSecurityGroupId",
					RustResourceNames: map[string]bool{
						"AwsEc2SecurityGroup": true,
						"Node":                true,
					},
				},
			},
		},
	}

	runSnapshotTest(
		t,
		testData,
		snapshotFileName,
		templateFileName,
		TemplateDir,
		template.FuncMap{"DerefResourceUnion": Deref[ResourceUnion]},
		`{{ define "union_enum.go.tmpl" }}// <{{ .RustUnionName }}> Mock{{ end }}`,
	)
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
					RustPropertyName: "identifier",
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
					TargetUnion: ResourceUnion{
						RustUnionName: "AwsCloudWatchAlarmConnections_DimensionsValue",
						RustResourceNames: map[string]bool{
							"AwsEc2Instance": true,
							"AwsS3Bucket":    true,
						},
					},
				},
				{
					RustSourcePropertyName: "metrics_metric_stat_metric_dimensions_value",
					RustReturnType:         "Vec<AwsCloudWatchAlarmConnections_MetricsMetricStatMetricDimensionsValue>",
					RustGenericType:        "Vec<AwsCloudWatchAlarmConnections_MetricsMetricStatMetricDimensionsValue>",
					TargetUnion: ResourceUnion{
						RustUnionName: "AwsCloudWatchAlarmConnections_MetricsMetricStatMetricDimensionsValue",
						RustResourceNames: map[string]bool{
							"AwsEc2Instance": true,
							"AwsS3Bucket":    true,
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
					RustPropertyName: "identifier",
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

	runSnapshotTest(
		t,
		testData,
		snapshotFileName,
		templateFileName,
		TemplateDir,
		template.FuncMap{"DerefResourceUnion": Deref[ResourceUnion]},
		`{{ define "interface_enum.go.tmpl" }}// Enum <Resource> Mock{{ end }}`,
		`{{ define "resource_struct.go.tmpl" }}// <{{ .RustResourceName }} Struct> Mock{{ end }}`,
		`{{ define "relationship.go.tmpl" }}// <{{ .RustResourceName }}Relationships> Mock{{ end }}`,
		`{{ define "aws_resource_impl.go.tmpl" }}// <impl AwsResource> Mock{{ end }}`,
		`{{ define "union_enum.go.tmpl" }}// <{{ .RustUnionName }}> Mock{{ end }}`,
	)
}

func runSnapshotTest[TestData any](
	t *testing.T,
	testData TestData,
	snapshotFileName string,
	templateName string,
	templateDir string,
	funcs template.FuncMap,
	nestedTemplates ...string,
) {
	expected, err := LoadSnapshot(snapshotFileName)
	if err != nil {
		t.Errorf("cannot load snapshot %s %v", snapshotFileName, err)
	}
	expected = strings.TrimSpace(expected)

	buffer, err := hydrateTemplate(testData, templateName, templateDir, funcs, nestedTemplates...)
	if err != nil {
		t.Errorf("For %s, expected no error but got\n%v", templateName, err)
	}
	actual := strings.TrimSpace(buffer.String())
	if expected != actual {
		t.Errorf("For %s, expected %s, but got\n%v", templateName, snapshotFileName, actual)
	}
}
