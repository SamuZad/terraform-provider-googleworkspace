package googleworkspace

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDynamicGroup() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(resourceDynamicGroup().Schema)
	addExactlyOneOfFieldsToSchema(dsSchema, "id", "email")

	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Dynamic group data source in the Terraform Googleworkspace provider. Dynamic Group resides under the " +
			"`https://www.googleapis.com/auth/cloud-identity.groups` client scope.",

		ReadContext: dataSourceDynamicGroupRead,

		Schema: dsSchema,
	}
}

func dataSourceDynamicGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.Get("id") != "" {
		d.SetId(d.Get("id").(string))
	} else {
		var diags diag.Diagnostics

		// use the meta value to retrieve your client from the provider configure method
		client := meta.(*apiClient)

		cloudIdentityService, diags := client.NewCloudIdentityService()
		if diags.HasError() {
			return diags
		}

		groupsService, diags := GetDynamicGroupsService(cloudIdentityService)
		if diags.HasError() {
			return diags
		}

		lookupGroupResponse, err := groupsService.Lookup().GroupKeyId(d.Get("email").(string)).Do()
		if err != nil {
			return diag.FromErr(err)
		}

		groupName := lookupGroupResponse.Name

		group, err := groupsService.Get(groupName).Do()
		if err != nil {
			return diag.FromErr(err)
		}

		if group == nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("No dynamic group was returned for %s.", d.Get("email").(string)),
			})

			return diags
		}

		d.SetId(group.Name)
	}

	return resourceDynamicGroupRead(ctx, d, meta)
}
