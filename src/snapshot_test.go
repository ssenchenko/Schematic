package schematic

import (
	"testing"
)

func TestHydrateResourceStruct(t *testing.T) {
	templateFileName := "resource_struct.go.tmpl"
	snapshotFileName := "resource_struct.rs.snap"
	testData := ResourceType{
		CfnResourceName:     "AWS::S3::Bucket",
		RustResourceName:    "AwsS3Bucket",
		GraphQlResourceName: "Aws_S3_Bucket",
		UseComplex:          ", complex",
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
	}

	runSnapshotTest(t, testData, templateFileName, snapshotFileName)
}

func TestHydrateResourceEnum(t *testing.T) {
	templateFileName := "interface_enum.go.tmpl"
	snapshotFileName := "interface_enum.rs.snap"
	testData := []ResourceType{
		{
			CfnResourceName:     "AWS::CloudWatch::Alarm",
			RustResourceName:    "AwsCloudWatchAlarm",
			GraphQlResourceName: "Aws_CloudWatch_Alarm",
			UseComplex:          "",
			Properties:          nil,
		},
		{
			CfnResourceName:     "AWS::EC2::Instance",
			RustResourceName:    "AwsEc2Instance",
			GraphQlResourceName: "Aws_Ec2_Instance",
			UseComplex:          "",
			Properties:          nil,
		},
	}

	runSnapshotTest(t, testData, templateFileName, snapshotFileName)
}

func TestHydrateResourceUnion(t *testing.T) {
	templateFileName := "union_enum.go.tmpl"
	snapshotFileName := "union_enum.rs.snap"
	testData := ResourceUnion{
		RustUnionName: "AwsEc2InstanceSecurityGroupIdsConnections",
		Resources: []ResourceType{
			{
				CfnResourceName:     "AWS::EC2::SecurityGroup",
				RustResourceName:    "AwsEc2SecurityGroup",
				GraphQlResourceName: "AwsEc2SecurityGroup",
				UseComplex:          ", complex",
				Properties:          nil,
			},
			{
				CfnResourceName:     "AWS::EC2::Subnet",
				RustResourceName:    "Node",
				GraphQlResourceName: "Node",
				UseComplex:          "",
				Properties:          nil,
			},
		},
	}

	runSnapshotTest(t, testData, templateFileName, snapshotFileName)
}

func runSnapshotTest[TestData any](t *testing.T, testData TestData, templateFileName string, snapshotFileName string) {
	expected, err := loadSnapshot(snapshotFileName)
	if err != nil {
		t.Errorf("cannot load snapshot %s %v", snapshotFileName, err)
	}

	actual, err := hydrate(testData, templateFileName)
	if err != nil {
		t.Errorf("For %s, expected no error but got %v", templateFileName, err)
	} else if actual != expected {
		t.Errorf("For %s, expected %s, but got %v", templateFileName, snapshotFileName, actual)
	}
}
