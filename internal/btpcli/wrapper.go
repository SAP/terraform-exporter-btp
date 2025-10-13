package btpcli

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/testutils"
)

func GetLoggedInClient() (*ClientFacade, error) {
	ctx := context.Background()

	username := os.Getenv("BTP_USERNAME")
	password := os.Getenv("BTP_PASSWORD")
	cliServerUrl := os.Getenv("BTP_CLI_SERVER_URL")
	globalAccount := os.Getenv("BTP_GLOBALACCOUNT")
	idp := os.Getenv("BTP_IDP")
	tlsClientCertificate := os.Getenv("BTP_TLS_CLIENT_CERTIFICATE")
	tlsClientKey := os.Getenv("BTP_TLS_CLIENT_KEY")
	tlsIdpURL := os.Getenv("BTP_TLS_IDP_URL")

	if cliServerUrl == "" {
		cliServerUrl = DefaultServerURL
	}

	u, _ := url.Parse(cliServerUrl)

	client := NewClientFacade(NewV2ClientWithHttpClient(http.DefaultClient, u))

	client.UserAgent = "terraform_exporter_btp"

	if username != "" && password != "" {

		if _, err := client.Login(ctx, NewLoginRequestWithCustomIDP(idp, globalAccount, username, password)); err != nil {
			return nil, fmt.Errorf("error logging in: %w", err)
		}
	}

	if tlsClientCertificate != "" && tlsClientKey != "" && tlsIdpURL != "" {

		passcodeLoginReq := &PasscodeLoginRequest{
			GlobalAccountSubdomain: globalAccount,
			IdentityProvider:       idp,
			IdentityProviderURL:    tlsIdpURL,
			Username:               username,
			PEMEncodedPrivateKey:   tlsClientKey,
			PEMEncodedCertificate:  tlsClientCertificate,
		}

		if _, err := client.PasscodeLogin(ctx, passcodeLoginReq); err != nil {
			return nil, fmt.Errorf("error logging in: %w", err)
		}
	}

	return client, nil
}

func GetGlobalAccountId(client *ClientFacade) (string, error) {

	cliRes, _, err := client.Accounts.GlobalAccount.Get(context.Background())

	if err != nil {
		return "", fmt.Errorf("error getting global account id: %w", err)
	}

	return cliRes.Guid, nil
}

func GetServiceDataByPlanId(client *ClientFacade, subaccountId string, planId string) (planName string, serviceName string, err error) {
	if testing.Testing() {
		planName, serviceName = testutils.GetServiceMockData(planId)
		return planName, serviceName, nil
	}

	cliRes, _, err := client.Services.Plan.GetById(context.Background(), subaccountId, planId)

	if err != nil {
		return "", "", fmt.Errorf("error getting service plan name: %w", err)
	}

	cliRes2, _, err := client.Services.Offering.GetById(context.Background(), subaccountId, cliRes.ServiceOfferingId)

	if err != nil {
		return "", "", fmt.Errorf("error getting service offering name: %w", err)
	}

	return cliRes.Name, cliRes2.Name, nil
}

func GetDefaultRoleCollectionsBySubaccount(subaccountId string, client *ClientFacade) (defaultRoleCollectionNames []string, err error) {
	var roleCollections []string

	cliRes, _, err := client.Security.RoleCollection.ListBySubaccount(context.Background(), subaccountId)

	if err != nil {
		return roleCollections, fmt.Errorf("error listing role collections for subaccount: %w", err)
	}

	for _, roleCollection := range cliRes {
		// The role collections that are marked as IsReadOnly as they are predefined and need not be exported
		if roleCollection.IsReadOnly {
			roleCollections = append(roleCollections, roleCollection.Name)
		}
	}
	return roleCollections, nil
}

func GetDefaultRoleCollectionsByDirectory(directoryId string, client *ClientFacade) (defaultRoleCollectionNames []string, err error) {
	var roleCollections []string

	cliRes, _, err := client.Security.RoleCollection.ListByDirectory(context.Background(), directoryId)

	if err != nil {
		return roleCollections, fmt.Errorf("error listing role collections for directory: %w", err)
	}

	for _, roleCollection := range cliRes {
		// The role collections that are marked as IsReadOnly as they are predefined and need not be exported
		if roleCollection.IsReadOnly {
			roleCollections = append(roleCollections, roleCollection.Name)
		}
	}
	return roleCollections, nil
}

func GetDefaultRolesBySubaccount(subaccountId string, client *ClientFacade) (defaults []string, err error) {
	var roles []string

	cliRes, _, err := client.Security.Role.ListBySubaccount(context.Background(), subaccountId)

	if err != nil {
		return roles, fmt.Errorf("error listing role collections for directory: %w", err)
	}

	for _, role := range cliRes {
		// The roles that are marked as IsReadOnly or are not read only but contain a non-empty attribute list are predefined and need not be exported
		if role.IsReadOnly || (!role.IsReadOnly && len(role.AttributeList) > 0) {
			roles = append(roles, role.Name)
		}
	}
	return roles, nil
}

func GetDefaultRolesByDirectory(directoryId string, client *ClientFacade) (defaultRoles []string, err error) {
	var roles []string

	cliRes, _, err := client.Security.Role.ListByDirectory(context.Background(), directoryId)

	if err != nil {
		return roles, fmt.Errorf("error listing role collections for directory: %w", err)
	}

	for _, role := range cliRes {
		// The roles that are marked as IsReadOnly or are not read only but contain a non-empty attribute list are predefined and need not be exported
		if role.IsReadOnly || (!role.IsReadOnly && len(role.AttributeList) > 0) {
			roles = append(roles, role.Name)
		}
	}
	return roles, nil
}

func GetAppNamesBySubaccountAndApp(subaccountId string, appName string, client *ClientFacade) (technicalAppName string, commercialAppName string, err error) {
	if testing.Testing() {
		return appName, appName, nil
	}

	appList, _, err := client.Accounts.Subscription.List(context.Background(), subaccountId)

	if err != nil {
		return "", "", fmt.Errorf("error getting app names: %w", err)
	}

	for _, app := range appList {
		if app.AppName == appName || app.CommercialAppName == appName {
			technicalAppName = app.AppName
			commercialAppName = app.CommercialAppName
			return technicalAppName, commercialAppName, nil
		}
	}

	return appName, appName, nil
}
