---
page_title: "btp_subaccount Data Source - terraform-provider-btp"
subcategory: ""
description: |-
  Gets details about a subaccount.
  Tip:
  You must be assigned to the admin or viewer role of the global account, directory, or subaccount.
---

# btp_subaccount (Data Source)

Gets details about a subaccount.

__Tip:__
You must be assigned to the admin or viewer role of the global account, directory, or subaccount.

## Example Usage

```terraform
# Read a subaccount by ID
data "btp_subaccount" "my_account_byid" {
  id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
}

# Read a subaccount by region and subdomain
data "btp_subaccount" "my_account_bysubdomain" {
  region    = "eu10"
  subdomain = "my-subaccount-subdomain"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `id` (String) The ID of the subaccount.
- `region` (String) The region in which the subaccount was created.
- `subdomain` (String) The subdomain that becomes part of the path used to access the authorization tenant of the subaccount. Must be unique within the defined region. Use only letters (a-z), digits (0-9), and hyphens (not at the start or end). Maximum length is 63 characters. Cannot be changed after the subaccount has been created.

### Read-Only

- `beta_enabled` (Boolean) Shows whether the subaccount can use beta services and applications.
- `created_by` (String) The details of the user that created the subaccount.
- `created_date` (String) The date and time when the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.
- `description` (String) The description of the subaccount.
- `labels` (Map of Set of String) Set of words or phrases assigned to the subaccount.
- `last_modified` (String) The date and time when the resource was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.
- `name` (String) A descriptive name of the subaccount for customer-facing UIs.
- `parent_features` (Set of String) The features of parent entity of the subaccount.
- `parent_id` (String) The ID of the subaccount’s parent entity. If the subaccount is located directly in the global account (not in a directory), then this is the ID of the global account.
- `state` (String) The current state of the subaccount. Possible values are: 

  | state | description | 
  | --- | --- | 
  | `OK` | The CRUD operation or series of operations completed successfully. | 
  | `STARTED` | CRUD operation on the subaccount has started. | 
  | `CANCELED` | The operation or processing was canceled by the operator. | 
  | `PROCESSING` | A series of operations related to the subaccount are in progress. | 
  | `PROCESSING_FAILED` | The processing operations failed. | 
  | `CREATING` | Creating the subaccount is in progress. | 
  | `CREATION_FAILED` | The creation operation failed, and the subaccount was not created or was created but cannot be used. | 
  | `UPDATING` | Updating the subaccount is in progress. | 
  | `UPDATE_FAILED` | The update operation failed, and the subaccount was not updated. | 
  | `UPDATE_DIRECTORY_TYPE_FAILED` | The update of the directory type failed. | 
  | `UPDATE_ACCOUNT_TYPE_FAILED` | The update of the account type failed. | 
  | `DELETING` | Deleting the subaccount is in progress. | 
  | `DELETION_FAILED` | The deletion of the subaccount failed, and the subaccount was not deleted. | 
  | `MOVING` | Moving the subaccount is in progress. | 
  | `MOVE_FAILED` | The moving of the subaccount failed. | 
  | `MOVING_TO_OTHER_GA` | Moving the subaccount to another global account is in progress. | 
  | `MOVE_TO_OTHER_GA_FAILED` | Moving the subaccount to another global account failed. | 
  | `PENDING_REVIEW` | The processing operation has been stopped for reviewing and can be restarted by the operator. | 
  | `MIGRATING` | Migrating the subaccount from Neo to Cloud Foundry. | 
  | `MIGRATED` | The migration of the subaccount completed. | 
  | `MIGRATION_FAILED` | The migration of the subaccount failed and the subaccount was not migrated. | 
  | `ROLLBACK_MIGRATION_PROCESSING` | The migration of the subaccount was rolled back and the subaccount is not migrated. | 
  | `SUSPENSION_FAILED` | The suspension operations failed. |
- `usage` (String) Shows whether the subaccount is used for production purposes. This flag can help your cloud operator to take appropriate action when handling incidents that are related to mission-critical accounts in production systems. Do not apply for subaccounts that are used for non-production purposes, such as development, testing, and demos. Applying this setting this does not modify the subaccount. Possible values are: 

  | value | description | 
  | --- | --- | 
  | `UNSET` | Global account or subaccount admin has not set the production-relevancy flag (default value). | 
  | `NOT_USED_FOR_PRODUCTION` | The subaccount is not used for production purposes. | 
  | `USED_FOR_PRODUCTION` | The subaccount is used for production purposes. |