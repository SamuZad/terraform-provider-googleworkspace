package googleworkspace

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"google.golang.org/api/cloudidentity/v1"
)

type Query struct {
	Query        string `json:"query"`
	ResourceType string `json:"resourceType"`
}

type DynamicGroupMeta struct {
	Queries []Query `json:"queries"`
}

type GroupKey struct {
	ID string `json:"id"`
}

type Group struct {
	Type                 string            `json:"@type"`
	CreateTime           string            `json:"createTime"`
	Description          string            `json:"description"`
	DisplayName          string            `json:"displayName"`
	DynamicGroupMetadata DynamicGroupMeta  `json:"dynamicGroupMetadata"`
	GroupKey             GroupKey          `json:"groupKey"`
	Labels               map[string]string `json:"labels"`
	Name                 string            `json:"name"`
	Parent               string            `json:"parent"`
	UpdateTime           string            `json:"updateTime"`
}

func resourceDynamicGroup() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Dynamic Group resource manages Google Workspace Groups with Dynamic memberships. Dynamic Group resides under the " +
			"`https://www.googleapis.com/auth/cloud-identity.groups` client scope.",

		CustomizeDiff: resourceExampleCustomizeDiff,
		CreateContext: resourceDynamicGroupCreate,
		ReadContext:   resourceDynamicGroupRead,
		UpdateContext: resourceDynamicGroupUpdate,
		DeleteContext: resourceDynamicGroupDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The unique ID of a group. A group id can be used as a group request URI's groupKey.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"email": {
				Description: "The group's email address. If your account has multiple domains," +
					"select the appropriate domain for the email address. The email must be unique.",
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Description: "The group's display name.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"query": {
				Description: "The dynamic group query.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "An extended description to help users determine the purpose of a group." +
					"For example, you can include information about who should join the group," +
					"the types of messages to send to the group, links to FAQs about the group, or related groups.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(0, 4096)),
			},
			"labels": {
				Description: "One or more label entries that apply to the Group. Currently supported labels contain a key with an empty value.",
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Default:     map[string]interface{}{"cloudidentity.googleapis.com/groups.discussion_forum": ""},
			},
		},
	}
}

func resourceDynamicGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// use the meta value to retrieve your client from the provider configure method
	client := meta.(*apiClient)

	email := d.Get("email").(string)
	log.Printf("[DEBUG] Creating Dynamic Group %q: %#v", email, email)

	cloudIdentityService, diags := client.NewCloudIdentityService()
	if diags.HasError() {
		return diags
	}

	groupsService, diags := GetDynamicGroupsService(cloudIdentityService)
	if diags.HasError() {
		return diags
	}

	labels, err := convertInterfaceMapToStringMap(d.Get("labels").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	groupObj := cloudidentity.Group{
		GroupKey: &cloudidentity.EntityKey{
			Id: d.Get("email").(string),
		},
		Parent:      "customerId/" + client.Customer,
		DisplayName: d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      labels,
		DynamicGroupMetadata: &cloudidentity.DynamicGroupMetadata{
			Queries: make([]*cloudidentity.DynamicGroupQuery, 0),
		},
	}

	query := d.Get("query").(string)

	groupObj.DynamicGroupMetadata.Queries = append(groupObj.DynamicGroupMetadata.Queries, &cloudidentity.DynamicGroupQuery{
		ResourceType: "USER",
		Query:        query,
	})

	group, err := groupsService.Create(&groupObj).InitialGroupConfig("EMPTY").Do()
	if err != nil {
		return diag.FromErr(err)
	}

	var groupResponse Group

	err = json.Unmarshal(group.Response, &groupResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(groupResponse.Name)

	log.Printf("[DEBUG] Finished creating Dynamic Group %q: %#v", d.Id(), email)

	return resourceDynamicGroupRead(ctx, d, meta)
}

func resourceDynamicGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	group, err := groupsService.Get(d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, d.Get("email").(string))
	}

	d.Set("name", group.DisplayName)
	d.Set("email", group.GroupKey.Id)
	d.Set("description", group.Description)
	d.Set("query", group.DynamicGroupMetadata.Queries[0].Query)
	d.Set("labels", group.Labels)

	d.SetId(group.Name)

	return diags
}

func resourceDynamicGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// use the meta value to retrieve your client from the provider configure method
	client := meta.(*apiClient)

	email := d.Get("email").(string)
	log.Printf("[DEBUG] Updating Dynamic Group %q: %#v", d.Id(), email)

	cloudIdentityService, diags := client.NewCloudIdentityService()
	if diags.HasError() {
		return diags
	}

	groupsService, diags := GetDynamicGroupsService(cloudIdentityService)
	if diags.HasError() {
		return diags
	}

	if d.HasChange("email") && d.HasChangesExcept("email") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "If you change the email address of a group, you must only change the email address.",
		})
		return diags

	}

	if d.HasChange("query") && d.HasChangesExcept("query") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "If you change the query of a group, you must only change the query.",
		})
		return diags
	}

	groupObj := cloudidentity.Group{}

	var updateMask []string

	if d.HasChange("email") {
		if groupObj.GroupKey == nil {
			groupObj.GroupKey = &cloudidentity.EntityKey{}
		}

		groupObj.GroupKey = &cloudidentity.EntityKey{
			Id: d.Get("email").(string),
		}

		updateMask = append(updateMask, "groupKey")
	}

	if d.HasChange("name") {
		groupObj.DisplayName = d.Get("name").(string)
		updateMask = append(updateMask, "displayName")
	}

	if d.HasChange("description") {
		groupObj.Description = d.Get("description").(string)
		updateMask = append(updateMask, "description")
	}

	if d.HasChange("query") {
		query := d.Get("query").(string)

		if groupObj.DynamicGroupMetadata == nil {
			groupObj.DynamicGroupMetadata = &cloudidentity.DynamicGroupMetadata{}
		}

		groupObj.DynamicGroupMetadata.Queries = []*cloudidentity.DynamicGroupQuery{
			{
				ResourceType: "USER",
				Query:        query,
			},
		}

		updateMask = append(updateMask, "dynamicGroupMetadata")
	}

	if d.HasChange("labels") {
		labels, err := convertInterfaceMapToStringMap(d.Get("labels").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}

		groupObj.Labels = labels
		updateMask = append(updateMask, "labels")
	}

	updateMaskStr := strings.Join(updateMask, ",")

	if &groupObj != new(cloudidentity.Group) {
		group, err := groupsService.Patch(d.Id(), &groupObj).UpdateMask(updateMaskStr).Do()
		if err != nil {
			return diag.FromErr(err)
		}

		var groupResponse Group

		err = json.Unmarshal(group.Response, &groupResponse)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(groupResponse.Name)
	}

	log.Printf("[DEBUG] Finished updating Dynamic Group %q: %#v", d.Id(), email)

	return resourceDynamicGroupRead(ctx, d, meta)
}

func resourceDynamicGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// use the meta value to retrieve your client from the provider configure method
	client := meta.(*apiClient)

	email := d.Get("email").(string)
	log.Printf("[DEBUG] Deleting Dynamic Group %q: %#v", d.Id(), email)

	cloudIdentityService, diags := client.NewCloudIdentityService()
	if diags.HasError() {
		return diags
	}

	groupsService, diags := GetDynamicGroupsService(cloudIdentityService)
	if diags.HasError() {
		return diags
	}

	_, err := groupsService.Delete(d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, d.Get("email").(string))
	}

	log.Printf("[DEBUG] Finished deleting Dynamic Group %q: %#v", d.Id(), email)

	return diags
}

// I'm leaving this here but it does not work. UpdatedKeys() returns an empty
// list for some reason. Need to debug this later.
func resourceExampleCustomizeDiff(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	mustBeChangedInIsolation := []string{"email", "query"}

	for _, attr := range mustBeChangedInIsolation {
		if diff.HasChange(attr) && len(diff.UpdatedKeys()) > 1 {
			return fmt.Errorf("If you change the %q of a group, you must only change the %q.", attr, attr)
		}
	}

	return nil
}

func convertInterfaceMapToStringMap(input map[string]interface{}) (map[string]string, error) {
	output := make(map[string]string)
	for key, value := range input {
		strValue, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("expected string for key %s but got %T", key, value)
		}
		output[key] = strValue
	}
	return output, nil
}
