package aws

import (
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicediscovery"
	"github.com/hashicorp/terraform/helper/schema"
)

// TODO: document
func buildSdAttributeFilterList(attrs map[string]string) []*servicediscovery.Filter {
	var filters []*ec2.Filter

	// sort the filters by name to make the output deterministic
	var names []string
	for filterName := range attrs {
		names = append(names, filterName)
	}

	sort.Strings(names)

	for _, filterName := range names {
		value := attrs[filterName]
		if value == "" {
			continue
		}

		filters = append(filters, &servicediscovery.Filter{
			Name:   aws.String(filterName),
			Values: []*string{aws.String(value)},
		})
	}

	return filters
}

// TODO: document
func buildSdTagFilterList(tags []*servicediscovery.Tag) []*servicediscovery.Filter {
	filters := make([]*servicediscovery.Filter, len(tags))

	for i, tag := range tags {
		filters[i] = &servicediscovery.Filter{
			Name:   aws.String(fmt.Sprintf("tag:%s", *tag.Key)),
			Values: []*string{tag.Value},
		}
	}

	return filters
}

// TODO: document
func sDCustomFiltersSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"values": {
					Type:     schema.TypeSet,
					Required: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

// TODO: document
func buildSDCustomFilterList(filterSet *schema.Set) []*servicediscovery.Filter {
	if filterSet == nil {
		return []*servicediscovery.Filter{}
	}

	customFilters := filterSet.List()
	filters := make([]*servicediscovery.Filter, len(customFilters))

	for filterIdx, customFilterI := range customFilters {
		customFilterMapI := customFilterI.(map[string]interface{})
		name := customFilterMapI["name"].(string)
		valuesI := customFilterMapI["values"].(*schema.Set).List()
		values := make([]*string, len(valuesI))
		for valueIdx, valueI := range valuesI {
			values[valueIdx] = aws.String(valueI.(string))
		}

		filters[filterIdx] = &servicediscovery.Filter{
			Name:   &name,
			Values: values,
		}
	}

	return filters
}
