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
