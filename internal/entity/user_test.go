package entity

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("Matheus", "mh.rossetti2002@gmail.com", "123")
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.Id)
	assert.NotEmpty(t, user.PasswordHash)
	assert.Equal(t, "Matheus", user.Name)
	assert.Equal(t, "mh.rossetti2002@gmail.com", user.Email)
}

func TestUser_ValidatePassword(t *testing.T) {
	user, err := NewUser("Matheus", "mh.rossetti2002@gmail.com", "123")
	assert.Nil(t, err)
	assert.True(t, user.ValidatePassword("123"))
	assert.False(t, user.ValidatePassword("1234"))
	assert.NotEqual(t, "123", user.PasswordHash)
}

func TestUserWhenNameIsRequired(t *testing.T) {
	user, err := NewUser("", "mh.rossetti2002@gmail.com", "123")
	assert.NotNil(t, err)
	assert.Nil(t, user)
	assert.Equal(t, ErrNameIsRequired, err)
}

func TestUserWhenHashFails(t *testing.T) {
	user, err := NewUser("Matheus", "mh.rossetti2002@gmail.com", strings.Repeat("A", 73))
	assert.NotNil(t, err)
	assert.Nil(t, user)
}

func TestUserWhenEmailIsRequired(t *testing.T) {
	user, err := NewUser("Matheus", "", "123")
	assert.NotNil(t, err)
	assert.Nil(t, user)
	assert.Equal(t, ErrEmailIsRequired, err)
}

func TestUserValidate(t *testing.T) {
	user, err := NewUser("Matheus", "mh.rossetti2002@gmail.com", "123")
	assert.NotNil(t, user)
	assert.Nil(t, err)
	assert.Nil(t, user.Validate())
}
