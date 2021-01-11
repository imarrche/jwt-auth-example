package pg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore_Get(t *testing.T) {
	assert.Equal(t, &Store{}, Get(nil))
}

func TestStore_Users(t *testing.T) {
	assert.Equal(t, newUserRepo(nil), Get(nil).Users())
}
