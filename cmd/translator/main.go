package main

import (
	"log"
	"os"

	"ssenchenko/schematic/translator"
)

const (
	// IaScopeOnly flag indicates that only resources in IaScopeResources should be processed
	IaScopeOnly = true
)

var (
	// IaScopeResources contains list of resources for IA
	// it's a map (set) to void adding the same resource twice by mistake
	IaScopeResources = map[string]bool{
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

	// AllRelationshipsOverrides allows to override names in relationship file.
	// Some property names seems to be wrong and has to be overridden.
	AllRelationshipsOverrides = map[string]map[string]string{
		"AWS::EC2::Instance": {"VolumeAttachments": "Volumes"},
	}
)

func main() {
	relationshipsFile := "data/all-schema-combined.json"
	if len(os.Args) > 1 {
		relationshipsFile = os.Args[1]
	}
	allRelationships := Must(translator.LoadAllRelationships(relationshipsFile))
	allRelationships.ApplyOverrides(AllRelationshipsOverrides)

	cfnSchemaDir := "data/cfn"
	if len(os.Args) > 2 {
		cfnSchemaDir = os.Args[2]
	}
	cfnSchemaCombined := Must(translator.LoadCfnSchemaCombined(cfnSchemaDir, nil))
	var filter map[string]bool
	if IaScopeOnly {
		filter = IaScopeResources
	} else {
		filter = make(map[string]bool)
	}
	rustModel, errors := translator.Translate(allRelationships, cfnSchemaCombined, filter)

	if len(errors) > 0 {
		log.Println("Translation Errors:")
		for _, err := range errors {
			log.Println(" >>", err.Error())
		}
	}

	bytes := Must(translator.HydrateTemplates(rustModel))
	err := os.WriteFile("configuration/model.rs", bytes.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
