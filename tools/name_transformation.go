package tools

import (
	"fmt"
	"regexp"
	"strings"
)

type StringCase int

const (
	UPPER           StringCase = 1
	PascalCase      StringCase = 2
	MIXEDPascalCase StringCase = 3
	Other           StringCase = 100
)

var (
	upperCaseRe       = regexp.MustCompile("^[A-Z]+(?:[0-9])*$")
	pascalCaseRe      = regexp.MustCompile("^[A-Z][a-z]+(?:[A-Z][a-z]+)*(?:[0-9])*$")
	mixedPascalCaseRe = regexp.MustCompile("^([A-Z]+)([A-Z][a-z]+(?:[A-Z][a-z]+)*[0-9]*)$")
)

type AwsResourceName struct {
	Partition string
	Service   string
	Resource  string
}

func NewAwsResourceName(resourceName string) (AwsResourceName, error) {
	name := AwsResourceName{}
	parts := strings.Split(resourceName, "::")
	if len(parts) != 3 {
		return name, fmt.Errorf("invalid resource name: %s", resourceName)
	}
	name.Partition = parts[0]
	name.Service = parts[1]
	name.Resource = parts[2]
	return name, nil
}

func (r AwsResourceName) AsCfn() string {
	return fmt.Sprintf("%s::%s::%s", r.Partition, r.Service, r.Resource)
}

func (r AwsResourceName) AsRust() string {
	partition := ToPascalCase(r.Partition)
	service := ToPascalCase(r.Service)
	resource := ToPascalCase(r.Resource)
	return fmt.Sprintf("%s%s%s", partition, service, resource)
}

func (r AwsResourceName) AsGraphQl() string {
	partition := ToPascalCase(r.Partition)
	service := ToPascalCase(r.Service)
	resource := ToPascalCase(r.Resource)
	return fmt.Sprintf("%s_%s_%s", partition, service, resource)
}

func GetStringCase(str string) StringCase {
	if upperCaseRe.MatchString(str) {
		return UPPER
	}
	if pascalCaseRe.MatchString(str) {
		return PascalCase
	}
	if mixedPascalCaseRe.MatchString(str) {
		return MIXEDPascalCase
	}
	return Other
}

func PascalCaseToSnakeCase(str string) string {
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func ToPascalCase(str string) string {
	switch GetStringCase(str) {
	case UPPER:
		return UpperToPascalCase(str)
	case MIXEDPascalCase:
		return MixedPascalToPascalCase(str)
	case PascalCase:
		return str
	default:
		return str // maybe raise an error?
	}
}

func MixedPascalToPascalCase(str string) string {
	matches := mixedPascalCaseRe.FindStringSubmatch(str)
	return fmt.Sprintf("%s%s", UpperToPascalCase(matches[1]), matches[2])
}

func UpperToPascalCase(str string) string {
	return strings.Title(strings.ToLower(str))
}
