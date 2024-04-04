package schematic

import (
	"testing"
)

func TestHydrateResourceStruct(t *testing.T) {
	templateFileName := "resource_struct.go.tmpl"
	snapshotFileName := "resource_struct.rs.snap"
	testResource := ResourceType{
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

	expected, err := loadSnapshot(snapshotFileName)
	if err != nil {
		t.Errorf("cannot load snapshot %s %v", snapshotFileName, err)
	}

	actual, err := hydrate[ResourceType](testResource, templateFileName)
	if err != nil {
		t.Errorf("For %s, expected no error but got %v", templateFileName, err)
	} else if actual != expected {
		t.Errorf("For %s, expected %s, but got %v", templateFileName, snapshotFileName, actual)
	}
}


