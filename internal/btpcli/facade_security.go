package btpcli

func newSecurityFacade(cliClient *v2Client) securityFacade {
	return securityFacade{
		Role:           newSecurityRoleFacade(cliClient),
		RoleCollection: newSecurityRoleCollectionFacade(cliClient),
	}
}

type securityFacade struct {
	Role           securityRoleFacade
	RoleCollection securityRoleCollectionFacade
}
