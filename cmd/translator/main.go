package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"ssenchenko/schematic/cfn"
	"ssenchenko/schematic/translator"
)

const (
	// A flag to indicate that only resources in IA_SCOPE_RESOURCES should be processed
	IA_SCOPE_ONLY = true
)

var (
	IA_SCOPE_RESOURCES = map[string]bool{
		"AWS::EC2::Instance":         true,
		"AWS::EC2::VolumeAttachment": true,
		"AWS::EC2::Volume":           true,
		"AWS::EC2::VPC":              true,
		"AWS::EC2::SecurityGroup":    true,
		"AWS::IAM::InstanceProfile":  true,
		"AWS::IAM::Role":             true,
		"AWS::IAM::Policy":           true,
		"AWS::CloudWatch::Alarm":     true,
		"AWS::S3::Bucket":            true,
	}

	// some property names seems to be wrong in relationship schema file
	// and has to be overridden
	ALL_RELATIONSHIPS_OVERRIDES = map[string]map[string]string{
		"AWS::EC2::Instance": {"VolumeAttachments": "Volumes"},
	}
)

func main() {
	relationshipsFile := "data/all-schema-combined.json"
	if len(os.Args) > 1 {
		relationshipsFile = os.Args[1]
	}
	relationshipsFile = GetFullPath(relationshipsFile)
	content := Must(os.ReadFile(relationshipsFile))

	var allRelationships translator.AllRelationships
	err := json.Unmarshal(content, &allRelationships)
	if err != nil {
		panic(err)
	}
	allRelationships.ApplyOverrides(ALL_RELATIONSHIPS_OVERRIDES)

	cfnSchemaDir := "data/cfn"
	if len(os.Args) > 2 {
		cfnSchemaDir = os.Args[2]
	}
	cfnSchemaDir = GetFullPath(cfnSchemaDir)
	cfnFiles := Must(os.ReadDir(cfnSchemaDir))
	cfnJsonSchema := make(cfn.Dict)
	cfnSchemaCombined := make(map[string]cfn.Dict)
	for _, file := range cfnFiles {
		content = Must(os.ReadFile(GetFullPath(file.Name())))
		err := json.Unmarshal(content, &cfnJsonSchema)
		if err != nil {
			panic(err)
		}
		cfnSchemaCombined[cfnJsonSchema["typeName"].(string)] = cfnJsonSchema
	}
}

// Transform relative path to an absolute one using current working directory as a base
func GetFullPath(path string) string {
	if !filepath.IsAbs(path) {
		wd := Must(os.Getwd())
		return filepath.Join(wd, path)
	}
	return path
}

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
