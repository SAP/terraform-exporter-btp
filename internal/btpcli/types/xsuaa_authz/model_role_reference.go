/*
 * Authorization
 *
 * Provides functions to administrate the Authorization and Trust Management service (XSUAA) of SAP BTP, Cloud Foundry environment. You can manage service instances of the Authorization and Trust Management service. You can also manage roles, role templates, and role collections of your subaccount.
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package xsuaa_authz

type RoleReference struct {
	RoleTemplateAppId string `json:"roleTemplateAppId,omitempty"`
	// The name has a maximum length of 64 characters. Only the following characters are allowed: alphanumeric characters (aA-zZ) and (0-9), underscore (_), period (.), and hyphen (-).
	RoleTemplateName string `json:"roleTemplateName,omitempty"`
	Description      string `json:"description,omitempty"`
	Name             string `json:"name,omitempty"`
}
