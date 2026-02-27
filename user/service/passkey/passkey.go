package passkey

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/titpetric/platform/pkg/ulid"

	"github.com/titpetric/platform-app/user/model"
)

// Service handles WebAuthn passkey ceremony state.
type Service struct {
	webAuthn       *webauthn.WebAuthn
	passkeyStorage model.PasskeyStorage
	userStorage    model.UserStorage

	mu       sync.Mutex
	sessions map[string]*ceremonySession
}

type ceremonySession struct {
	SessionData *webauthn.SessionData
	UserRequest *model.UserCreateRequest
	ExpiresAt   time.Time
}

// New creates a new passkey Service.
func New(wa *webauthn.WebAuthn, ps model.PasskeyStorage, us model.UserStorage) *Service {
	return &Service{
		webAuthn:       wa,
		passkeyStorage: ps,
		userStorage:    us,
		sessions:       make(map[string]*ceremonySession),
	}
}

// RegistrationResult is returned from FinishRegistration.
type RegistrationResult struct {
	UserID string
}

// LoginResult is returned from FinishLogin.
type LoginResult struct {
	UserID string
}

// BeginRegistration starts the WebAuthn registration ceremony for a new user.
func (s *Service) BeginRegistration(req *model.UserCreateRequest) (token string, options *protocol.CredentialCreation, err error) {
	tempID := ulid.String()
	waUser := &model.WebAuthnUser{
		User: &model.User{
			ID:       tempID,
			Username: req.Username,
			FullName: req.FullName,
		},
	}

	creation, sessionData, err := s.webAuthn.BeginRegistration(waUser,
		webauthn.WithResidentKeyRequirement(protocol.ResidentKeyRequirementPreferred),
	)
	if err != nil {
		return "", nil, fmt.Errorf("begin registration: %w", err)
	}

	token = ulid.String()
	s.mu.Lock()
	s.sessions[token] = &ceremonySession{
		SessionData: sessionData,
		UserRequest: req,
		ExpiresAt:   time.Now().Add(5 * time.Minute),
	}
	s.mu.Unlock()

	return token, creation, nil
}

// FinishRegistration completes the WebAuthn registration, creates the user and stores the passkey.
func (s *Service) FinishRegistration(token string, r *http.Request) (*RegistrationResult, error) {
	cs, err := s.consumeSession(token)
	if err != nil {
		return nil, err
	}

	waUser := &model.WebAuthnUser{
		User: &model.User{
			ID:       string(cs.SessionData.UserID),
			Username: cs.UserRequest.Username,
			FullName: cs.UserRequest.FullName,
		},
	}

	credential, err := s.webAuthn.FinishRegistration(waUser, *cs.SessionData, r)
	if err != nil {
		return nil, &Error{Status: http.StatusBadRequest, Err: fmt.Errorf("finish registration: %w", err)}
	}

	ctx := r.Context()

	// Create the user (password not required for passkey registration).
	createReq := cs.UserRequest
	createReq.Password = ulid.String()
	user, err := s.userStorage.Create(ctx, createReq)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	// Store the passkey credential.
	passkey := &model.UserPasskey{
		UserID:          user.ID,
		CredentialID:    credential.ID,
		PublicKey:       credential.PublicKey,
		AttestationType: credential.AttestationType,
		Transport:       model.TransportJSON(credential.Transport),
		SignCount:       int64(credential.Authenticator.SignCount),
	}
	if _, err := s.passkeyStorage.Create(ctx, passkey); err != nil {
		return nil, fmt.Errorf("store passkey: %w", err)
	}

	return &RegistrationResult{UserID: user.ID}, nil
}

// BeginLogin starts a discoverable WebAuthn login ceremony.
func (s *Service) BeginLogin() (token string, options *protocol.CredentialAssertion, err error) {
	assertion, sessionData, err := s.webAuthn.BeginDiscoverableLogin()
	if err != nil {
		return "", nil, fmt.Errorf("begin login: %w", err)
	}

	token = ulid.String()
	s.mu.Lock()
	s.sessions[token] = &ceremonySession{
		SessionData: sessionData,
		ExpiresAt:   time.Now().Add(5 * time.Minute),
	}
	s.mu.Unlock()

	return token, assertion, nil
}

// FinishLogin completes the discoverable WebAuthn login.
func (s *Service) FinishLogin(token string, r *http.Request) (*LoginResult, error) {
	cs, err := s.consumeSession(token)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()

	handler := func(rawID, userHandle []byte) (webauthn.User, error) {
		userID := string(userHandle)
		passkeys, err := s.passkeyStorage.ListByUser(ctx, userID)
		if err != nil {
			return nil, err
		}
		user, err := s.userStorage.Get(ctx, userID)
		if err != nil {
			return nil, err
		}
		return &model.WebAuthnUser{User: user, Passkeys: passkeys}, nil
	}

	_, credential, err := s.webAuthn.FinishPasskeyLogin(handler, *cs.SessionData, r)
	if err != nil {
		return nil, &Error{Status: http.StatusUnauthorized, Err: fmt.Errorf("finish login: %w", err)}
	}

	// Update sign count.
	stored, err := s.passkeyStorage.GetByCredentialID(ctx, credential.ID)
	if err != nil {
		return nil, fmt.Errorf("lookup passkey: %w", err)
	}
	_ = s.passkeyStorage.UpdateSignCount(ctx, stored.ID, int64(credential.Authenticator.SignCount))

	return &LoginResult{UserID: stored.UserID}, nil
}

func (s *Service) consumeSession(token string) (*ceremonySession, error) {
	s.mu.Lock()
	cs, ok := s.sessions[token]
	if ok {
		delete(s.sessions, token)
	}
	s.mu.Unlock()

	if !ok || time.Now().After(cs.ExpiresAt) {
		return nil, &Error{Status: http.StatusBadRequest, Err: fmt.Errorf("invalid or expired passkey session")}
	}
	return cs, nil
}
