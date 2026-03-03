package passkey

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/titpetric/platform/pkg/require"

	"github.com/titpetric/platform-app/user/model"
)

type mockPasskeyStorage struct {
	passkeys map[string]*model.UserPasskey
	byCredID map[string]*model.UserPasskey
}

func newMockPasskeyStorage() *mockPasskeyStorage {
	return &mockPasskeyStorage{
		passkeys: make(map[string]*model.UserPasskey),
		byCredID: make(map[string]*model.UserPasskey),
	}
}

func (m *mockPasskeyStorage) Create(_ context.Context, passkey *model.UserPasskey) (*model.UserPasskey, error) {
	m.passkeys[passkey.ID] = passkey
	m.byCredID[string(passkey.CredentialID)] = passkey
	return passkey, nil
}

func (m *mockPasskeyStorage) Delete(_ context.Context, id string) error {
	delete(m.passkeys, id)
	return nil
}

func (m *mockPasskeyStorage) ListByUser(_ context.Context, userID string) ([]model.UserPasskey, error) {
	var result []model.UserPasskey
	for _, p := range m.passkeys {
		if p.UserID == userID {
			result = append(result, *p)
		}
	}
	return result, nil
}

func (m *mockPasskeyStorage) GetByCredentialID(_ context.Context, credentialID []byte) (*model.UserPasskey, error) {
	if p, ok := m.byCredID[string(credentialID)]; ok {
		return p, nil
	}
	return nil, errors.New("passkey not found")
}

func (m *mockPasskeyStorage) UpdateSignCount(_ context.Context, id string, signCount int64) error {
	if p, ok := m.passkeys[id]; ok {
		p.SignCount = signCount
	}
	return nil
}

type mockUserStorage struct {
	users map[string]*model.User
}

func newMockUserStorage() *mockUserStorage {
	return &mockUserStorage{
		users: make(map[string]*model.User),
	}
}

func (m *mockUserStorage) Create(_ context.Context, req *model.UserCreateRequest) (*model.User, error) {
	u := &model.User{
		ID:       "user-" + req.Username,
		Username: req.Username,
		FullName: req.FullName,
	}
	m.users[u.ID] = u
	return u, nil
}

func (m *mockUserStorage) Update(_ context.Context, u *model.User) (*model.User, error) {
	m.users[u.ID] = u
	return u, nil
}

func (m *mockUserStorage) Get(_ context.Context, id string) (*model.User, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return nil, errors.New("user not found")
}

func (m *mockUserStorage) GetByUsername(_ context.Context, username string) (*model.User, error) {
	for _, u := range m.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *mockUserStorage) GetByStub(_ context.Context, _ string) (*model.User, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserStorage) GetGroups(_ context.Context, _ string) ([]model.UserGroup, error) {
	return nil, nil
}

func (m *mockUserStorage) List(_ context.Context) ([]model.User, error) {
	var result []model.User
	for _, u := range m.users {
		result = append(result, *u)
	}
	return result, nil
}

func (m *mockUserStorage) Authenticate(_ context.Context, _ model.UserAuth) (*model.User, error) {
	return nil, errors.New("not implemented")
}

func newTestWebAuthn(t *testing.T) *webauthn.WebAuthn {
	wa, err := webauthn.New(&webauthn.Config{
		RPID:          "localhost",
		RPDisplayName: "Test App",
		RPOrigins:     []string{"http://localhost:3000"},
	})
	require.Nil(t, err)
	return wa
}

func TestNew(t *testing.T) {
	wa := newTestWebAuthn(t)
	ps := newMockPasskeyStorage()
	us := newMockUserStorage()

	svc := New(wa, ps, us)
	require.NotNil(t, svc)
}

func TestBeginRegistration(t *testing.T) {
	wa := newTestWebAuthn(t)
	ps := newMockPasskeyStorage()
	us := newMockUserStorage()
	svc := New(wa, ps, us)

	req := &model.UserCreateRequest{
		Username: "testuser",
		FullName: "Test User",
	}

	token, options, err := svc.BeginRegistration(req)
	require.Nil(t, err)
	require.NotEmpty(t, token)
	require.NotNil(t, options)
	require.NotNil(t, options.Response.Challenge)
}

func TestFinishRegistrationInvalidToken(t *testing.T) {
	wa := newTestWebAuthn(t)
	ps := newMockPasskeyStorage()
	us := newMockUserStorage()
	svc := New(wa, ps, us)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	_, err := svc.FinishRegistration("invalid-token", req)
	require.NotNil(t, err)

	var passkeyErr *Error
	require.True(t, errors.As(err, &passkeyErr))
	require.Equal(t, http.StatusBadRequest, passkeyErr.Status)
}

func TestBeginLogin(t *testing.T) {
	wa := newTestWebAuthn(t)
	ps := newMockPasskeyStorage()
	us := newMockUserStorage()
	svc := New(wa, ps, us)

	token, options, err := svc.BeginLogin()
	require.Nil(t, err)
	require.NotEmpty(t, token)
	require.NotNil(t, options)
	require.NotNil(t, options.Response.Challenge)
}

func TestFinishLoginInvalidToken(t *testing.T) {
	wa := newTestWebAuthn(t)
	ps := newMockPasskeyStorage()
	us := newMockUserStorage()
	svc := New(wa, ps, us)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	_, err := svc.FinishLogin("invalid-token", req)
	require.NotNil(t, err)

	var passkeyErr *Error
	require.True(t, errors.As(err, &passkeyErr))
	require.Equal(t, http.StatusBadRequest, passkeyErr.Status)
}

func TestConsumeSessionExpired(t *testing.T) {
	wa := newTestWebAuthn(t)
	ps := newMockPasskeyStorage()
	us := newMockUserStorage()
	svc := New(wa, ps, us)

	// Try to consume a session that doesn't exist
	_, err := svc.consumeSession("nonexistent")
	require.NotNil(t, err)

	var passkeyErr *Error
	require.True(t, errors.As(err, &passkeyErr))
	require.Equal(t, http.StatusBadRequest, passkeyErr.Status)
}

func TestError(t *testing.T) {
	err := &Error{
		Status: http.StatusBadRequest,
		Err:    errors.New("test error"),
	}
	require.Equal(t, "test error", err.Error())
}
