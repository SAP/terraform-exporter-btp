# Security Considerations

When using the Terraform Exporter for SAP BTP, following security best practices is key to ensuring that your SAP BTP resources and Terraform configurations are managed safely. Here are some recommended security practices:

1. **Access Control and Least Privilege**: Ensure that the user running the btptf CLI has only the necessary permissions, following the principle of least privilege. This helps minimize the risk of accidental or malicious modifications to resources​.
​
2. **Secure Storage for State Files**: Enable encryption for additional security for Terraform stat files. This is to protect sensitive data in Terraform’s state files, such as resource IDs and secrets. Store these state files in secure, remote back-ends configured with access policies.

3. **Network Security and Encryption**: When using remote state storage, ensure that data is encrypted both in transit and at rest.

4. **Terraform Code and Secrets Management**: Avoid hard coding secrets directly in Terraform configurations. Instead, use a secure vault services to manage secrets and inject them into configurations at runtime. This reduces the risk of exposing sensitive information in code repositories​

5. **Review and Update Terraform Configurations Regularly**: Ensure that exported configurations align with current security and compliance requirements. Regularly review and update configurations, removing any deprecated or unused resources that could expose vulnerabilities.

Following these guidelines can greatly reduce the security risks associated with managing SAP BTP resources using Terraform.

## Export Directory Containment

The exporter's file operations use retained directory handles to keep generated files inside the selected output directory. Export paths must be relative descendants of the working directory. Recursive overwrite is confined to that directory tree and still requires confirmation. Temporary query directories have random names and are cleaned up only through their allocation handles.

Before invoking Terraform, the exporter verifies that its directory has not been replaced and rejects escaping or dangling symbolic links in the directory tree. This includes external provider-cache links under `.terraform/providers`; shared provider caches that create those links must be disabled for exports. Relative symbolic links whose targets remain inside the directory are supported. Absolute symbolic links are rejected by the rooted filesystem operations, even if their target is inside the directory.

Terraform execution is not isolated. The directory must be trusted against concurrent modification while Terraform runs. These checks do not sandbox Terraform configuration, provider executables, environment settings, hard links, or filesystem mounts.
