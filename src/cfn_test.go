package schematic

import (
	"errors"
	"reflect"
	"testing"
)

func TestExtractRefTypeName(t *testing.T) {
	testCases := []struct {
		ref          string
		expectedType string
		expectedErr  error
	}{
		{"#/definitions/Type", "Type", nil},
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
	// Mock JSON schema fragments
	fragment := map[string]any{
		DEFINITIONS: map[string]any{
			"type1": map[string]any{
				TYPE: STRING,
			},
			"type2": map[string]any{
				REF: "#/definitions/type1",
			},
		},
	}

	testCases := []struct {
		ref           string
		expectedTypes []map[string]any
		expectedError error
	}{
		{
			ref: "#/definitions/type1",
			expectedTypes: []map[string]any{
				{"type": "string"},
			},
			expectedError: nil,
		},
		{
			ref: "#/definitions/type2",
			expectedTypes: []map[string]any{
				{"type": "string"},
			},
			expectedError: nil,
		},
		{
			ref:           "#/definitions/invalidType",
			expectedTypes: nil,
			expectedError: errors.New(`no type found for invalidType in
{
  "definitions": {
    "type1": {
      "type": "string"
    },
    "type2": {
      "$ref": "#/definitions/type1"
    }
  }
}`,
			),
		},
	}

	for _, testCase := range testCases {
		actualTypes, actualError := ResolveRef(testCase.ref, fragment)

		if !reflect.DeepEqual(actualTypes, testCase.expectedTypes) {
			t.Errorf("For ref %s, expected types %v, but got %v", testCase.ref, testCase.expectedTypes, actualTypes)
		}

		if (actualError == nil && testCase.expectedError != nil) || (actualError != nil && testCase.expectedError == nil) || (actualError != nil && actualError.Error() != testCase.expectedError.Error()) {
			t.Errorf("For ref %s, expected error '%v', but got '%v'", testCase.ref, testCase.expectedError, actualError)
		}
	}
}
