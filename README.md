![Golang](https://img.shields.io/badge/Go-1.23-informational)
[![REUSE status](https://api.reuse.software/badge/github.com/SAP/terraform-exporter-btp)](https://api.reuse.software/info/github.com/SAP/terraform-exporter-btp)
[![Go Report Card](https://goreportcard.com/badge/github.com/SAP/terraform-exporter-btp)](https://goreportcard.com/report/github.com/SAP/terraform-exporter-btp)
[![CodeQL](https://github.com/SAP/terraform-exporter-btp/actions/workflows/codeql.yml/badge.svg)](https://github.com/SAP/terraform-exporter-btp/actions/workflows/codeql.yml)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/9673/badge)](https://www.bestpractices.dev/projects/9673)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=SAP_terraform-exporter-btp&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=SAP_terraform-exporter-btp)


# Terraform Exporter for SAP BTP

> [!CAUTION]
> The Terraform Exporter for SAP BTP is still under development. It is not intended for productive useage in its current state.

## Overview
The *Terraform Exporter for SAP BTP* (btptf CLI) is a handy tool that makes it easier to bring your existing SAP Business Technology Platform (BTP) resources into Terraform. With it, you can take things like subaccounts and directories in BTP and turn them into Terraform state and configuration files. It's especially useful for teams who are moving to Terraform but still need to manage older infrastructure or SAP BTP accounts that are already set up.

Here's how it works:

- **Resource Identification**: Terraform Exporter for SAP BTP identifies the SAP BTP resources and maps them to corresponding Terraform resources using the BTP CLI Server APIs.
- **Import Process**: The tool utilizes Terraform's import function to integrate each resource into Terraform's state.
- **Configuration Generation**: After import, it generates the Terraform code (in HashiCorp Configuration Language - HCL) for each resource, enabling further customizations as needed.

You can install btptf CLI across various operating systems as described below.

## Documentation

You find the documentation of the Terraform Exporter for SAP BTP on the [GitHub page of this repository](https://sap.github.io/terraform-exporter-btp/).

## Developer Guide

If you want to contribute to the code of the Terraform Exporter for SAP BTP, please check our [Contribution Guidelines](CONTRIBUTING.md). The technical setup and how to get started are described in the [Developer Guide](./guidelines/DEVELOPER-GUIDE.md)

## Support, Feedback, Contributing

This project is open to feature requests/suggestions, bug reports, and so on, via [GitHub issues](https://github.com/SAP/terraform-exporter-for-sap-btp/issues). Contribution and feedback are encouraged and always welcome. For more information about how to contribute, the project structure, as well as additional contribution information, see our [Contribution Guidelines](CONTRIBUTING.md).

## Security / Disclosure
If you find any bug that may be a security problem, please follow our instructions at [in our security policy](https://github.com/SAP/terraform-exporter-for-sap-btp/security/policy) on how to report it. Please do not create GitHub issues for security-related doubts or problems.

## Code of Conduct

We as members, contributors, and leaders pledge to make participation in our community a harassment-free experience for everyone. By participating in this project, you agree to abide by its [Code of Conduct](https://github.com/SAP/.github/blob/main/CODE_OF_CONDUCT.md) at all times.

## Licensing

Copyright 2025 SAP SE or an SAP affiliate company and terraform-exporter-for-sap-btp contributors. Please see our [LICENSE](LICENSE) for copyright and license information. Detailed information including third-party components and their licensing/copyright information is available [via the REUSE tool](https://api.reuse.software/info/github.com/SAP/terraform-exporter-btp).
