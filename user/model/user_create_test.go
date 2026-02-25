package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateUsername(t *testing.T) {
	allowed := []string{
		"john",
		"john_doe",
		"john-doe",
		"j0hn",
		"abcd",
		"user1234",
		"a1b2c3d4",
		"test-user-name",
		"test_user_name",
		"abcdefghijklmnopqrst", // 20 chars, at limit
	}

	disallowed := []struct {
		username string
		err      error
	}{
		{"", ErrUsernameMissing},
		{"ab", ErrUsernameMinLength},
		{"a", ErrUsernameMinLength},
		{"abcdefghijklmnopqrstu", ErrUsernameMaxLength}, // 21 chars
		{"-john", ErrUsernameInvalid},
		{"john-", ErrUsernameInvalid},
		{"_john", ErrUsernameInvalid},
		{"john_", ErrUsernameInvalid},
		{"john.doe", ErrUsernameInvalid},
		{"John", ErrUsernameInvalid},
		{"JOHN", ErrUsernameInvalid},
		{"john doe", ErrUsernameInvalid},
		{"john%doe", ErrUsernameInvalid},
		{"<script>", ErrUsernameInvalid},
		{"john&doe", ErrUsernameInvalid},
		{"john/doe", ErrUsernameInvalid},
		{"john@doe", ErrUsernameInvalid},
		{"john+doe", ErrUsernameInvalid},
		{"über", ErrUsernameInvalid},
		{"john=doe", ErrUsernameInvalid},
		{"john?doe", ErrUsernameInvalid},
		{"john#doe", ErrUsernameInvalid},
	}

	for _, username := range allowed {
		t.Run("allow/"+username, func(t *testing.T) {
			req := &UserCreateRequest{Username: username}
			assert.NoError(t, req.ValidateUsername())
		})
	}

	for _, tt := range disallowed {
		t.Run("reject/"+tt.username, func(t *testing.T) {
			req := &UserCreateRequest{Username: tt.username}
			err := req.ValidateUsername()
			assert.ErrorIs(t, err, tt.err)
		})
	}
}
