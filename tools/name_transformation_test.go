package tools

import "testing"

func TestGetStringCase(t *testing.T) {
	testCases := []struct {
		input    string
		expected StringCase
	}{
		{"AWS", UPPER},
		{"AWS1", UPPER},
		{"A123", UPPER},
		{"AWSService", MIXEDPascalCase},
		{"AWSService12", MIXEDPascalCase},
		{"AWSServiceResource1", MIXEDPascalCase},
		{"AwsService", PascalCase},
		{"AwsService2", PascalCase},
		{"AwsService123", PascalCase},
	}

	for _, tc := range testCases {
		actual := GetStringCase(tc.input)
		if actual != tc.expected {
			t.Errorf("Expected: %v, Got: %v", tc.expected, actual)
		}
	}
}

func TestPascalCaseToSnakeCase(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"AwsService", "aws_service"},
		{"KMS", "kms"},
		{"AWSServiceMixed", "aws_service_mixed"},
		{"AwsServiceResource", "aws_service_resource"},
		{"AwsServiceResource12", "aws_service_resource12"},
	}

	for _, tc := range testCases {
		actual := PascalCaseToSnakeCase(tc.input)
		if actual != tc.expected {
			t.Errorf("Expected: %v, Got: %v", tc.expected, actual)
		}
	}
}

func TestMixedPascalToPascalCase(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"AWSService", "AwsService"},
		{"AWSService12", "AwsService12"},
		{"AWSServiceProperty12", "AwsServiceProperty12"},
	}

	for _, tc := range testCases {
		actual := MixedPascalToPascalCase(tc.input)
		if actual != tc.expected {
			t.Errorf("Expected: %v, Got: %v", tc.expected, actual)
		}
	}
}
