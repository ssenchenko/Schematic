package schematic

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

var schema = Dict{
	PROPS: Dict{
		"NotDefined": Dict{
			ONE_OF: []Dict{
				{TYPE: STRING},
				{
					TYPE: ARRAY,
					ITEMS: Dict{
						TYPE: STRING,
					},
				},
			},
		},
		"ConfusingType": Dict{
			ONE_OF: []Dict{
				{REF: "#/definitions/SomeObject"},
				{REF: "#/definitions/AnotherObject"},
			},
		},
		"MoreConfusingType": Dict{
			ONE_OF: []Dict{
				{REF: "#/definitions/AnotherObject"},
				{REF: "#/definitions/SimpleString"},
			},
		},
		"SimpleArray": Dict{
			TYPE: ARRAY,
			ITEMS: Dict{
				TYPE: STRING,
			},
		},
		"ArrayWithRef": Dict{
			TYPE: ARRAY,
			ITEMS: Dict{
				REF: "#/definitions/RefToString",
			},
		},
		"ArrayWithObject": Dict{
			TYPE: ARRAY,
			ITEMS: Dict{
				REF: "#/definitions/SomeObject",
			},
		},
		"UntypedObject": Dict{
			TYPE: OBJECT,
		},
		"UntypedWeirdObject": Dict{
			TYPE: []string{OBJECT, STRING},
		},
		"ArrayWithWeirdType": Dict{
			TYPE: []string{ARRAY, STRING},
			ITEMS: Dict{
				TYPE: NUMBER,
			},
		},
		"DeepArrayWithWeirdType": Dict{
			TYPE: []string{ARRAY, STRING},
			ITEMS: Dict{
				REF: "#/definitions/AnotherObject",
			},
		},
		"ArrayWithOne": Dict{
			TYPE: ARRAY,
			ITEMS: Dict{
				ONE_OF: []Dict{
					{REF: "#/definitions/SomeObject"},
					{REF: "#/definitions/SimpleString"},
				},
			},
		},
		"AnyType": Dict{
			ANY_OF: []Dict{
				{TYPE: STRING},
				{TYPE: OBJECT},
			},
		},
		"AnyReferences": Dict{
			ANY_OF: []Dict{
				{REF: "#/definitions/SimpleString"},
				{REF: "#/definitions/RefToString"},
				{REF: "#/definitions/RefWithAny"},
			},
		},
		"Nested": Dict{
			ONE_OF: []Dict{
				{TYPE: STRING},
				{
					ANY_OF: []Dict{
						{TYPE: BOOL},
						{TYPE: OBJECT},
					},
				},
			},
		},
		"JustString": Dict{TYPE: STRING},
		"FutureRelationshipLink": Dict{
			REF: "#/definitions/FutureRelationshipLink",
		},
		"ShouldNotHappen": Dict{
			"Unexpected": Dict{
				"Stuff": "WTF",
			},
		},
	},
	DEFINITIONS: Dict{
		"SomeObject": Dict{
			TYPE: OBJECT,
			PROPS: Dict{
				"SomeProperty": Dict{
					TYPE: INT,
				},
			},
		},
		"AnotherObject": Dict{
			TYPE: OBJECT,
			PROPS: Dict{
				"ArrayProperty": Dict{
					TYPE: ARRAY,
					ITEMS: Dict{
						REF: "#/definitions/SomeObject",
					},
				},
			},
		},
		"SimpleString": Dict{TYPE: STRING},
		"RefToString":  Dict{REF: "#/definitions/SimpleString"},
		"RefWithAny": Dict{
			ANY_OF: []Dict{
				{TYPE: INT},
				{TYPE: NUMBER},
			},
		},
		"FutureRelationshipLink": Dict{
			"type": "object",
			"properties": Dict{
				"KeyId": Dict{
					"type": "string",
					"anyOf": []Dict{
						{
							"relationshipRef": Dict{
								"typeName":     "AWS::KMS::Key",
								"propertyPath": "/properties/Arn",
							},
						},
						{
							"relationshipRef": Dict{
								"typeName":     "AWS::KMS::Key",
								"propertyPath": "/properties/KeyId",
							},
						},
					},
				},
			},
		},
	},
}

func TestExtractRefTypeName(t *testing.T) {
	testCases := []struct {
		ref          string
		expectedType string
		expectedErr  error
	}{
		{"#/definitions/AnotherType", "AnotherType", nil},
		{"", "", errors.New("unexpected ref: ")},
		{"#/definitions/", "", errors.New("unexpected ref: #/definitions/")},                     // Invalid ref format
		{"#/invalid/Type", "", errors.New("unexpected ref: #/invalid/Type")},                     // Invalid ref format
		{"#/definitions/Type/Extra", "", errors.New("unexpected ref: #/definitions/Type/Extra")}, // Invalid ref format
	}

	for _, testCase := range testCases {
		actualType, actualErr := extractRefTypeName(testCase.ref)

		if actualType != testCase.expectedType {
			t.Errorf("For ref %s, expected type %s, but got %s", testCase.ref, testCase.expectedType, actualType)
		}

		if (actualErr == nil && testCase.expectedErr != nil) || (actualErr != nil && testCase.expectedErr == nil) || (actualErr != nil && actualErr.Error() != testCase.expectedErr.Error()) {
			t.Errorf("For ref %s, expected error '%v', but got '%v'", testCase.ref, testCase.expectedErr, actualErr)
		}
	}
}

func TestResolveRef(t *testing.T) {
	testCases := []struct {
		ref           string
		expectedTypes []Dict
		expectedError error
	}{
		{
			ref: "#/definitions/SimpleString",
			expectedTypes: []Dict{
				{TYPE: STRING},
			},
			expectedError: nil,
		},
		{
			ref: "#/definitions/RefToString",
			expectedTypes: []Dict{
				{TYPE: STRING},
			},
			expectedError: nil,
		},
		{
			ref:           "#/definitions/invalidType",
			expectedTypes: nil,
			expectedError: fmt.Errorf("no type found for invalidType in %s", schema[DEFINITIONS].(Dict)),
		},
	}

	for _, testCase := range testCases {
		actualTypes, actualError := ResolveRef(testCase.ref, schema)

		if !reflect.DeepEqual(actualTypes, testCase.expectedTypes) {
			t.Errorf("For ref %s, expected types %v, but got %v", testCase.ref, testCase.expectedTypes, actualTypes)
		}

		if (actualError == nil && testCase.expectedError != nil) || (actualError != nil && testCase.expectedError == nil) || (actualError != nil && actualError.Error() != testCase.expectedError.Error()) {
			t.Errorf("For ref %s, expected error '%v', but got '%v'", testCase.ref, testCase.expectedError, actualError)
		}
	}
}

func TestResolveAnyOneOf(t *testing.T) {
	testCases := []struct {
		fragment      []Dict
		expectedTypes []Dict
		expectedError error
	}{
		{
			fragment: schema[PROPS].(Dict)["AnyType"].(Dict)[ANY_OF].([]Dict),
			expectedTypes: []Dict{
				{TYPE: STRING},
				{TYPE: OBJECT},
			},
			expectedError: nil,
		},
		{
			fragment: schema[PROPS].(Dict)["AnyReferences"].(Dict)[ANY_OF].([]Dict),
			expectedTypes: []Dict{
				{TYPE: STRING},
				{TYPE: STRING},
				{TYPE: INT},
				{TYPE: NUMBER},
			},
			expectedError: nil,
		},
		{
			fragment: schema[PROPS].(Dict)["Nested"].(Dict)[ONE_OF].([]Dict),
			expectedTypes: []Dict{
				{TYPE: STRING},
				{TYPE: BOOL},
				{TYPE: OBJECT},
			},
			expectedError: nil,
		},
	}

	for _, testCase := range testCases {
		actualTypes, actualError := ResolveAnyOneOf(testCase.fragment, schema)

		if !reflect.DeepEqual(actualTypes, testCase.expectedTypes) {
			t.Errorf("For ref %s, expected types %v, but got %v", testCase.fragment, testCase.expectedTypes, actualTypes)
		}

		if (actualError == nil && testCase.expectedError != nil) || (actualError != nil && testCase.expectedError == nil) || (actualError != nil && actualError.Error() != testCase.expectedError.Error()) {
			t.Errorf("For ref %s, expected error '%v', but got '%v'", testCase.fragment, testCase.expectedError, actualError)
		}
	}
}

func TestOneStepAtATime(t *testing.T) {
	testCases := []struct {
		propertyName string
		fragment     Dict
		expected     []Dict
		expectedErr  error
	}{
		{
			propertyName: "JustString",
			fragment:     schema,
			expected: []Dict{
				{TYPE: STRING},
			},
			expectedErr: nil,
		},
		{
			propertyName: "AnyReferences",
			fragment:     schema,
			expected: []Dict{
				{TYPE: STRING},
				{TYPE: STRING},
				{TYPE: INT},
				{TYPE: NUMBER},
			},
			expectedErr: nil,
		},
		{
			propertyName: "SimpleArray",
			fragment:     schema,
			expected: []Dict{
				{
					TYPE: ARRAY,
					ITEMS: Dict{
						TYPE: STRING,
					},
				},
			},
			expectedErr: nil,
		},
		{
			propertyName: "SomeProperty",
			fragment:     schema[PROPS].(Dict)["ArrayWithObject"].(Dict),
			expected: []Dict{
				{
					TYPE: INT,
				},
			},
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		actual, err := OneStepAtATime(testCase.propertyName, testCase.fragment, schema)

		if err != nil {
			if testCase.expectedErr == nil {
				t.Errorf("For property %s, expected no error, but got '%v'", testCase.propertyName, err)
			} else if err.Error() != testCase.expectedErr.Error() {
				t.Errorf("For property %s, expected error '%v', but got '%v'", testCase.propertyName, testCase.expectedErr, err)
			}
		} else if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("For property %s, expected %v, but got %v", testCase.propertyName, testCase.expected, actual)
		}
	}
}

func TestFollowPath(t *testing.T) {
	testCases := []struct {
		path        string
		expected    map[TypeBits]bool
		expectedErr error
	}{
		{
			path: "JustString",
			expected: map[TypeBits]bool{
				OBJECTB | STRINGB: true,
			},
			expectedErr: nil,
		},
		{
			path: "SimpleArray",
			expected: map[TypeBits]bool{
				OBJECTB | ARRAYB: true,
			},
			expectedErr: nil,
		},
		{
			path: "ArrayWithRef",
			expected: map[TypeBits]bool{
				OBJECTB | ARRAYB: true,
			},
			expectedErr: nil,
		},
		{
			path: "AnyReferences",
			expected: map[TypeBits]bool{
				OBJECTB | STRINGB: true,
				OBJECTB | INTB:    true,
				OBJECTB | NUMBERB: true,
			},
			expectedErr: nil,
		},
		{
			path: "ArrayWithObject/SomeProperty",
			expected: map[TypeBits]bool{
				OBJECTB | ARRAYB | INTB: true,
			},
			expectedErr: nil,
		},
		{
			path: "ArrayWithOne/SomeProperty",
			expected: map[TypeBits]bool{
				OBJECTB | ARRAYB | INTB: true,
			},
			expectedErr: nil,
		},
		{
			path: "ConfusingType/SomeProperty",
			expected: map[TypeBits]bool{
				OBJECTB | INTB: true,
			},
			expectedErr: nil,
		},
		{
			path: "ConfusingType/ArrayProperty/SomeProperty",
			expected: map[TypeBits]bool{
				OBJECTB | ARRAYB | INTB: true,
			},
			expectedErr: nil,
		},
		{
			path: "MoreConfusingType",
			expected: map[TypeBits]bool{
				OBJECTB | STRINGB: true,
				OBJECTB:           true,
			},
			expectedErr: nil,
		},
		{
			path: "MoreConfusingType/ArrayProperty/SomeProperty",
			expected: map[TypeBits]bool{
				OBJECTB | ARRAYB | INTB: true,
			},
			expectedErr: nil,
		},
		{
			path:        "NonexistantProperty",
			expected:    map[TypeBits]bool{},
			expectedErr: nil,
		},
		{
			path: "UntypedObject",
			expected: map[TypeBits]bool{
				OBJECTB: true,
			},
			expectedErr: nil,
		},
		{
			path: "UntypedWeirdObject",
			expected: map[TypeBits]bool{
				OBJECTB | STRINGB: true,
			},
			expectedErr: nil,
		},
		{
			path: "ArrayWithWeirdType",
			expected: map[TypeBits]bool{
				OBJECTB | ARRAYB | STRINGB: true,
			},
			expectedErr: nil,
		},
		{
			path: "DeepArrayWithWeirdType/ArrayProperty/SomeProperty",
			expected: map[TypeBits]bool{
				OBJECTB | ARRAYB | STRINGB | INTB: true,
			},
			expectedErr: nil,
		},
		{
			path:        "ShouldNotHappen/Unexpected",
			expected:    nil,
			expectedErr: fmt.Errorf("error in path ShouldNotHappen/Unexpected no idea how to handle %s", schema[PROPS].(Dict)["ShouldNotHappen"]),
		},
	}

	for _, testCase := range testCases {
		actual, err := FollowPath(testCase.path, schema)

		if err != nil {
			if testCase.expectedErr == nil {
				t.Errorf("For path %s, expected no error, but got '%v'", testCase.path, err)
			} else if err.Error() != testCase.expectedErr.Error() {
				t.Errorf("For path %s, expected error '%v', but got '%v'", testCase.path, testCase.expectedErr, err)
			}
		} else if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("For path %s, expected %v, but got %v", testCase.path, testCase.expected, actual)
		}
	}
}

func TestIsArray(t *testing.T) {
	var resources = map[string]Dict{
		"AWS::Test::Test": schema,
	}
	testCases := []struct {
		path        string
		expected    bool
		expectedErr error
	}{
		{
			path:        "JustString",
			expected:    false,
			expectedErr: nil,
		},
		{
			path:        "SimpleArray",
			expected:    true,
			expectedErr: nil,
		},
		{
			path:        "NotDefined",
			expected:    false,
			expectedErr: fmt.Errorf(
				"what should I do with it? %s in %s has branches; some of them end with array and some - don't", "NotDefined", "AWS::Test::Test"),
		},
		{
			path:        "NonexistantProperty",
			expected:    false,
			expectedErr: fmt.Errorf("no type found for %s in %s", "NonexistantProperty", "AWS::Test::Test"),
		},
		{
			path:        "ShouldNotHappen/Unexpected",
			expected:    false,
			expectedErr: fmt.Errorf("error in path ShouldNotHappen/Unexpected no idea how to handle %s", schema[PROPS].(Dict)["ShouldNotHappen"]),
		},
	}

	for _, testCase := range testCases {
		actual, err := IsArray(testCase.path, "AWS::Test::Test", resources)

		if err != nil {
			if testCase.expectedErr == nil {
				t.Errorf("For path %s, expected no error, but got '%v'", testCase.path, err)
			} else if err.Error() != testCase.expectedErr.Error() {
				t.Errorf("For path %s, expected error '%v', but got '%v'", testCase.path, testCase.expectedErr, err)
			}
		} else if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("For path %s, expected %v, but got %v", testCase.path, testCase.expected, actual)
		}
	}
}
