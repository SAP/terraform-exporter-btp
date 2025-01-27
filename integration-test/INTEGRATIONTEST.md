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

In case of the exports we must compare the resulting Terraform states after executing the import. As the `tfstate` format is an internal representation we must make the the state files comparable by transfering them to the canonical JSON format:

```bash
terraform show -json > <Some Name>state.json
```

As the state contains sensitive data, the reference state is stored as a GitHub secret. As it is JSON format, we transfer it into a base64-encoded.

> [!NOTE]
> Be aware that a GitHub Secret can only contain up to 48 kb of data

For the comparison we cannot rely on a simple `diff` as the sequence of the resources might differ. The same is true for the metadata of the resources like `address` or `name` of the resource in the state file. Consequently we provide a JS script that compares the `value` and `sensitive` value part of the state files in JSON format. If all entries match, the validation passes.

The script is available at [`.github/scripts/compareJson.js`](../.github/scripts/compareJson.js).

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
terraform -chdir=directory-export-by-resource-ref show -json > directory-export-by-resource-ref.json
base64 -i directory-export-by-resource-ref.json > directory-export-by-resource-ref
```

> [!NOTE]
> The base 64 encoded string gets stored as a GitHub secret

For integration test:

```bash
btptf export -d <directoryID> -c directory-export-by-resource -r=directory
terraform -chdir=directory-export-by-resource init
terraform -chdir=directory-export-by-resource apply -auto-approve
terraform -chdir=directory-export-by-resource show -json > directory-export-by-resource.json
```

#### Subaccount

For reference data:

```bash
btptf export -s <subaccountID> -c subaccount-export-by-resource-ref -r='subaccount,subscriptions'
terraform -chdir=subaccount-export-by-resource-ref init
terraform -chdir=subaccount-export-by-resource-ref apply -auto-approve
terraform -chdir=subaccount-export-by-resource-ref show -json > subaccount-export-by-resource-ref.json
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
terraform -chdir=cforg-export-by-resource-ref show -json > cforg-export-by-resource-ref.json
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
terraform -chdir=directory-export-by-json show -json > directory-export-by-json-ref.json
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
terraform -chdir=subaccount-export-by-json-ref show -json > subaccount-export-by-json-ref.json
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
terraform -chdir=cforg-export-by-json-ref show -json > cforg-export-by-json-ref.json
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

To keep the reference files comparable with the newly created files, we must pin the Terraform CLI version to a specific version to avoid inconsistencies especially regarding the JSON format of the state.
