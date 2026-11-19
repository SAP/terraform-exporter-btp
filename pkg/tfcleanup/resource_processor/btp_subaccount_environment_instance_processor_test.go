package resourceprocessor

import (
	"fmt"
	"testing"

	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvironmentInstancePostProcessing(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		blockType   string
		resource    string
		parameters  string
		want        string
	}{
		{
			name: "exported Cloud Foundry instance",
			parameters: `jsonencode({
  instance_name = "dev-inttest-exporter-8b2967c9-c3f3-68a4-4604-16702098e1d1"
  status = "ACTIVE"
})`,
			want: `jsonencode({
  instance_name = "dev-inttest-exporter-8b2967c9-c3f3-68a4-4604-16702098e1d1"
})`,
		},
		{name: "inactive status", parameters: `jsonencode({ status = "INACTIVE", instance_name = "cf" })`, want: `jsonencode({
  instance_name = "cf"
})`},
		{name: "null status only", parameters: `jsonencode({ status = null })`},
		{name: "referenced status only", parameters: `jsonencode({ status = var.status })`},
		{name: "object status only", parameters: `jsonencode({ status = { value = "ACTIVE" } })`},
		{name: "quoted status key", parameters: `jsonencode({ "status" = "ACTIVE", "instance_name" = "cf" })`, want: `jsonencode({
  "instance_name" = "cf"
})`},
		{
			name: "preserve nested objects and expressions",
			parameters: `jsonencode({
  status = "ACTIVE"
  instance_name = var.instance_name
  metadata = { status = "keep", tags = ["Grüße", "ACTIVE"] }
  message = "hello ${var.name}"
  enabled = true
  count = 2
})`,
			want: `jsonencode({
  instance_name = var.instance_name
  metadata = { status = "keep", tags = ["Grüße", "ACTIVE"] }
  message = "hello ${var.name}"
  enabled = true
  count = 2
})`,
		},
		{name: "missing status", parameters: `jsonencode({ instance_name = "cf" })`, want: `jsonencode({ instance_name = "cf" })`},
		{name: "empty parameters", parameters: `jsonencode({})`},
		{name: "missing parameters"},
		{name: "other environment", environment: `"kyma"`, parameters: `jsonencode({ status = "ACTIVE" })`, want: `jsonencode({ status = "ACTIVE" })`},
		{name: "missing environment type", environment: "omit", parameters: `jsonencode({ status = "ACTIVE" })`, want: `jsonencode({ status = "ACTIVE" })`},
		{name: "unknown environment type", environment: "var.environment_type", parameters: `jsonencode({ status = "ACTIVE" })`, want: `jsonencode({ status = "ACTIVE" })`},
		{name: "other resource", resource: "btp_subaccount_service_binding", parameters: `jsonencode({ status = "ACTIVE" })`, want: `jsonencode({ status = "ACTIVE" })`},
		{name: "data source", blockType: "data", parameters: `jsonencode({ status = "ACTIVE" })`, want: `jsonencode({ status = "ACTIVE" })`},
		{name: "referenced parameters", parameters: "var.parameters", want: "var.parameters"},
		{name: "computed parameters", parameters: `jsonencode(var.parameters)`, want: `jsonencode(var.parameters)`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			environment := tt.environment
			if environment == "" {
				environment = `"cloudfoundry"`
			}
			blockType := tt.blockType
			if blockType == "" {
				blockType = "resource"
			}
			resource := tt.resource
			if resource == "" {
				resource = tfutils.SubaccountEnvironmentInstanceType
			}
			config := func(parameters, subaccount string) string {
				source := fmt.Sprintf(`%s "%s" "instance" {
  provider = btp
  name = "cf-instance"
  landscape_label = "cf-us10"
  plan_name = "standard"
  service_name = "cloudfoundry"
  subaccount_id = %s
`, blockType, resource, subaccount)
				if environment != "omit" {
					source += "  environment_type = " + environment + "\n"
				}
				if parameters != "" {
					source += "  parameters = " + parameters + "\n"
				}
				return source + "}\n"
			}
			file, diagnostics := hclwrite.ParseConfig([]byte(config(tt.parameters, `"sa"`)), "test.tf", hcl.InitialPos)
			require.False(t, diagnostics.HasErrors(), diagnostics.Error())
			variables := generictools.VariableContent{}
			dependencies := generictools.NewDependencyAddresses()
			dependencies.SubaccountAddress = "btp_subaccount.export_0"
			ProcessResources(file, tfutils.SubaccountLevel, &variables, &dependencies, nil, generictools.LevelIds{SubaccountId: "sa"})
			want := hclwrite.Format([]byte(config(tt.want, "btp_subaccount.export_0.id")))
			assert.Equal(t, string(want), string(file.Bytes()))
			_, diagnostics = hclwrite.ParseConfig(file.Bytes(), "result.tf", hcl.InitialPos)
			require.False(t, diagnostics.HasErrors(), diagnostics.Error())
			before := string(file.Bytes())
			ProcessResources(file, tfutils.SubaccountLevel, &variables, &dependencies, nil, generictools.LevelIds{SubaccountId: "sa"})
			assert.Equal(t, before, string(file.Bytes()), "cleanup must be idempotent")
		})
	}
}
