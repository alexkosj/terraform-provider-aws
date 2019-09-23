package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicediscovery"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAwsPrivateDnsNamespace() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAwsPrivateDnsNamespaceRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"arn": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"filter": sdCustomFiltersSchema(),
		},
	}
}

func dataSourceAwsPrivateDnsNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).sdconn

	req := &servicediscovery.ListNamespaces()

	if id, ok := d.GetOk("id"); ok {
		req.Namespaces = []*string{aws.String(id.(string))}
	}

	filters := map[string]string{
		"name":   d.Get("name").(string)
	}

	req.Filters = buildSdAttributeFilterList(filters)
	if len(req.Filters) == 0 {
		// Don't send an empty filters list; the Service Discovery API won't accept it.
		req.Filters = nil
	}

	log.Printf("[DEBUG] Reading Namespace: %s", input)

	resp, err := conn.ListNamespaces(req)
	if err != nil {
		return fmt.Errorf("Error retrieving namespace: %s", err)
	}
	if resp == nil || len(resp.Namespaces) == 0 {
		return fmt.Errorf("no matching namespace found")
	}
	if len(resp.Namespaces) > 1 {
		return fmt.Errorf("multiple namespaces matched; use additional constraints to reduce matches to a single namespace")
	}

	namespace := resp.Namespaces[0]

	d.SetId(*namespace.NamespaceId)
	d.Set("name", namespace.NamespaceName)
	d.Set("arn", namespace.NamespaceArn)

	return nil
}
