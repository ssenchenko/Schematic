package translator

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
)

type AwsResourceName struct {
	Partition string
	Service   string
	Resource  string
}

func NewAwsResourceName(resourceName string) (name AwsResourceName, err error) {
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
	partition, service, resource := r.asPascalCase()
	return fmt.Sprintf("%s%s%s", partition, service, resource)
}

func (r AwsResourceName) AsGraphQl() string {
	partition, service, resource := r.asPascalCase()
	return fmt.Sprintf("%s_%s_%s", partition, service, resource)
}

func (r AwsResourceName) asPascalCase() (
	partition string,
	service string,
	resource string,
) {
	partition = strcase.ToCamel(r.Partition)
	service = strcase.ToCamel(r.Service)
	resource = strcase.ToCamel(r.Resource)
	return partition, service, resource
}
