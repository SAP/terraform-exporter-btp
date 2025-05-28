package btpcli

import (
	"context"

	"github.com/SAP/terraform-exporter-btp/internal/btpcli/types/saas_manager_service"
)

func newAccountsSubscriptionFacade(cliClient *v2Client) accountsSubscriptionFacade {
	return accountsSubscriptionFacade{cliClient: cliClient}
}

type accountsSubscriptionFacade struct {
	cliClient *v2Client
}

func (f *accountsSubscriptionFacade) getCommand() string {
	return "accounts/subscription"
}

func (f *accountsSubscriptionFacade) List(ctx context.Context, subaccountId string) ([]saas_manager_service.EntitledApplicationsResponseObject, CommandResponse, error) {
	type wrapper struct { // TODO should be in types package
		Applications []saas_manager_service.EntitledApplicationsResponseObject `json:"applications"`
	}

	data, res, err := doExecute[wrapper](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
	}))

	return data.Applications, res, err
}
