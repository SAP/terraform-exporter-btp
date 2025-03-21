package btpcli

import (
	"context"

	"github.com/SAP/terraform-exporter-btp/internal/btpcli/types/xsuaa_authz"
)

func newSecurityRoleFacade(cliClient *v2Client) securityRoleFacade {
	return securityRoleFacade{cliClient: cliClient}
}

type securityRoleFacade struct {
	cliClient *v2Client
}

func (f *securityRoleFacade) getCommand() string {
	return "security/role"
}

func (f *securityRoleFacade) ListByGlobalAccount(ctx context.Context) ([]xsuaa_authz.Role, CommandResponse, error) {
	return doExecute[[]xsuaa_authz.Role](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
	}))
}

func (f *securityRoleFacade) GetByGlobalAccount(ctx context.Context, roleName string, roleTemplateAppId string, roleTemplateName string) (xsuaa_authz.Role, CommandResponse, error) {
	return doExecute[xsuaa_authz.Role](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"globalAccount":    f.cliClient.GetGlobalAccountSubdomain(),
		"roleName":         roleName,
		"appId":            roleTemplateAppId,
		"roleTemplateName": roleTemplateName,
	}))
}

func (f *securityRoleFacade) ListBySubaccount(ctx context.Context, subaccountId string) ([]xsuaa_authz.Role, CommandResponse, error) {
	return doExecute[[]xsuaa_authz.Role](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
	}))
}

func (f *securityRoleFacade) GetBySubaccount(ctx context.Context, subaccountId string, roleName string, roleTemplateAppId string, roleTemplateName string) (xsuaa_authz.Role, CommandResponse, error) {
	return doExecute[xsuaa_authz.Role](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"subaccount":       subaccountId,
		"roleName":         roleName,
		"appId":            roleTemplateAppId,
		"roleTemplateName": roleTemplateName,
	}))
}

func (f *securityRoleFacade) ListByDirectory(ctx context.Context, directoryId string) ([]xsuaa_authz.Role, CommandResponse, error) {
	return doExecute[[]xsuaa_authz.Role](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"directory": directoryId,
	}))
}

func (f *securityRoleFacade) GetByDirectory(ctx context.Context, directoryId string, roleName string, roleTemplateAppId string, roleTemplateName string) (xsuaa_authz.Role, CommandResponse, error) {
	return doExecute[xsuaa_authz.Role](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"directory":        directoryId,
		"roleName":         roleName,
		"appId":            roleTemplateAppId,
		"roleTemplateName": roleTemplateName,
	}))
}
