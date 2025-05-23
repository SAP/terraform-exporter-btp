package btpcli

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityRoleFacade_ListByGlobalAccount(t *testing.T) {
	command := "security/role"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionList, map[string]string{
				"globalAccount": "795b53bb-a3f0-4769-adf0-26173282a975",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Role.ListByGlobalAccount(context.TODO())

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleFacade_ListBySubaccount(t *testing.T) {
	command := "security/role"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionList, map[string]string{
				"subaccount": subaccountId,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Role.ListBySubaccount(context.TODO(), subaccountId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleFacade_ListByDirectory(t *testing.T) {
	command := "security/role"

	directoryId := "f6c7137d-c5a0-48c2-b2a4-fd64e6b35d3d"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionList, map[string]string{
				"directory": directoryId,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Role.ListByDirectory(context.TODO(), directoryId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleFacade_GetByGlobalAccount(t *testing.T) {
	command := "security/role"

	globalAccountId := "795b53bb-a3f0-4769-adf0-26173282a975"
	roleName := "User and Role Auditor"
	roleTemplateAppId := "xsuaa!t1"
	roleTemplateName := "xsuaa_auditor"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"globalAccount":    globalAccountId,
				"appId":            roleTemplateAppId,
				"roleName":         roleName,
				"roleTemplateName": roleTemplateName,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Role.GetByGlobalAccount(context.TODO(), roleName, roleTemplateAppId, roleTemplateName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleFacade_GetBySubaccount(t *testing.T) {
	command := "security/role"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	roleName := "User and Role Auditor"
	roleTemplateAppId := "xsuaa!t1"
	roleTemplateName := "xsuaa_auditor"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"subaccount":       subaccountId,
				"appId":            roleTemplateAppId,
				"roleName":         roleName,
				"roleTemplateName": roleTemplateName,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Role.GetBySubaccount(context.TODO(), subaccountId, roleName, roleTemplateAppId, roleTemplateName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleFacade_GetByDirectory(t *testing.T) {
	command := "security/role"

	directoryId := "f6c7137d-c5a0-48c2-b2a4-fd64e6b35d3d"
	roleName := "User and Role Auditor"
	roleTemplateAppId := "xsuaa!t1"
	roleTemplateName := "xsuaa_auditor"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"directory":        directoryId,
				"appId":            roleTemplateAppId,
				"roleName":         roleName,
				"roleTemplateName": roleTemplateName,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Role.GetByDirectory(context.TODO(), directoryId, roleName, roleTemplateAppId, roleTemplateName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}
