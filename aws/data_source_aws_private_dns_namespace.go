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

			"filter": sdCustomFiltersSchema(), // TODO: create file for and define

			"tags": tagsSchemaComputed(),
		},
	}
}

func dataSourceAwsPrivateDnsNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).sdconn

	nsName := d.Get("name").(string)

	input := &servicediscovery.ListNamespacesInput{
		NamespaceNames: []*string{aws.String(nsName)},
	}

	// if id, ok := d.GetOk("id"); ok {
	// 	req.SubnetIds = []*string{aws.String(id.(string))}
	// }

	log.Printf("[DEBUG] Reading Namespace: %s", input)

	resp, err := conn.ListNamespaces(input)
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
	d.Set("arn", namespace.NamespaceArn)

	return nil
}
