package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextKeyString(t *testing.T) {
	t.Parallel()

	testCtxKey := &contextKey{"test"}
	assert.Equal(t, "bastion/middleware context value test", testCtxKey.String())
}
