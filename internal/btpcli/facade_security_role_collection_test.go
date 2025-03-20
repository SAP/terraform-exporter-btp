package btpcli

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityRoleCollectionFacade_ListByGlobalAccount(t *testing.T) {
	command := "security/role-collection"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionList, map[string]string{
				"globalAccount": "795b53bb-a3f0-4769-adf0-26173282a975",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.ListByGlobalAccount(context.TODO())

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_ListBySubaccount(t *testing.T) {
	command := "security/role-collection"

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

		_, res, err := uut.Security.RoleCollection.ListBySubaccount(context.TODO(), subaccountId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_ListByDirectory(t *testing.T) {
	command := "security/role-collection"

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

		_, res, err := uut.Security.RoleCollection.ListByDirectory(context.TODO(), directoryId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_GetByGlobalAccount(t *testing.T) {
	command := "security/role-collection"

	globalAccountId := "795b53bb-a3f0-4769-adf0-26173282a975"
	roleCollectionName := "Global Account Administrator"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"globalAccount":      globalAccountId,
				"roleCollectionName": roleCollectionName,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.GetByGlobalAccount(context.TODO(), roleCollectionName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_GetBySubaccount(t *testing.T) {
	command := "security/role-collection"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	roleCollectionName := "Subaccount Administrator"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"subaccount":         subaccountId,
				"roleCollectionName": roleCollectionName,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.GetBySubaccount(context.TODO(), subaccountId, roleCollectionName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_GetByDirectory(t *testing.T) {
	command := "security/role-collection"

	directoryId := "f6c7137d-c5a0-48c2-b2a4-fd64e6b35d3d"
	roleCollectionName := "Directory Administrator"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"directory":          directoryId,
				"roleCollectionName": roleCollectionName,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.GetByDirectory(context.TODO(), directoryId, roleCollectionName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}
