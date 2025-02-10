package btpcli

func NewClientFacade(cliClient *v2Client) *ClientFacade {
	return &ClientFacade{
		v2Client: cliClient,
		Accounts: newAccountsFacade(cliClient),
		Services: newServicesFacade(cliClient),
	}
}

type ClientFacade struct {
	*v2Client
	Accounts accountsFacade
	Services servicesFacade
}
