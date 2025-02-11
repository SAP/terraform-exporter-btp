package cfcli

import (
	"github.com/cloudfoundry/go-cfclient/v3/client"
	config "github.com/cloudfoundry/go-cfclient/v3/config"
)

type CloudFoundryClientConfig struct {
	Endpoint       string
	User           string
	Password       string
	CFClientID     string
	CFClientSecret string
	Origin         string
	AccessToken    string
	RefreshToken   string
}

type Session struct {
	CFClient *client.Client
}

func (c *CloudFoundryClientConfig) NewSession() (*Session, error) {
	var cfg *config.Config
	var err error
	var opts []config.Option

	switch {
	case c.User != "" && c.Password != "":
		opts = append(opts, config.UserPassword(c.User, c.Password))
		if c.Origin != "" {
			opts = append(opts, config.Origin(c.Origin))
		}
		cfg, err = config.New(c.Endpoint, opts...)
	case c.CFClientID != "" && c.CFClientSecret != "":
		opts = append(opts, config.ClientCredentials(c.CFClientID, c.CFClientSecret))
		cfg, err = config.New(c.Endpoint, opts...)
	case c.AccessToken != "":
		opts = append(opts, config.Token(c.AccessToken, c.RefreshToken))
		cfg, err = config.New(c.Endpoint, opts...)
	default:
		cfg, err = config.NewFromCFHome(opts...)
	}
	if err != nil {
		return nil, err
	}
	cf, err := client.New(cfg)
	if err != nil {
		return nil, err
	}
	s := Session{
		CFClient: cf,
	}
	return &s, nil
}
