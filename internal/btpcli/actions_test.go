package btpcli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertAction(t *testing.T, expectAction Action, initFn func(command string, args any) *CommandRequest) {
	t.Helper()

	assert.Equal(t, expectAction, initFn("", map[string]string{}).Action)
}

func TestNewGetRequest(t *testing.T) {
	assertAction(t, ActionGet, NewGetRequest)
}

func TestNewListRequest(t *testing.T) {
	assertAction(t, ActionList, NewListRequest)
}
