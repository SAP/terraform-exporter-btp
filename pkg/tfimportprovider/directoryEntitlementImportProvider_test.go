package tfimportprovider

import (
	"testing"

	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/stretchr/testify/assert"
)

func TestCreateDirectoryEntitlementImportBlock(t *testing.T) {
	resourceDoc := tfutils.EntityDocs{
		Import: "import {\n\t\t\t\tto = btp_directory_entitlement.<resource_name>\n\t\t\t\tid = \"<directory_id>,<service_name>,<plan_name>\"\n\t\t\t  }\n",
	}

	jsonString := "{\"alert-notification:build-code\": {\"category\":\"SERVICE\",\"plan_description\":\"This is a dedicated plan for the SAP Build Code. If you want to use this service with the SAP Build Code, use the build-code plan. Please note if you use the regular plan, it will incur charges outside the build code capacity units metric. Allows production & consumption of custom events\",\"plan_display_name\":\"build-code\",\"plan_name\":\"build-code\",\"quota_assigned\": 1,\"quota_remaining\": 1,\"service_display_name\":\"Alert Notification\",\"service_name\":\"alert-notification\"}}"
	dataEntitlement, _ := GetDataFromJsonString(jsonString)

	jsonStringMultipleEntitlements := "	{\"alert-notification:build-code\": {\"category\":\"SERVICE\",\"plan_description\":\"This is a dedicated plan for the SAP Build Code. If you want to use this service with the SAP Build Code, use the build-code plan. Please note if you use the regular plan, it will incur charges outside the build code capacity units metric. Allows production & consumption of custom events\",\"plan_display_name\":\"build-code\",\"plan_name\":\"build-code\",\"quota_assigned\": 1,\"quota_remaining\": 1,\"service_display_name\":\"Alert Notification\",\"service_name\":\"alert-notification\"},\"alert-notification:standard\": {\"category\":\"ELASTIC_SERVICE\",\"plan_description\":\"Allows production & consumption of custom events\",\"plan_display_name\":\"standard\",\"plan_name\":\"standard\",\"quota_assigned\": 13,\"quota_remaining\": 13,\"service_display_name\":\"Alert Notification\",\"service_name\":\"alert-notification\"},\"auditlog:standard\": {\"category\":\"SERVICE\",\"plan_description\":\"Technical plan to enable the possibility of ingestion of audit data from SAP BTP Applications and services. The standard plan uses basic authentication to authenticate in front of the audit log server, which has some security vulnerabilities identified and is about to be deprecated and not recommended to be used.\",\"plan_display_name\":\"standard\",\"plan_name\":\"standard\",\"quota_assigned\": 1,\"quota_remaining\": 1,\"service_display_name\":\"Audit Log Service\",\"service_name\":\"auditlog\"}}"
	dataMultipleEntitlements, _ := GetDataFromJsonString(jsonStringMultipleEntitlements)

	tests := []struct {
		name          string
		data          map[string]interface{}
		directoryId   string
		filterValues  []string
		expectedBlock string
		expectError   bool
	}{
		{
			name:          "No filter values",
			data:          dataEntitlement,
			directoryId:   "12345",
			expectedBlock: "import {\n\t\t\t\tto = btp_directory_entitlement.alert-notification_build-code\n\t\t\t\tid = \"12345,alert-notification,build-code\"\n\t\t\t  }\n\n",
			expectError:   false,
		},
		{
			name:          "With filter values",
			data:          dataMultipleEntitlements,
			directoryId:   "12345",
			filterValues:  []string{"alert-notification_build-code"},
			expectedBlock: "import {\n\t\t\t\tto = btp_directory_entitlement.alert-notification_build-code\n\t\t\t\tid = \"12345,alert-notification,build-code\"\n\t\t\t  }\n\n",
			expectError:   false,
		},
		{
			name:          "Missing entitlement in filter values",
			data:          dataMultipleEntitlements,
			directoryId:   "12345",
			filterValues:  []string{"non-existent-entitlement"},
			expectedBlock: "",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			importBlock, err := CreateDirEntitlementImportBlock(tt.data, tt.directoryId, tt.filterValues, resourceDoc)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBlock, importBlock)
			}
		})
	}
}
