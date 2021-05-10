---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "googleworkspace_group Resource - terraform-provider-googleworkspace"
subcategory: ""
description: |-
  Group resource manages Google Workspace Groups.
---

# googleworkspace_group (Resource)

Group resource manages Google Workspace Groups.

## Example Usage

```terraform
resource "googleworkspace_group" "sales" {
  email       = "sales@example.com"
  name        = "Sales"
  description = "Sales Group"

  aliases = ["paper-sales@example.com", "sales-dept@example.com"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **email** (String) The group's email address. If your account has multiple domains,select the appropriate domain for the email address. The email must be unique.

### Optional

- **aliases** (List of String) asps.list of group's email addresses.
- **description** (String) An extended description to help users determine the purpose of a group.For example, you can include information about who should join the group,the types of messages to send to the group, links to FAQs about the group, or related groups.
- **name** (String) The group's display name.
- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- **admin_created** (Boolean) Value is true if this group was created by an administrator rather than a user.
- **direct_members_count** (Number) The number of users that are direct members of the group.If a group is a member (child) of this group (the parent),members of the child group are not counted in the directMembersCount property of the parent group.
- **etag** (String) ETag of the resource.
- **id** (String) The unique ID of a group. A group id can be used as a group request URI's groupKey.
- **non_editable_aliases** (List of String) asps.list of the group's non-editable alias email addresses that are outside of theaccount's primary domain or subdomains. These are functioning email addresses used by the group.

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- **create** (String)
- **update** (String)

## Import

Import is supported using the following syntax:

```shell
terraform import googleworkspace_group.sales 01abcde23fg4h5i
```