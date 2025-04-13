package utils_test

import (
	"testing"

	"github.com/davidyannick86/grpc-api-mongodb/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestPassword(t *testing.T) {
	password := "my_very_secret_password"

	// Hash the password
	hashPassword, err := utils.HashPassword(password)
	assert.NoError(t, err)

	// Verify the password
	err = utils.VerifyPassword(password, hashPassword)
	assert.NoError(t, err)
}
