package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/the-urge-tech/terraform-provider-contentful/pkg/contentful"
)

const warningMessage = "[DO NOT EDIT: Managed by Terraform] "

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"cma_token": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("CONTENTFUL_MANAGEMENT_TOKEN", nil),
					Description: "The Contentful Management API token",
				},
				"organization_id": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("CONTENTFUL_ORGANIZATION_ID", nil),
					Description: "The organization ID",
				},
				"env": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("CONTENTFUL_ENVIRONMENT", nil),
					Description: "The target environment id",
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				"contentful_contenttype": resourceContentfulContentType(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		c := contentful.NewClient(d.Get("cma_token").(string), d.Get("organization_id").(string), d.Get("env").(string))
		return c, nil
	}
}
