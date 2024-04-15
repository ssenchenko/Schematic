package translator

import (
	"fmt"
	"reflect"
	"strings"
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
				{
					// to make sure that Volumes is not added to all relationships
					"OtherKey": []Reference{
						{
							TypeName:  "AWS::EC2::BlahBlah",
							Attribute: "SomeOtherAttribute",
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
				{
					// to make sure that Volumes is not added to all relationships
					"OtherKey": []Reference{
						{
							TypeName:  "AWS::EC2::BlahBlah",
							Attribute: "SomeOtherAttribute",
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

func TestIsInFilter(t *testing.T) {
	testCases := []struct {
		element  string
		filter   map[string]bool
		expected bool
	}{
		{
			element:  "i-am-in",
			filter:   map[string]bool{"i-am-in": true, "i-am-in-too": true},
			expected: true,
		},
		{
			element:  "i-am-not-in",
			filter:   map[string]bool{"i-am-in": true, "i-am-in-too": true},
			expected: false,
		},
		{
			element:  "i-am-in",
			filter:   make(map[string]bool),
			expected: true,
		},
	}
	for _, testCase := range testCases {
		actual := isInFilter(testCase.element, testCase.filter)
		if actual != testCase.expected {
			t.Errorf("For %s in filter %v expected %v but got %v", testCase.element, testCase.filter, testCase.expected, actual)
		}
	}
}

func TestAllRelationships_HasRelationships(t *testing.T) {
	rel := AllRelationships{
		"WithRelationshipsInFilter": Resource{
			PrimaryKeys: nil,
			Relationships: []Relationship{
				{
					"Property": []Reference{
						{
							TypeName:  "InFilter",
							Attribute: "blah",
						},
					},
				},
			},
		},
		"WithoutRelationshipsInFilter": Resource{
			PrimaryKeys:   nil,
			Relationships: make([]Relationship, 0),
		},
		"WithRelationshipNotInFilter": Resource{
			PrimaryKeys: nil,
			Relationships: []Relationship{
				{
					"Property": []Reference{
						{
							TypeName:  "NotInFilter",
							Attribute: "blah",
						},
					},
				},
			},
		},
	}
	filter := map[string]bool{"InFilter": true}
	testCases := []struct {
		name     string
		expected bool
	}{
		{
			name:     "WithRelationshipsInFilter",
			expected: true,
		},
		{
			name:     "WithoutRelationshipsInFilter",
			expected: false,
		},
		{
			name:     "WithRelationshipNotInFilter",
			expected: false,
		},
	}

	for _, testCase := range testCases {
		actual := rel.HasRelationships(testCase.name, filter)
		if actual != testCase.expected {
			t.Errorf("For %s expected %v but got %v", testCase.name, testCase.expected, actual)
		}
	}
}

func prepareTranslateTestData(t *testing.T) (map[string]map[string]any, AllRelationships) {
	resources := []string{
		"AWS::Connect::UserHierarchyGroup",
		"AWS::EC2::TransitGatewayRouteTablePropagation",
		"AWS::RAM::ResourceShare",
		"AWS::ECR::ReplicationConfiguration",
	}

	files := make([]string, len(resources))
	for i, resource := range resources {
		res := strings.ReplaceAll(resource, "::", "-")
		res = strings.ToLower(res)
		files[i] = fmt.Sprintf("%s.json", res)
	}

	cfnSchemaTest, err := LoadCfnSchemaCombined("../data/cfn", files)
	if err != nil {
		t.Errorf("Error loading combined schema: %v", err)
	}

	allRelationships, err := LoadAllRelationships("../data/all-schema-combined.json")
	if err != nil {
		t.Errorf("Error loading all relationships: %v", err)
	}

	testRelationships := make(AllRelationships)
	for _, resource := range resources {
		testRelationships[resource] = allRelationships[resource]
		for _, relationship := range allRelationships[resource].Relationships {
			for _, references := range relationship {
				for _, ref := range references {
					testRelationships[ref.TypeName] = allRelationships[ref.TypeName]
				}
			}
		}
	}

	return cfnSchemaTest, testRelationships
}

func TestTranslate(t *testing.T) {
	cfnSchemaTest, testRelationships := prepareTranslateTestData(t)
	testCase := []struct {
		description      string
		resourceName     string
		expectedResource ResourceType
		expectedErrors   []error
	}{
		{
			description:  "resource with single relationship",
			resourceName: "AWS::Connect::UserHierarchyGroup",
			expectedResource: ResourceType{
				CfnResourceName:     "AWS::Connect::UserHierarchyGroup",
				RustResourceName:    "AwsConnectUserHierarchyGroup",
				GraphQlResourceName: "Aws_Connect_UserHierarchyGroup",
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
						RustSourcePropertyName: "instance_arn",
						RustReturnType:         "Option<Node>",
						RustGenericType:        "Node",
						TargetUnion:            ResourceUnion{},
					},
				},
			},
			expectedErrors: nil,
		},
		{
			description:  "resource with single array relationship and nested property",
			resourceName: "AWS::ECR::ReplicationConfiguration",
			expectedResource: ResourceType{
				CfnResourceName:     "AWS::ECR::ReplicationConfiguration",
				RustResourceName:    "AwsEcrReplicationConfiguration",
				GraphQlResourceName: "Aws_Ecr_ReplicationConfiguration",
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
						RustSourcePropertyName: "replication_configuration_rules_repository_filters_filter",
						RustReturnType:         "Vec<AwsEcrRepository>",
						RustGenericType:        "Vec<AwsEcrRepository>",
						TargetUnion:            ResourceUnion{},
					},
				},
			},
			expectedErrors: nil,
		},
		{
			description:  "resource with multiple relationships",
			resourceName: "AWS::EC2::TransitGatewayRouteTablePropagation",
			expectedResource: ResourceType{
				CfnResourceName:     "AWS::EC2::TransitGatewayRouteTablePropagation",
				RustResourceName:    "AwsEc2TransitGatewayRouteTablePropagation",
				GraphQlResourceName: "Aws_Ec2_TransitGatewayRouteTablePropagation",
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
						RustSourcePropertyName: "transit_gateway_route_table_id",
						RustReturnType:         "Option<AwsEc2TransitGatewayRouteTable>",
						RustGenericType:        "AwsEc2TransitGatewayRouteTable",
						TargetUnion:            ResourceUnion{},
					},
					{
						RustSourcePropertyName: "transit_gateway_attachment_id",
						RustReturnType:         "Option<AwsEc2TransitGatewayAttachment>",
						RustGenericType:        "AwsEc2TransitGatewayAttachment",
						TargetUnion:            ResourceUnion{},
					},
				},
			},
			expectedErrors: nil,
		},
		{
			description:  "resource with multiple relationships and one union return type",
			resourceName: "AWS::RAM::ResourceShare",
			expectedResource: ResourceType{
				CfnResourceName:     "AWS::RAM::ResourceShare",
				RustResourceName:    "AwsRamResourceShare",
				GraphQlResourceName: "Aws_Ram_ResourceShare",
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
						RustSourcePropertyName: "resource_arns",
						RustReturnType:         "Vec<AwsRamResourceShareConnections_ResourceArns>",
						RustGenericType:        "Vec<AwsRamResourceShareConnections_ResourceArns>",
						TargetUnion: ResourceUnion{
							RustUnionName: "AwsRamResourceShareConnections_ResourceArns",
							RustResourceNames: map[string]bool{
								"AwsEc2PrefixList":                    true,
								"AwsRoute53ResolverFirewallRuleGroup": true,
								"Node":                                true,
							},
						},
					},
					{
						RustSourcePropertyName: "tags_value",
						RustReturnType:         "Vec<Node>",
						RustGenericType:        "Vec<Node>",
						TargetUnion:            ResourceUnion{},
					},
				},
			},
			expectedErrors: nil,
		},
	}

	for _, tc := range testCase {
		translated, errors := translateResource(
			tc.resourceName, testRelationships, cfnSchemaTest, nil)
		if len(errors) != len(tc.expectedErrors) {
			t.Errorf("For %s expected no errors but got %v", tc.description, errors)
		}
		if !compareTranslatedResources(translated, tc.expectedResource) {
			t.Errorf("For %s expected %v but got %v", tc.description, tc.expectedResource, translated)
		}
	}
}

func TestTranslateAllResources(t *testing.T) {
	mockCfnSchema := make(map[string]map[string]any) // not important for this test
	mockRelationships := AllRelationships{
		"AWS::Personalize::Dataset": Resource{
			PrimaryKeys:   nil,
			Relationships: nil,
		},
		"AWS::IoT1Click::Placement": {
			PrimaryKeys:   nil,
			Relationships: nil,
		},
		"AWS::EC2::PrefixList": {
			PrimaryKeys:   nil,
			Relationships: nil,
		},
	}
	mockTranslateResource := func(resourceName string, relationships AllRelationships, cfnSchema map[string]map[string]any, filter map[string]bool) (ResourceType, []error) {
		awsResource, err := NewAwsResourceName(resourceName)
		if err != nil {
			return ResourceType{}, []error{err}
		}
		return ResourceType{
			CfnResourceName:     awsResource.AsCfn(),
			RustResourceName:    awsResource.AsRust(),
			GraphQlResourceName: awsResource.AsGraphQl(),
			Properties:          nil,
			Relationships:       nil,
		}, nil
	}
	
	testCases := []struct {
		filter   map[string]bool
		expected []ResourceType
	}{
		{
			filter: nil,
			expected: []ResourceType{
				{
					CfnResourceName:     "AWS::EC2::PrefixList",
					RustResourceName:    "AwsEc2PrefixList",
					GraphQlResourceName: "Aws_Ec2_PrefixList",
					Properties:          nil,
					Relationships:       nil,
				},
				{
					CfnResourceName:     "AWS::IoT1Click::Placement",
					RustResourceName:    "AwsIoT1ClickPlacement",
					GraphQlResourceName: "Aws_IoT1_Click_Placement",
					Properties:          nil,
					Relationships:       nil,
				},
				{
					CfnResourceName:     "AWS::Personalize::Dataset",
					RustResourceName:    "AwsPersonalizeDataset",
					GraphQlResourceName: "Aws_Personalize_Dataset",
					Properties:          nil,
					Relationships:       nil,
				},
			},
		},
		{
			filter: map[string]bool{"AWS::IoT1Click::Placement": true, "AWS::EC2::PrefixList": true,},
			expected: []ResourceType{
				{
					CfnResourceName:     "AWS::EC2::PrefixList",
					RustResourceName:    "AwsEc2PrefixList",
					GraphQlResourceName: "Aws_Ec2_PrefixList",
					Properties:          nil,
					Relationships:       nil,
				},
				{
					CfnResourceName:     "AWS::IoT1Click::Placement",
					RustResourceName:    "AwsIoT1ClickPlacement",
					GraphQlResourceName: "Aws_IoT1_Click_Placement",
					Properties:          nil,
					Relationships:       nil,
				},
			},
		},
	}

	for _, tc := range testCases {
		translated, errors := translateAllResources(mockRelationships, mockCfnSchema, nil, mockTranslateResource)
		if len(errors) > 0 {
			t.Errorf("Expected no errors but got %v", errors)
		}
		if len(translated) != len(tc.expected) {
			t.Errorf("Expected %v but got %v", tc.expected, translated)
		}
		// we mostly test that the list order is the same
		for i, resource := range translated {
			if !compareTranslatedResources(resource, tc.expected[i]) {
				t.Errorf("Expected %v but got %v", tc.expected[i], resource)
			}
		}
	}
}

func compareTranslatedResources(left ResourceType, right ResourceType) bool {
	if left.CfnResourceName != right.CfnResourceName {
		return false
	}
	if left.RustResourceName != right.RustResourceName {
		return false
	}
	if left.GraphQlResourceName != right.GraphQlResourceName {
		return false
	}
	if !reflect.DeepEqual(left.Properties, right.Properties) {
		return false
	}
	if len(left.Relationships) != len(right.Relationships) {
		return false
	}
	leftRelDict := relationshipsToMap(left.Relationships)
	rightRelDict := relationshipsToMap(right.Relationships)
	return reflect.DeepEqual(leftRelDict, rightRelDict)
}

func relationshipsToMap(relationships []ResourceRelationship) map[string]ResourceRelationship {
	relDict := make(map[string]ResourceRelationship)
	for _, relationship := range relationships {
		relDict[relationship.RustSourcePropertyName] = ResourceRelationship{
			RustSourcePropertyName: relationship.RustSourcePropertyName,
			RustReturnType:         relationship.RustReturnType,
			RustGenericType:        relationship.RustGenericType,
			TargetUnion:            relationship.TargetUnion,
		}
	}
	return relDict
}
