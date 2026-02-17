![Golang](https://img.shields.io/badge/Go-1.26-informational)
[![REUSE status](https://api.reuse.software/badge/github.com/SAP/terraform-exporter-btp)](https://api.reuse.software/info/github.com/SAP/terraform-exporter-btp)
[![Go Report Card](https://goreportcard.com/badge/github.com/SAP/terraform-exporter-btp)](https://goreportcard.com/report/github.com/SAP/terraform-exporter-btp)
[![CodeQL](https://github.com/SAP/terraform-exporter-btp/actions/workflows/codeql.yml/badge.svg)](https://github.com/SAP/terraform-exporter-btp/actions/workflows/codeql.yml)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/9673/badge)](https://www.bestpractices.dev/projects/9673)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=SAP_terraform-exporter-btp&metric=coverage)](https://sonarcloud.io/summary/new_code?id=SAP_terraform-exporter-btp)


# Terraform Exporter for SAP BTP

## Overview
The *Terraform Exporter for SAP BTP* (btptf CLI) is a handy tool that makes it easier to bring your existing SAP Business Technology Platform (BTP) resources into Terraform. With it, you can take things like subaccounts and directories in BTP and turn them into configuration files including the import of the state. It's especially useful for teams who are moving to Terraform but still need to manage older infrastructure or SAP BTP accounts that are already set up.

Here's how it works:

- **Resource Identification**: Terraform Exporter for SAP BTP identifies the SAP BTP resources and maps them to corresponding Terraform resources.
- **Import Process**: The tool uses Terraform's import function to create the import blocks required for the state import.
- **Configuration Generation**: It generates the Terraform code (in HashiCorp Configuration Language - HCL) for each resource, enabling further customizations as needed.

You can install btptf CLI across various operating systems as described below.

## Documentation

You find the documentation of the Terraform Exporter for SAP BTP on the [GitHub page of this repository](https://sap.github.io/terraform-exporter-btp/).

Here are some quick links to the most important documentation sections:

- [Installation](https://sap.github.io/terraform-exporter-btp/install/)
- Quickstart
   - [Prerequisites](https://sap.github.io/terraform-exporter-btp/prerequisites/)
   - [Getting started](https://sap.github.io/terraform-exporter-btp/gettingstarted/)
- [Limitations](https://sap.github.io/terraform-exporter-btp/limitations/)
- [Troubleshooting](https://sap.github.io/terraform-exporter-btp/troubleshooting/)

## Developer Guide

If you want to contribute to the code of the Terraform Exporter for SAP BTP, please check our [Contribution Guidelines](CONTRIBUTING.md). The technical setup and how to get started are described in the [Developer Guide](./guidelines/DEVELOPER-GUIDE.md)

## Support, Feedback, Contributing

This project is open to feature requests/suggestions, bug reports, and so on, via [GitHub issues](https://github.com/SAP/terraform-exporter-for-sap-btp/issues):

❓ - If you have a *question* you can ask it here in [GitHub Discussions](https://github.com/SAP/terraform-exporter-btp/discussions/).

🐞 - If you find a bug, feel free to create a [bug report](https://github.com/SAP/terraform-exporter-btp/issues/new?assignees=&labels=bug%2Cneeds-triage&projects=&template=bug_report.yml&title=%5BBUG%5D).

💡 - If you have an idea for improvement or a feature request, please open a [feature request](https://github.com/SAP/terraform-exporter-btp/issues/new?assignees=&labels=enhancement%2Cneeds-triage&projects=&template=feature_request.yml&title=%5BFEATURE%5D).

Contribution and feedback are encouraged and always welcome. For more information about how to contribute, the project structure, as well as additional contribution information, see our [Contribution Guidelines](CONTRIBUTING.md).

## Security / Disclosure
If you find any bug that may be a security problem, please follow our instructions at [in our security policy](https://github.com/SAP/terraform-exporter-for-sap-btp/security/policy) on how to report it. Please do not create GitHub issues for security-related doubts or problems.

## Code of Conduct

We as members, contributors, and leaders pledge to make participation in our community a harassment-free experience for everyone. By participating in this project, you agree to abide by its [Code of Conduct](https://github.com/SAP/.github/blob/main/CODE_OF_CONDUCT.md) at all times.

## Licensing

Copyright 2026 SAP SE or an SAP affiliate company and terraform-exporter-for-sap-btp contributors. Please see our [LICENSE](LICENSE) for copyright and license information. Detailed information including third-party components and their licensing/copyright information is available [via the REUSE tool](https://api.reuse.software/info/github.com/SAP/terraform-exporter-btp).
