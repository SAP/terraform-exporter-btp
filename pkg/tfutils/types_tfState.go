package tfutils

import "encoding/json"

// Reduced definition of state as provided by https://github.com/hashicorp/terraform-json

type State struct {
	// The version of the state format. This should always match the
	// StateFormatVersion constant in this package, or else am
	// unmarshal will be unstable.
	FormatVersion string `json:"format_version,omitempty"`

	// The Terraform version used to make the state.
	TerraformVersion string `json:"terraform_version,omitempty"`

	// The values that make up the state.
	Values *StateValues `json:"values,omitempty"`
}

type StateValues struct {
	// The root module in this state representation.
	RootModule *StateModule `json:"root_module,omitempty"`
}

type StateModule struct {
	// All resources or data sources within this module.
	Resources []*StateResource `json:"resources,omitempty"`

	// The absolute module address, omitted for the root module.
	Address string `json:"address,omitempty"`

	// Any child modules within this module.
	ChildModules []*StateModule `json:"child_modules,omitempty"`
}

type StateResource struct {
	// The absolute resource address.
	Address string `json:"address,omitempty"`

	// The resource mode.
	Mode ResourceMode `json:"mode,omitempty"`

	// The resource type, example: "aws_instance" for aws_instance.foo.
	Type string `json:"type,omitempty"`

	// The resource name, example: "foo" for aws_instance.foo.
	Name string `json:"name,omitempty"`

	// The instance key for any resources that have been created using
	// "count" or "for_each". If neither of these apply the key will be
	// empty.
	//
	// This value can be either an integer (int) or a string.
	Index interface{} `json:"index,omitempty"`

	// The name of the provider this resource belongs to. This allows
	// the provider to be interpreted unambiguously in the unusual
	// situation where a provider offers a resource type whose name
	// does not start with its own name, such as the "googlebeta"
	// provider offering "google_compute_instance".
	ProviderName string `json:"provider_name,omitempty"`

	// The version of the resource type schema the "values" property
	// conforms to.
	SchemaVersion uint64 `json:"schema_version,"`

	// The JSON representation of the attribute values of the resource,
	// whose structure depends on the resource type schema. Any unknown
	// values are omitted or set to null, making them indistinguishable
	// from absent values.
	AttributeValues map[string]interface{} `json:"values,omitempty"`

	// The JSON representation of the sensitivity of the resource's
	// attribute values. Only attributes which are sensitive
	// are included in this structure.
	SensitiveValues json.RawMessage `json:"sensitive_values,omitempty"`

	// The addresses of the resources that this resource depends on.
	DependsOn []string `json:"depends_on,omitempty"`

	// If true, the resource has been marked as tainted and will be
	// re-created on the next update.
	Tainted bool `json:"tainted,omitempty"`

	// DeposedKey is set if the resource instance has been marked Deposed and
	// will be destroyed on the next apply.
	DeposedKey string `json:"deposed_key,omitempty"`
}

type ResourceMode string
