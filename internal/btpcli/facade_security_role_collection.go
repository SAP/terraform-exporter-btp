package btpcli

import (
	"context"

	"github.com/SAP/terraform-exporter-btp/internal/btpcli/types/xsuaa_authz"
)

func newSecurityRoleCollectionFacade(cliClient *v2Client) securityRoleCollectionFacade {
	return securityRoleCollectionFacade{cliClient: cliClient}
}

type securityRoleCollectionFacade struct {
	cliClient *v2Client
}

func (f *securityRoleCollectionFacade) getCommand() string {
	return "security/role-collection"
}

func (f *securityRoleCollectionFacade) ListByGlobalAccount(ctx context.Context) ([]xsuaa_authz.RoleCollection, CommandResponse, error) {
	return doExecute[[]xsuaa_authz.RoleCollection](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
	}))
}

func (f *securityRoleCollectionFacade) GetByGlobalAccount(ctx context.Context, roleCollectionName string) (xsuaa_authz.RoleCollection, CommandResponse, error) {
	return doExecute[xsuaa_authz.RoleCollection](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"globalAccount":      f.cliClient.GetGlobalAccountSubdomain(),
		"roleCollectionName": roleCollectionName,
	}))
}

func (f *securityRoleCollectionFacade) ListBySubaccount(ctx context.Context, subaccountId string) ([]xsuaa_authz.RoleCollection, CommandResponse, error) {
	return doExecute[[]xsuaa_authz.RoleCollection](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
	}))
}

func (f *securityRoleCollectionFacade) GetBySubaccount(ctx context.Context, subaccountId string, roleCollectionName string) (xsuaa_authz.RoleCollection, CommandResponse, error) {
	return doExecute[xsuaa_authz.RoleCollection](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"subaccount":         subaccountId,
		"roleCollectionName": roleCollectionName,
	}))
}

func (f *securityRoleCollectionFacade) ListByDirectory(ctx context.Context, directoryId string) ([]xsuaa_authz.RoleCollection, CommandResponse, error) {
	return doExecute[[]xsuaa_authz.RoleCollection](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"directory": directoryId,
	}))
}

func (f *securityRoleCollectionFacade) GetByDirectory(ctx context.Context, directoryId string, roleCollectionName string) (xsuaa_authz.RoleCollection, CommandResponse, error) {
	return doExecute[xsuaa_authz.RoleCollection](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"directory":          directoryId,
		"roleCollectionName": roleCollectionName,
	}))
}
