package btpcli

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountsSubscriptionFacade_List(t *testing.T) {
	command := "accounts/subscription"

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

		_, res, err := uut.Accounts.Subscription.List(context.TODO(), subaccountId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}
