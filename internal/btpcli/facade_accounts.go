package btpcli

func newAccountsFacade(cliClient *v2Client) accountsFacade {
	return accountsFacade{
		GlobalAccount: newAccountsGlobalAccountFacade(cliClient),
	}
}

type accountsFacade struct {
	GlobalAccount accountsGlobalAccountFacade
}
