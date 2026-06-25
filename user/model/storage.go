package model

import (
	"context"
)

// SessionStorage defines the storage operations for user sessions.
type SessionStorage interface {
	Create(ctx context.Context, userID string) (*UserSession, error)
	Get(ctx context.Context, sessionID string) (*UserSession, error)
	Delete(ctx context.Context, sessionID string) error
}

// UserStorage defines the storage operations for users.
type UserStorage interface {
	Create(context.Context, *UserCreateRequest) (*User, error)
	Update(context.Context, *User) (*User, error)

	Get(context.Context, string) (*User, error)
	GetByUsername(context.Context, string) (*User, error)
	GetByStub(context.Context, string) (*User, error)

	GetGroups(context.Context, string) ([]UserGroup, error)

	List(context.Context) ([]User, error)

	Authenticate(context.Context, UserAuth) (*User, error)
	Activate(ctx context.Context, token string) (*User, error)

	ResetActivation(ctx context.Context, email string) error
}

// PasskeyStorage defines the storage operations for WebAuthn passkeys.
type PasskeyStorage interface {
	Create(ctx context.Context, passkey *UserPasskey) (*UserPasskey, error)
	Delete(ctx context.Context, id string) error
	ListByUser(ctx context.Context, userID string) ([]UserPasskey, error)
	GetByCredentialID(ctx context.Context, credentialID []byte) (*UserPasskey, error)
	UpdateSignCount(ctx context.Context, id string, signCount int64) error
}
