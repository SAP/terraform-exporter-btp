# Terraform Setup of BTP Exporter Integration test

## Setup of infrastructure

The setup comprises:

- A managed directory on SAP BTP
- A subaccount including a Cloud Foundry runtime
- A Cloud Foundry organization with two spaces

The setup of the infrastructure is handled via Terramate. The following steps are needed to execute the setup of the sample used for integration testing:

1. Initialize the workspaces via

  ```bash
  terramate run -X --tags dev --parallel 2 terraform init
  ```

2.  Export variables for setup e.g., via `export $(xargs <.env)`

3. Execute the deployment

  ```bash
  terramate script run --tags dev -X deploy
  ```

In case of a local storage of the state you can use the `teardown` script via

```bash
terramate script run --tags dev -X --reverse teardown
```

> [!NOTE]
> While the setup on SAP BTP could be done automatically, it makes sense to keep it static as we do not want to test if the setup via Terraform works, but if the Terraform Exporter for BTP is working as expected.

## Integration test

The integration test for the exporter is based on the reference setup described in the previous section

Configure the `CF_ORIGIN` secret in the GitHub `release` environment to select the Cloud Foundry identity provider. All three CF flows (JSON inventory, export by JSON, and export by resource) expose it as the `CF_ORIGIN` environment variable, following the same pattern as `BTP_IDP` for BTP. For local runs and reference generation, export `CF_ORIGIN` alongside the other CF credentials.

The test comprises the different flows and compares the results with a reference file:

- JSON inventory: reference file for the resulting JSON
- Export: reference state file

We must distinguish the two scenarios with respect to the checks in the integration test.

### Validation of the JSON inventory

The check if the export is valid namely if the reference JSON inventory corresponds to the current export can be done using `diff`:

```
diff <(jq -S . btpResources_new.json) <(jq -S . btpResources_reference.json)
```

As the JSON files are flat and contain no metadata the `diff` is sufficient.

### Validation of the "Export By" flows

For exports, we compare the identities of the resources in the Terraform states after import. Capture raw state, which includes the `identity` block of each resource instance:

```bash
terraform state pull > reference-state.json
```

As the state contains sensitive data, the reference state is base64-encoded and stored as a GitHub secret.

> [!NOTE]
> Be aware that a GitHub Secret can only contain up to 48 kb of data

The comparison matches resource instances one-to-one using their resource type, mode, and complete `identity` object. Resource order, identity object key order, generated resource names, module addresses, and instance indexes do not affect matching. Every instance is included, including resources in modules and resources using `count` or `for_each`. Missing or extra instances, including surplus duplicates, fail the comparison.

Attributes, sensitive-attribute markers, dependencies, schema versions, and state metadata are not compared. This validates that the exporter imports the expected objects; it does not validate attribute equality. Every instance must have an identity containing at least one non-null value; optional identity fields may be null. Missing or empty identities, unfinished replacements with deposed instances, and unsupported state formats fail with an error.

Existing reference secrets generated with `terraform show -json` must be regenerated using `terraform state pull` after refreshing or importing with providers that support resource identities. Update all six secrets: `DIRECTORY_EXPORT_BY_JSON_REF`, `SUBACCOUNT_EXPORT_BY_JSON_REF`, `CFORG_EXPORT_BY_JSON_REF`, `DIRECTORY_EXPORT_BY_RESOURCE_REF`, `SUBACCOUNT_EXPORT_BY_RESOURCE_REF`, and `CFORG_EXPORT_BY_RESOURCE_REF`. The reference-generation commands below produce the required format.

The script is available at [`.github/scripts/compareJson.js`](../.github/scripts/compareJson.js).

Run the comparison and its unit tests with:

```bash
node .github/scripts/compareJson.js reference-state.json actual-state.json
npm --prefix .github/scripts test
```

### Flows and levels to check

The following matrix needs to be checked:

Supported levels:

- Directory level
- Subaccount level
- CF Organization level

Supported flows

- Creation of JSON inventory
- Export by JSON
- Export by resource

The export by JSON will be executed based on a curated JSON file to reduce the number of entitlements and roles.
The export by resource will be restricted so some resources as the basic export logic is already covered in the export by JSON flow using the same function for the export under the hood


## Terraform exporter commands

### JSON inventory

#### Directory

For reference data:

```bash
btptf create-json -d <directoryID> -p integrationTestDirectoryReference.json
jq '.BtpResources[] |= (.Values |= sort)' integrationTestDirectoryReference.json > temp.json && mv temp.json integrationTestDirectoryReference.json
```

For integration test data:

```bash
btptf create-json -d <directoryID> -p integrationTestDirectory.json
```

#### Subaccount

For reference data:

```bash
btptf create-json -s <subaccountID> -p integrationTestSubaccountReference.json
jq '.BtpResources[] |= (.Values |= sort)' integrationTestSubaccountReference.json > temp.json && mv temp.json integrationTestSubaccountReference.json
```

For integration test data:

```bash
btptf create-json -s <subaccountID> -p integrationTestSubaccount.json
```

#### Cloud Foundry Organization

For reference data:

```bash
btptf create-json -o <organizationID> -p integrationTestCfOrgReference.json
jq '.BtpResources[] |= (.Values |= sort)' integrationTestCfOrgReference.json > temp.json && mv temp.json integrationTestCfOrgReference.json
```

For integration test data:

```bash
btptf create-json -o <organizationID> -p integrationTestCfOrg.json
```

### Export by Resource

#### Directory

For reference data:

```bash
btptf export -d <directoryID> -c directory-export-by-resource-ref -r=directory
terraform -chdir=directory-export-by-resource-ref init
terraform -chdir=directory-export-by-resource-ref apply -auto-approve
terraform -chdir=directory-export-by-resource-ref state pull > directory-export-by-resource-ref.json
base64 -i directory-export-by-resource-ref.json > directory-export-by-resource-ref
```

> [!NOTE]
> The base 64 encoded string gets stored as a GitHub secret

For integration test:

```bash
btptf export -d <directoryID> -c directory-export-by-resource -r=directory
terraform -chdir=directory-export-by-resource init
terraform -chdir=directory-export-by-resource apply -auto-approve
terraform -chdir=directory-export-by-resource state pull > directory-export-by-resource.json
```

#### Subaccount

For reference data:

```bash
btptf export -s <subaccountID> -c subaccount-export-by-resource-ref -r='subaccount,subscriptions'
terraform -chdir=subaccount-export-by-resource-ref init
terraform -chdir=subaccount-export-by-resource-ref apply -auto-approve
terraform -chdir=subaccount-export-by-resource-ref state pull > subaccount-export-by-resource-ref.json
base64 -i subaccount-export-by-resource-ref.json > subaccount-export-by-resource-ref
```

> [!NOTE]
> The base 64 encoded string gets stored as a GitHub secret

For integration test:

```bash
btptf export -s <subaccountID> -c subaccount-export-by-resource -r='subaccount,subscriptions'
terraform -chdir=subaccount-export-by-resource init
terraform -chdir=subaccount-export-by-resource apply -auto-approve
```

#### Cloud Foundry Organization

For reference data:

```bash
btptf export -o <organizationID> -c cforg-export-by-resource-ref -r='spaces'
terraform -chdir=cforg-export-by-resource-ref init
terraform -chdir=cforg-export-by-resource-ref apply -auto-approve
terraform -chdir=cforg-export-by-resource-ref state pull > cforg-export-by-resource-ref.json
base64 -i cforg-export-by-resource-ref.json > cforg-export-by-resource-ref
```

> [!NOTE]
> The base 64 encoded string gets stored as a GitHub secret

For integration test:

```bash
btptf export -o <organizationID> -c cforg-export-by-resource -r='spaces'
terraform -chdir=cforg-export-by-resource init
terraform -chdir=cforg-export-by-resource apply -auto-approve
```

### Export by JSON

#### Directory

For reference data (state created by Terramate setup):

```bash
btptf export-by-json -d <directoryID> -p integrationTestDirectoryCurated.json -c directory-export-by-json
terraform -chdir=directory-export-by-json init
terraform -chdir=directory-export-by-json apply -auto-approve
terraform -chdir=directory-export-by-json state pull > directory-export-by-json-ref.json
base64 -i directory-export-by-json-ref.json > directory-export-by-json-ref
```

> [!NOTE]
> The base 64 encoded string gets stored as a GitHub secret

For integration test:

```bash
btptf export-by-json -d <directoryID> -p integrationTestDirectoryCurated.json -c directory-export-by-json
terraform -chdir=directory-export-by-json init
terraform -chdir=directory-export-by-json apply -auto-approve
```

#### Subaccount

For reference data:

```bash
btptf export-by-json -s <subaccountID> -p integrationTestSubaccountCurated.json -c subaccount-export-by-json
terraform -chdir=subaccount-export-by-json init
terraform -chdir=subaccount-export-by-json apply -auto-approve
terraform -chdir=subaccount-export-by-json state pull > subaccount-export-by-json-ref.json
base64 -i subaccount-export-by-json-ref.json > subaccount-export-by-json-ref
```

> [!NOTE]
> The base 64 encoded string gets stored as a GitHub secret

For integration test:

```bash
btptf export-by-json -s <subaccountID> -p integrationTestSubaccountCurated.json -c subaccount-export-by-json
terraform -chdir=subaccount-export-by-json init
terraform -chdir=subaccount-export-by-json apply -auto-approve
```

#### Cloud Foundry Organization

For reference data:

```bash
btptf export-by-json -o <organizationID> -p integrationTestCfOrgCurated.json -c cforg-export-by-json
terraform -chdir=cforg-export-by-json init
terraform -chdir=cforg-export-by-json apply -auto-approve
terraform -chdir=cforg-export-by-json state pull > cforg-export-by-json-ref.json
base64 -i cforg-export-by-json-ref.json > cforg-export-by-json-ref
```

> [!NOTE]
> The base 64 encoded string gets stored as a GitHub secret

For integration test:

```bash
btptf export-by-json -o <organizationID> -p integrationTestCfOrgCurated.json -c cforg-export-by-json
terraform -chdir=cforg-export-by-json init
terraform -chdir=cforg-export-by-json apply -auto-approve
```

### Terraform metadata and JSON format

The workflow pins the Terraform CLI version. The comparator supports raw state version 4 and deliberately rejects other formats so format changes cannot silently pass validation.
