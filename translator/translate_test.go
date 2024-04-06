package translator

import (
	"reflect"
	"testing"
)

func TestAllRelationshipApplyOverrides(t *testing.T) {
	// Test code here
	override := map[string]map[string]string{
		"AWS::EC2::Instance": {"VolumeAttachments": "Volumes"},
	}
	testCase := AllRelationships{
		"AWS::EC2::Instance": {
			PrimaryKeys: []string{"InstanceId"},
			Relationships: []Relationship{
				{
					"VolumeAttachments": []Reference{
						{
							TypeName:  "AWS::EC2::Volume",
							Attribute: "SomeAttribute",
						},
					},
				},
			},
		},
	}
	expected := AllRelationships{
		"AWS::EC2::Instance": {
			PrimaryKeys: []string{"InstanceId"},
			Relationships: []Relationship{
				{
					"Volumes": []Reference{
						{
							TypeName:  "AWS::EC2::Volume",
							Attribute: "SomeAttribute",
						},
					},
				},
			},
		},
	}

	testCase.ApplyOverrides(override)
	if !reflect.DeepEqual(testCase, expected) {
		t.Errorf("Expected: %v, Got: %v", expected, testCase)
	}
}
