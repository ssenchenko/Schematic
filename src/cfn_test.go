package schematic

import (
	"errors"
	"reflect"
	"testing"
)

var schema = Dict{
	PROPS: Dict{
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
		"NoOneOrAny": Dict{TYPE: STRING},
		"FutureRelationshipLink": Dict{
			REF: "#/definitions/FutureRelationshipLink",
		},
	},
	DEFINITIONS: Dict{
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
					"type" : "string",
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
		actualType, actualErr := ExtractRefTypeName(testCase.ref)

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
			expectedError: errors.New(`no type found for invalidType in {"FutureRelationshipLink":{"properties":{"KeyId":{"anyOf":[{"relationshipRef":{"propertyPath":"/properties/Arn","typeName":"AWS::KMS::Key"}},{"relationshipRef":{"propertyPath":"/properties/KeyId","typeName":"AWS::KMS::Key"}}],"type":"string"}},"type":"object"},"RefToString":{"$ref":"#/definitions/SimpleString"},"RefWithAny":{"anyOf":[{"type":"integer"},{"type":"number"}]},"SimpleString":{"type":"string"}}`),
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
		fragment      Dict
		expectedTypes []Dict
		expectedError error
	}{
		{
			fragment: schema[PROPS].(Dict)["AnyType"].(Dict),
			expectedTypes: []Dict{
				{TYPE: STRING},
				{TYPE: OBJECT},
			},
			expectedError: nil,
		},
		{
			fragment: schema[PROPS].(Dict)["AnyReferences"].(Dict),
			expectedTypes: []Dict{
				{TYPE: STRING},
				{TYPE: STRING},
				{TYPE: INT},
				{TYPE: NUMBER},
			},
			expectedError: nil,
		},
		{
			fragment: schema[PROPS].(Dict)["Nested"].(Dict),
			expectedTypes: []Dict{
				{TYPE: STRING},
				{TYPE: BOOL},
				{TYPE: OBJECT},
			},
			expectedError: nil,
		},
		{
			fragment:      schema[PROPS].(Dict)["NoOneOrAny"].(Dict),
			expectedTypes: nil,
			expectedError: errors.New(`no anyOf or oneOf found {"type":"string"}`),
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
