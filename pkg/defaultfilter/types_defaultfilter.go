package defaultfilter

type EntitlementFilterData []struct {
	ServiceName string
	PlanName    string
}

var DefaultEntitlements EntitlementFilterData = EntitlementFilterData{
	{
		ServiceName: "auditlog-api",
		PlanName:    "default",
	},
	{
		ServiceName: "autoscaler",
		PlanName:    "standard",
	},
	{
		ServiceName: "application-logs",
		PlanName:    "lite",
	},
	{
		ServiceName: "auditlog",
		PlanName:    "standard",
	},
	{
		ServiceName: "cias",
		PlanName:    "standard",
	},
	{
		ServiceName: "cias",
		PlanName:    "oauth2",
	},
	{
		ServiceName: "connectivity",
		PlanName:    "lite",
	},
	{
		ServiceName: "connectivity",
		PlanName:    "connectivity_proxy",
	},
	{
		ServiceName: "content-agent",
		PlanName:    "free",
	},
	{
		ServiceName: "credstore",
		PlanName:    "proxy",
	},
	{
		ServiceName: "destination",
		PlanName:    "lite",
	},
	{
		ServiceName: "feature-flags-dashboard",
		PlanName:    "dashboard",
	},
	{
		ServiceName: "feature-flags",
		PlanName:    "lite",
	},
	{
		ServiceName: "feature-flags",
		PlanName:    "standard",
	},
	{
		ServiceName: "html5-apps-repo",
		PlanName:    "app-host",
	},
	{
		ServiceName: "html5-apps-repo",
		PlanName:    "app-runtime",
	},
	{
		ServiceName: "one-mds",
		PlanName:    "sap-integration",
	},
	{
		ServiceName: "mdo-one-mds-master",
		PlanName:    "standard",
	},
	{
		ServiceName: "print",
		PlanName:    "receiver",
	},
	{
		ServiceName: "saas-registry",
		PlanName:    "application",
	},
	{
		ServiceName: "sap-identity-services-onboarding",
		PlanName:    "default",
	},
	{
		ServiceName: "identity",
		PlanName:    "application",
	},
	{
		ServiceName: "service-manager",
		PlanName:    "global-offerings-audit",
	},
	{
		ServiceName: "service-manager",
		PlanName:    "subaccount-admin",
	},
	{
		ServiceName: "service-manager",
		PlanName:    "subaccount-audit",
	},
	{
		ServiceName: "service-manager",
		PlanName:    "container",
	},
	{
		ServiceName: "service-manager",
		PlanName:    "service-operator-access",
	},
	{
		ServiceName: "content-agent",
		PlanName:    "standard",
	},
	{
		ServiceName: "content-agent",
		PlanName:    "application",
	},
	{
		ServiceName: "one-mds-master",
		PlanName:    "standard",
	},
	{
		ServiceName: "uas",
		PlanName:    "reporting-directory",
	},
	{
		ServiceName: "xsuaa",
		PlanName:    "apiaccess",
	},
	{
		ServiceName: "xsuaa",
		PlanName:    "application",
	},
	{
		ServiceName: "xsuaa",
		PlanName:    "space",
	},
	{
		ServiceName: "xsuaa",
		PlanName:    "broker",
	},
}
