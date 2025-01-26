# Terraform Setup of BTP Exporter Integration test

## Setup of infrastructure

The setup comprises:

- A managed dircetory on SAP BTP
- A subaccount inclusinga Cloud Foundry runtime
- A Cloud Foundry organization with two spaces

The setup of the infrastructure is handled via Terramate. The following steps are needed to execute the setup of the sample used for integration testing:

1. Initialize the workpsaces via

  ```bash
  terramate run -X --tags dev --parallel 2 terraform init
  ```

2.  Export variables for setup e.g., via `export $(xargs <.env)`

3. Execute the deployment

  ```bash
  terramate script run --tags dev -X deploy
  ```

In case of a local storage of the state you can use the *teardown* scrtipt via

```bash
terramate script run --tags dev -X --reverse teardown
```

## Integration test

The integration test for the exporter is based on the refernce setup described in the previous section

The test comprises the different flows and copmpares the results with a reference file:

- JSON inventory: reference file for the resulting JSON
- Export: reference state file


The check if the export is working at the point in time of the creation of the reference file is done by comparing a newly created file with the refernence files using `diff`.

For the JSON inventory this is achieved vai the the following statement

```
diff <(jq -S . btpResources_new.json) <(jq -S . btpResources_reference.json)
```

In coase of the exports we must compare the resulting Terraform states after executing the import. As the state contains sensitive data, the reference state is stored as a GitHub secret. As it is JSON format, we transfer it into a string using base64 encoding.

The resulting comparison is then done in analogy to the JSON inventory file after decoding the string from the secret

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
```

For integration test data:

```bash
btptf create-json -d <directoryID> -p integrationTestDirectory.json
```

#### Subaccount

For reference data:

```bash
btptf create-json -s <subaccountID> -p integrationTestSubaccountReference.json
```

For integration test data:

```bash
btptf create-json -s <subaccountID> -p integrationTestSubaccount.json
```

#### Cloud Foundry Organization

For reference data:

```bash
btptf create-json -o <organizationID> -p integrationTestCfOrgReference.json
```

For integration test data:

```bash
btptf create-json -o <organizationID> -p integrationTestCfOrg.json
```

### Export by Resource

#### Directory

For integration test:

```bash
btptf export -d <directoryID> -c directory-export-by-resource -r=directory
terraform -chdir=directory-export-by-resource init
terraform -chdir=directory-export-by-resource apply -auto-approve
```

#### Subaccount

For integration test:

```bash
btptf export -s <subaccountID> -c subaccount-export-by-resource -r='subaccount,subscriptions'
terraform -chdir=subaccount-export-by-resource init
terraform -chdir=subaccount-export-by-resource apply -auto-approve
```

#### Cloud Foundry Organization

For integration test:

```bash
btptf export -o <organizationID> -c cforg-export-by-resource -r='spaces'
terraform -chdir=cforg-export-by-resource init
terraform -chdir=cforg-export-by-resource apply -auto-approve
```

### Export by JSON

#### Directory

For integration test:

```bash
btptf export-by-json -d <directoryID> -p integrationTestDirectoryCurated.json -c directory-export-by-json
terraform -chdir=directory-export-by-json init
terraform -chdir=directory-export-by-json apply -auto-approve
```

#### Subaccount

For integration test:

```bash
btptf export-by-json -s <subaccountID> -p integrationTestSubaccountCurated.json -c subaccount-export-by-json
terraform -chdir=subaccount-export-by-json init
terraform -chdir=subaccount-export-by-json apply -auto-approve
```

#### Cloud Foundry Organization

For integration test:

```bash
btptf export-by-json -o <organizationID> -p integrationTestCfOrgCurated.json -c cforg-export-by-json
terraform -chdir=cforg-export-by-json init
terraform -chdir=cforg-export-by-json apply -auto-approve
```

### Terraform metadata

As the Terraform state contains some metadata e.g., around the used Terraform CLI version, we must make sure that the setup in the Gh Action comprises the same Terraform version as the one used to record the reference data
