package tfimportprovider

import (
	"testing"

	"github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/stretchr/testify/assert"
)

func TestCreateSubscriptionImportBlock(t *testing.T) {
	resourceDoc := tfutils.EntityDocs{
		Import: "import {\n\t\t\t\tto = btp_subaccount_subscription.<resource_name>\n\t\t\t\tid = \"<subaccount_id>,<app_name>,<plan_name>\"\n\t\t\t  }\n",
	}

	jsonStringSingleSubscription := "{\"id\":\"1234\",\"subaccount_id\":\"1234\",\"values\": [{\"additional_plan_features\": null,\"app_id\":\"feature-flags!b1765\",\"app_name\":\"feature-flags-dashboard\",\"authentication_provider\":\"XSUAA\",\"automation_state\": null,\"automation_state_message\": null,\"category\":\"Foundation / Cross Services\",\"category_display_name\": null,\"commercial_app_name\":\"feature-flags-dashboard\",\"created_date\":\"2024-09-16T08:06:51Z\",\"customer_developed\": false,\"description\":\"The Feature Flags service allows you to enable or disable new features at runtime without redeploying or restarting the application. You can use feature flags to control code delivery, synchronized rollout, direct shipment, and fast rollback of features.\",\"display_name\":\"Feature Flags Service\",\"formation_solution_name\":\"\",\"globalaccount_id\":\"1234\",\"id\":\"5e6e5ce6-1847-4193-b6e7-f6d3f8a4339f\",\"incident_tracking_component\": null,\"labels\": null,\"last_modified\":\"2024-09-25T11:03:37Z\",\"plan_description\": null,\"plan_name\":\"dashboard\",\"platform_entity_id\":\"\",\"quota\": 1,\"short_description\": null,\"state\":\"SUBSCRIBED\",\"subscribed_subaccount_id\":\"1234\",\"subscribed_tenant_id\":\"1234\",\"subscription_url\":\"https://btptfexporter-validate-xhyu0cme.feature-flags-dashboard.cfapps.us10.hana.ondemand.com\",\"supports_parameters_updates\": false,\"supports_plan_updates\": false,\"tenant_id\":\"faefaf02-0b2f-48e0-87e7-2d89c121dc25\"}]}"
	dataSingleSubscription, _ := GetDataFromJsonString(jsonStringSingleSubscription)

	jsonStringDoubleSubscription := "{\"id\":\"1234\",\"subaccount_id\":\"1234\",\"values\": [{\"additional_plan_features\": [\"content management\"],\"app_id\":\"cas-ui-xsuaa-prod!t6249\",\"app_name\":\"content-agent-ui\",\"authentication_provider\":\"XSUAA\",\"automation_state\": null,\"automation_state_message\": null,\"category\":\"Integration Suite\",\"category_display_name\": null,\"commercial_app_name\":\"content-agent-ui\",\"created_date\":\"2024-09-25T12:15:36Z\",\"customer_developed\": false,\"description\":\"SAP Content Agent service is a tool for SAP BTP applications offering generic content management operations such as view, export and import content with inter-dependencies and integration with SAP Cloud Transport Management service. It offers a view to track all activities along with logs, status and other information.\",\"display_name\":\"Content Agent Service\",\"formation_solution_name\":\"\",\"globalaccount_id\":\"07b521f7-c0f7-4306-a95a-0adefdb86d23\",\"id\":\"afdf5081-9e84-47e8-ac2a-1298d347155d\",\"incident_tracking_component\": null,\"labels\": null,\"last_modified\":\"2024-09-25T12:15:39Z\",\"plan_description\": null,\"plan_name\":\"free\",\"platform_entity_id\":\"\",\"quota\": 1,\"short_description\": null,\"state\":\"SUBSCRIBED\",\"subscribed_subaccount_id\":\"1234\",\"subscribed_tenant_id\":\"1234\",\"subscription_url\":\"https://btptfexporter-validate-xhyu0cme.us10.content-agent.cloud.sap\",\"supports_parameters_updates\": false,\"supports_plan_updates\": false,\"tenant_id\":\"50107178-c526-498e-b54c-577312a5f79b\"},    {\"additional_plan_features\": null,\"app_id\":\"feature-flags!b1765\",\"app_name\":\"feature-flags-dashboard\",\"authentication_provider\":\"XSUAA\",\"automation_state\": null,\"automation_state_message\": null,\"category\":\"Foundation / Cross Services\",\"category_display_name\": null,\"commercial_app_name\":\"feature-flags-dashboard\",\"created_date\":\"2024-09-16T08:06:51Z\",\"customer_developed\": false,\"description\":\"The Feature Flags service allows you to enable or disable new features at runtime without redeploying or restarting the application. You can use feature flags to control code delivery, synchronized rollout, direct shipment, and fast rollback of features.\",\"display_name\":\"Feature Flags Service\",\"formation_solution_name\":\"\",\"globalaccount_id\":\"1234\",\"id\":\"5e6e5ce6-1847-4193-b6e7-f6d3f8a4339f\",\"incident_tracking_component\": null,\"labels\": null,\"last_modified\":\"2024-09-25T11:03:37Z\",\"plan_description\": null,\"plan_name\":\"dashboard\",\"platform_entity_id\":\"\",\"quota\": 1,\"short_description\": null,\"state\":\"SUBSCRIBED\",\"subscribed_subaccount_id\":\"1234\",\"subscribed_tenant_id\":\"1234\",\"subscription_url\":\"https://btptfexporter-validate-xhyu0cme.feature-flags-dashboard.cfapps.us10.hana.ondemand.com\",\"supports_parameters_updates\": false,\"supports_plan_updates\": false,\"tenant_id\":\"faefaf02-0b2f-48e0-87e7-2d89c121dc25\"}]}"
	dataDoubleSubscription, _ := GetDataFromJsonString(jsonStringDoubleSubscription)

	jsonStringSubscriptionInProcess := "{\"id\":\"1234\",\"subaccount_id\":\"1234\",\"values\": [{\"additional_plan_features\": null,\"app_id\":\"feature-flags!b1765\",\"app_name\":\"feature-flags-dashboard\",\"authentication_provider\":\"XSUAA\",\"automation_state\": null,\"automation_state_message\": null,\"category\":\"Foundation / Cross Services\",\"category_display_name\": null,\"commercial_app_name\":\"feature-flags-dashboard\",\"created_date\":\"2024-09-16T08:06:51Z\",\"customer_developed\": false,\"description\":\"The Feature Flags service allows you to enable or disable new features at runtime without redeploying or restarting the application. You can use feature flags to control code delivery, synchronized rollout, direct shipment, and fast rollback of features.\",\"display_name\":\"Feature Flags Service\",\"formation_solution_name\":\"\",\"globalaccount_id\":\"1234\",\"id\":\"5e6e5ce6-1847-4193-b6e7-f6d3f8a4339f\",\"incident_tracking_component\": null,\"labels\": null,\"last_modified\":\"2024-09-25T11:03:37Z\",\"plan_description\": null,\"plan_name\":\"dashboard\",\"platform_entity_id\":\"\",\"quota\": 1,\"short_description\": null,\"state\":\"SUBSCRIBE_FAILED\",\"subscribed_subaccount_id\":\"1234\",\"subscribed_tenant_id\":\"1234\",\"subscription_url\":\"https://btptfexporter-validate-xhyu0cme.feature-flags-dashboard.cfapps.us10.hana.ondemand.com\",\"supports_parameters_updates\": false,\"supports_plan_updates\": false,\"tenant_id\":\"faefaf02-0b2f-48e0-87e7-2d89c121dc25\"}]}"
	dataSubscriptionInProcess, _ := GetDataFromJsonString(jsonStringSubscriptionInProcess)

	jsonStringSubscriptionFailed := "{\"id\":\"1234\",\"subaccount_id\":\"1234\",\"values\": [{\"additional_plan_features\": null,\"app_id\":\"feature-flags!b1765\",\"app_name\":\"feature-flags-dashboard\",\"authentication_provider\":\"XSUAA\",\"automation_state\": null,\"automation_state_message\": null,\"category\":\"Foundation / Cross Services\",\"category_display_name\": null,\"commercial_app_name\":\"feature-flags-dashboard\",\"created_date\":\"2024-09-16T08:06:51Z\",\"customer_developed\": false,\"description\":\"The Feature Flags service allows you to enable or disable new features at runtime without redeploying or restarting the application. You can use feature flags to control code delivery, synchronized rollout, direct shipment, and fast rollback of features.\",\"display_name\":\"Feature Flags Service\",\"formation_solution_name\":\"\",\"globalaccount_id\":\"1234\",\"id\":\"5e6e5ce6-1847-4193-b6e7-f6d3f8a4339f\",\"incident_tracking_component\": null,\"labels\": null,\"last_modified\":\"2024-09-25T11:03:37Z\",\"plan_description\": null,\"plan_name\":\"dashboard\",\"platform_entity_id\":\"\",\"quota\": 1,\"short_description\": null,\"state\":\"IN_PROCESS\",\"subscribed_subaccount_id\":\"1234\",\"subscribed_tenant_id\":\"1234\",\"subscription_url\":\"https://btptfexporter-validate-xhyu0cme.feature-flags-dashboard.cfapps.us10.hana.ondemand.com\",\"supports_parameters_updates\": false,\"supports_plan_updates\": false,\"tenant_id\":\"faefaf02-0b2f-48e0-87e7-2d89c121dc25\"}]}"
	dataSubscriptionFailed, _ := GetDataFromJsonString(jsonStringSubscriptionFailed)

	tests := []struct {
		name          string
		data          map[string]interface{}
		subaccountId  string
		filterValues  []string
		expectedBlock string
		expectedCount int
		expectError   bool
	}{
		{
			name:          "Single subscription",
			data:          dataSingleSubscription,
			subaccountId:  "1234",
			expectedBlock: "import {\n\t\t\t\tto = btp_subaccount_subscription.feature_flags_dashboard\n\t\t\t\tid = \"1234,feature-flags-dashboard,dashboard\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Subscription with filter values",
			data:          dataDoubleSubscription,
			subaccountId:  "1234",
			filterValues:  []string{"feature-flags-dashboard_dashboard"},
			expectedBlock: "import {\n\t\t\t\tto = btp_subaccount_subscription.feature_flags_dashboard\n\t\t\t\tid = \"1234,feature-flags-dashboard,dashboard\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Subscription not found in filter values",
			data:          dataSingleSubscription,
			subaccountId:  "1234",
			filterValues:  []string{"nonexistentapp_nonexistentplan"},
			expectedBlock: "",
			expectedCount: 0,
			expectError:   true,
		},
		{
			name:          "Subscription in progress",
			data:          dataSubscriptionInProcess,
			subaccountId:  "1234",
			expectedBlock: "",
			expectedCount: 0,
			expectError:   false,
		},
		{
			name:          "Subscription failed",
			data:          dataSubscriptionFailed,
			subaccountId:  "1234",
			expectedBlock: "",
			expectedCount: 0,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			importBlock, count, err := createSubscriptionImportBlock(tt.data, tt.subaccountId, tt.filterValues, resourceDoc)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBlock, importBlock)
				assert.Equal(t, tt.expectedCount, count)
			}
		})
	}
}
