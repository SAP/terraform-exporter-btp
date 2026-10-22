# Migrating to Version 2

Version 2 uses Terraform's list-resource query workflow for discovery, import blocks, and resource configuration.

## Breaking changes

- Terraform CLI 1.14 or newer is required.
- OpenTofu is no longer supported because it does not provide the required `query` command.
- Cloud Foundry users are discovered with the `cloudfoundry_user` list resource, then generated resource and import blocks are converted to `cloudfoundry_user_cf`.
- The JSON inventory name for Cloud Foundry users is `users` instead of `users_cf`.
- When Cloud Foundry users have the same username from multiple identity-provider origins, their inventory values use `username__origin` so they can be selected independently.
- Terraform assigns generated resource addresses from query labels. Resource addresses can therefore differ from Version 1 output.

The command structure, editable JSON inventory workflow, cleanup, resume support, backend configuration, and summary output remain available. Create a new Version 2 inventory before exporting instead of reusing an inventory containing `users_cf`.
