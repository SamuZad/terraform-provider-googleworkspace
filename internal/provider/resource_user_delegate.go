package googleworkspace

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/gmail/v1"
)

func resourceUserDelegate() *schema.Resource {
	return &schema.Resource{
		Description: "User Delegate resource manages delagation of access to a Gmail user mailbox. User Delegate resides under the " +
			"`https://www.googleapis.com/auth/gmail.settings.sharing` client scope.",

		CreateContext: resourceUserDelegateCreate,
		ReadContext:   resourceUserDelegateRead,
		DeleteContext: resourceUserDelegateDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: resourceUserDelegateImport,
		},

		Schema: map[string]*schema.Schema{
			"user_id": {
				Description: "The user's email address.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"delegate_email": {
				Description: "The email address of the delegate.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"delegate_verification_status": {
				Description: "Indicates whether this address has been verified and can act as a delegate for the account. Read-only.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"id": {
				Description: "The ID of this resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceUserDelegateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)

	userId := d.Get("user_id").(string)
	gmailService, diags := client.NewGmailService(ctx, userId)
	if diags.HasError() {
		return diags
	}

	usersSettingsDelegatesService, diags := GetGmailUsersSettingsDelegatesService(gmailService)
	if diags.HasError() {
		return diags
	}

	delegateEmail := d.Get("delegate_email").(string)
	log.Printf("[INFO] Creating delegate %q for user %q", delegateEmail, userId)

	delegate, err := usersSettingsDelegatesService.Create(userId, &gmail.Delegate{
		DelegateEmail: delegateEmail,
	}).Do()
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("delegate_verification_status", delegate.VerificationStatus)
	d.SetId(fmt.Sprintf("%s:%s", userId, delegateEmail))

	log.Printf("[INFO] Created delegate %q for user %q", delegateEmail, userId)

	return resourceUserDelegateRead(ctx, d, meta)
}

func resourceUserDelegateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)

	userId := d.Get("user_id").(string)
	delegateEmail := d.Get("delegate_email").(string)

	gmailService, diags := client.NewGmailService(ctx, userId)
	if diags.HasError() {
		return diags
	}

	usersSettingsDelegatesService, diags := GetGmailUsersSettingsDelegatesService(gmailService)
	if diags.HasError() {
		return diags
	}

	log.Printf("[INFO] Reading delegate %q for user %q", delegateEmail, userId)

	delegate, err := usersSettingsDelegatesService.Get(userId, delegateEmail).Do()
	if err != nil {
		return handleNotFoundError(err, d, d.Id())
	}

	log.Printf("[INFO] Read delegate %q for user %q", delegateEmail, userId)

	d.SetId(fmt.Sprintf("%s:%s", userId, delegateEmail))
	d.Set("user_id", userId)
	d.Set("delegate_email", delegate.DelegateEmail)
	d.Set("delegate_verification_status", delegate.VerificationStatus)

	return nil
}

func resourceUserDelegateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)

	userId := d.Get("user_id").(string)
	delegateEmail := d.Get("delegate_email").(string)

	gmailService, diags := client.NewGmailService(ctx, userId)
	if diags.HasError() {
		return diags
	}

	usersSettingsDelegatesService, diags := GetGmailUsersSettingsDelegatesService(gmailService)
	if diags.HasError() {
		return diags
	}

	log.Printf("[INFO] Deleting delegate %q for user %q", delegateEmail, userId)

	err := usersSettingsDelegatesService.Delete(userId, delegateEmail).Do()
	if err != nil {
		return handleNotFoundError(err, d, d.Id())
	}

	log.Printf("[INFO] Deleted delegate %q for user %q", delegateEmail, userId)

	return nil
}

func resourceUserDelegateImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	idParts := strings.Split(d.Id(), ":")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		return nil, fmt.Errorf("Unexpected format of ID (%q), expected user_id:delegate_email", d.Id())
	}
	d.Set("user_id", idParts[0])
	d.Set("delegate_email", idParts[1])
	return []*schema.ResourceData{d}, nil
}
