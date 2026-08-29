Hi there,

Thank you for opening an issue.

**Please note:** this is a community-maintained fork of the Google Workspace provider and is **not** affiliated with, supported by, or maintained by HashiCorp. Issues here should be specific to this provider.

For anything that isn't specific to the Google Workspace provider, please use the appropriate channel instead:
- **Terraform core** (bugs in `terraform` itself): https://github.com/hashicorp/terraform/issues
- **General usage questions and help**: https://discuss.hashicorp.com/c/terraform-core
- **Other providers**: open an issue on that provider's own repository

This issue tracker is reserved for bug reports and feature requests specific to this provider.

### Terraform Version
Run `terraform -v` to show the version. If you are not running the latest version of Terraform, please upgrade because your issue may have already been fixed.

### Affected Resource(s)
Please list the resources as a list, for example:
- opc_instance
- opc_storage_volume

If this issue appears to affect multiple resources, it may be an issue with Terraform's core, so please mention this.

### Terraform Configuration Files
```hcl
# Copy-paste your Terraform configurations here - for large Terraform configs,
# please use a service like Dropbox and share a link to the ZIP file. For
# security, you can also encrypt the files using our GPG public key.
```

### Debug Output
Please provider a link to a GitHub Gist containing the complete debug output: https://www.terraform.io/docs/internals/debugging.html. Please do NOT paste the debug output in the issue; just paste a link to the Gist.

### Panic Output
If Terraform produced a panic, please provide a link to a GitHub Gist containing the output of the `crash.log`.

### Expected Behavior
What should have happened?

### Actual Behavior
What actually happened?

### Steps to Reproduce
Please list the steps required to reproduce the issue, for example:
1. `terraform apply`

### Important Factoids
Are there anything atypical about your accounts that we should know? For example: Running in EC2 Classic? Custom version of OpenStack? Tight ACLs?

### References
Are there any other GitHub issues (open or closed) or Pull Requests that should be linked here? For example:
- GH-1234
