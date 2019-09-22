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

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"filter": sdCustomFiltersSchema(), // TODO: create file for and define

			"tags": tagsSchemaComputed(),
		},
	}
}

func dataSourceAwsPrivateDnsNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).sdconn

	req := &servicediscovery.ListNamespaces{}
}
