package cfcli

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cloudfoundry/go-cfclient/v3/client"
)

func getAndSetValues() *CloudFoundryClientConfig {

	endpoint := os.Getenv("CF_API_URL")
	user := os.Getenv("CF_USER")
	password := os.Getenv("CF_PASSWORD")
	origin := os.Getenv("CF_ORIGIN")
	cfclientid := os.Getenv("CF_CLIENT_ID")
	cfclientsecret := os.Getenv("CF_CLIENT_SECRET")
	cfaccesstoken := os.Getenv("CF_ACCESS_TOKEN")
	cfrefreshtoken := os.Getenv("CF_REFRESH_TOKEN")

	c := CloudFoundryClientConfig{
		Endpoint:       strings.TrimSuffix(endpoint, "/"),
		User:           user,
		Password:       password,
		CFClientID:     cfclientid,
		CFClientSecret: cfclientsecret,
		Origin:         origin,
		AccessToken:    cfaccesstoken,
		RefreshToken:   cfrefreshtoken,
	}
	return &c
}

func configureClient() (*Session, error) {
	cloudFoundryClientConfig := getAndSetValues()
	session, err := cloudFoundryClientConfig.NewSession()
	if err != nil {

	}
	return session, nil
}

// get all spaces under an org
func GetSpaceList(OrgId string) (map[string]string, error) {
	session, err := configureClient()
	if err != nil {
		return nil, fmt.Errorf("error configure client: %v", err)
	}

	options := client.NewSpaceListOptions()
	options.OrganizationGUIDs = client.Filter{
		Values: []string{
			OrgId,
		},
	}
	spaces, err := session.CFClient.Spaces.ListAll(context.Background(), options)

	if err != nil {
		log.Printf("error listing spaces: %v", err)
		return nil, err
	}

	spaceDetails := make(map[string]string)

	for _, space := range spaces {
		spaceDetails[space.Name] = space.GUID

	}
	return spaceDetails, nil

}

// get username from user id
func GetUser(userId string) (string, error) {
	session, err := configureClient()
	if err != nil {
		return "", fmt.Errorf("error configure client: %s", err)
	}

	userListOption := client.NewUserListOptions()
	userListOption.GUIDs = client.Filter{
		Values: []string{
			userId,
		},
	}
	user, err := session.CFClient.Users.ListAll(context.Background(), userListOption)
	if err != nil {
		return "", err
	}
	if len(user) == 0 {
		return "", fmt.Errorf("user not found")
	}
	return user[0].Username, nil
}

// get space name from space id
func GetSpaceName(spaceId string) (string, error) {
	session, err := configureClient()
	if err != nil {
		return "", fmt.Errorf("error configure client")
	}

	spaceListOption := client.NewSpaceListOptions()
	spaceListOption.GUIDs = client.Filter{
		Values: []string{
			spaceId,
		},
	}

	space, err := session.CFClient.Spaces.ListAll(context.Background(), spaceListOption)
	if err != nil {
		return "", err
	}
	if len(space) == 0 {
		return "", fmt.Errorf("space not found")
	}
	return space[0].Name, nil
}
