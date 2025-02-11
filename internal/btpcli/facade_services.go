package btpcli

func newServicesFacade(cliClient *v2Client) servicesFacade {
	return servicesFacade{
		Plan:     newServicesPlanFacade(cliClient),
		Offering: newServicesOfferingFacade(cliClient),
	}
}

type servicesFacade struct {
	Plan     servicesPlanFacade
	Offering servicesOfferingFacade
}
