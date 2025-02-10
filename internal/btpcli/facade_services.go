package btpcli

func newServicesFacade(cliClient *v2Client) servicesFacade {
	return servicesFacade{
		Plan: newServicesPlanFacade(cliClient),
	}
}

type servicesFacade struct {
	Plan servicesPlanFacade
}
