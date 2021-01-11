package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_Auth(t *testing.T) {
	assert.Equal(t, newAuthService(nil), NewService(nil).Auth())
}
