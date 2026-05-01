package apikey_test

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"hostamat/internal/apikey"

	platformauth "github.com/platforma-dev/platforma/auth"
)

func TestServiceCreateAndAuthenticate(t *testing.T) {
	repo := newRepositoryMock("user-1")
	service := apikey.NewService(repo, "pepper")

	response, err := service.Create(testContext("user-1"), apikey.CreateRequest{
		Name:   "CI",
		Scopes: []apikey.Scope{apikey.ScopeSitesRead, apikey.ScopeSitesUpload},
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if response.Secret == "" {
		t.Fatal("Create() secret is empty")
	}
	if !strings.HasPrefix(response.Secret, "hmat_"+response.Key.ID+"_") {
		t.Fatalf("Create() secret prefix = %q, want hmat_<id>_", response.Secret)
	}

	stored := repo.keys[response.Key.ID]
	if stored.SecretDigest == "" {
		t.Fatal("stored SecretDigest is empty")
	}
	if strings.Contains(response.Secret, stored.SecretDigest) {
		t.Fatal("secret contains digest unexpectedly")
	}
	if strings.Contains(stored.SecretDigest, strings.TrimPrefix(response.Secret, "hmat_"+response.Key.ID+"_")) {
		t.Fatal("digest contains plaintext secret")
	}

	authenticated, err := service.Authenticate(context.Background(), response.Secret)
	if err != nil {
		t.Fatalf("Authenticate() error = %v", err)
	}
	if authenticated.Owner.ID != "user-1" {
		t.Fatalf("Authenticate() owner = %q, want user-1", authenticated.Owner.ID)
	}
	if repo.lastUsedUpdates != 1 {
		t.Fatalf("last used updates = %d, want 1", repo.lastUsedUpdates)
	}
}

func TestServiceCreateValidation(t *testing.T) {
	repo := newRepositoryMock("user-1")
	service := apikey.NewService(repo, "pepper")
	past := time.Now().Add(-time.Minute)

	tests := []struct {
		name    string
		request apikey.CreateRequest
		wantErr error
	}{
		{
			name:    "missing name",
			request: apikey.CreateRequest{Name: "", Scopes: []apikey.Scope{apikey.ScopeSitesRead}},
			wantErr: apikey.ErrMissingName,
		},
		{
			name:    "missing scopes",
			request: apikey.CreateRequest{Name: "CI"},
			wantErr: apikey.ErrMissingScopes,
		},
		{
			name:    "invalid scope",
			request: apikey.CreateRequest{Name: "CI", Scopes: []apikey.Scope{"unknown"}},
			wantErr: apikey.ErrInvalidScope,
		},
		{
			name:    "past expiration",
			request: apikey.CreateRequest{Name: "CI", Scopes: []apikey.Scope{apikey.ScopeSitesRead}, ExpiresAt: &past},
			wantErr: apikey.ErrInvalidExpiration,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := service.Create(testContext("user-1"), tt.request); err == nil || !strings.Contains(err.Error(), tt.wantErr.Error()) {
				t.Fatalf("Create() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceAuthenticateRejectsExpiredRevokedAndInactiveOwner(t *testing.T) {
	repo := newRepositoryMock("user-1")
	service := apikey.NewService(repo, "pepper")

	expiredAt := time.Now().Add(time.Minute)
	created, err := service.Create(testContext("user-1"), apikey.CreateRequest{
		Name:      "CI",
		Scopes:    []apikey.Scope{apikey.ScopeSitesRead},
		ExpiresAt: &expiredAt,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	repo.keys[created.Key.ID].ExpiresAt = sql.NullTime{Time: time.Now().Add(-time.Minute), Valid: true}
	if _, err := service.Authenticate(context.Background(), created.Secret); err != apikey.ErrExpiredKey {
		t.Fatalf("Authenticate() expired error = %v, want %v", err, apikey.ErrExpiredKey)
	}

	repo.keys[created.Key.ID].ExpiresAt = sql.NullTime{}
	repo.keys[created.Key.ID].RevokedAt = sql.NullTime{Time: time.Now(), Valid: true}
	if _, err := service.Authenticate(context.Background(), created.Secret); err != apikey.ErrRevokedKey {
		t.Fatalf("Authenticate() revoked error = %v, want %v", err, apikey.ErrRevokedKey)
	}

	repo.keys[created.Key.ID].RevokedAt = sql.NullTime{}
	repo.owner.Status = platformauth.StatusInactive
	if _, err := service.Authenticate(context.Background(), created.Secret); err != apikey.ErrOwnerInactive {
		t.Fatalf("Authenticate() inactive owner error = %v, want %v", err, apikey.ErrOwnerInactive)
	}
}

type repositoryMock struct {
	keys            map[string]*apikey.APIKey
	owner           *apikey.Owner
	lastUsedUpdates int
}

func newRepositoryMock(userID string) *repositoryMock {
	return &repositoryMock{
		keys: make(map[string]*apikey.APIKey),
		owner: &apikey.Owner{
			ID:       userID,
			Username: "tester",
			Created:  time.Now(),
			Updated:  time.Now(),
			Status:   platformauth.StatusActive,
		},
	}
}

func (m *repositoryMock) Create(_ context.Context, key *apikey.APIKey) error {
	copyKey := *key
	m.keys[key.ID] = &copyKey
	return nil
}

func (m *repositoryMock) Get(_ context.Context, id string) (*apikey.APIKey, error) {
	key, ok := m.keys[id]
	if !ok {
		return nil, sql.ErrNoRows
	}

	return key, nil
}

func (m *repositoryMock) GetForUser(_ context.Context, userID, id string) (*apikey.APIKey, error) {
	key, ok := m.keys[id]
	if !ok || key.UserID != userID {
		return nil, sql.ErrNoRows
	}

	return key, nil
}

func (m *repositoryMock) ListVisibleForUser(_ context.Context, userID string) ([]*apikey.APIKey, error) {
	keys := make([]*apikey.APIKey, 0)
	for _, key := range m.keys {
		if key.UserID == userID && !key.RevokedAt.Valid {
			keys = append(keys, key)
		}
	}

	return keys, nil
}

func (m *repositoryMock) CountActiveForUser(_ context.Context, userID string, now time.Time) (int, error) {
	count := 0
	for _, key := range m.keys {
		if key.UserID == userID && !key.RevokedAt.Valid && (!key.ExpiresAt.Valid || key.ExpiresAt.Time.After(now)) {
			count++
		}
	}

	return count, nil
}

func (m *repositoryMock) UpdateName(_ context.Context, key *apikey.APIKey) error {
	m.keys[key.ID] = key
	return nil
}

func (m *repositoryMock) Revoke(_ context.Context, id string, now time.Time) error {
	key, ok := m.keys[id]
	if !ok {
		return sql.ErrNoRows
	}

	key.RevokedAt = sql.NullTime{Time: now, Valid: true}
	return nil
}

func (m *repositoryMock) UpdateLastUsed(_ context.Context, id string, now, _ time.Time) error {
	key, ok := m.keys[id]
	if !ok {
		return sql.ErrNoRows
	}

	key.LastUsedAt = sql.NullTime{Time: now, Valid: true}
	m.lastUsedUpdates++
	return nil
}

func (m *repositoryMock) GetOwner(_ context.Context, userID string) (*apikey.Owner, error) {
	if m.owner == nil || m.owner.ID != userID {
		return nil, sql.ErrNoRows
	}

	return m.owner, nil
}

func testContext(userID string) context.Context {
	return context.WithValue(context.Background(), platformauth.UserContextKey, &platformauth.User{
		ID:       userID,
		Username: "tester",
		Status:   platformauth.StatusActive,
	})
}
