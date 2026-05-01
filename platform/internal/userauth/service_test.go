package userauth_test

import (
	"context"
	"database/sql"
	"testing"

	"hostamat/internal/adapters/yandexoauth"
	"hostamat/internal/userauth"

	"github.com/platforma-dev/platforma/log"
)

type mockUserauthRepository struct {
	user             *userauth.User
	getByProviderErr error
	createdUser      *userauth.User
	updatedID        string
	updatedUsername  string
	updatedEmail     *string
}

func (m *mockUserauthRepository) Get(ctx context.Context, id string) (*userauth.User, error) {
	return nil, sql.ErrNoRows
}

func (m *mockUserauthRepository) GetByProviderUserID(ctx context.Context, provider userauth.Provider, providerUserID string) (*userauth.User, error) {
	if m.getByProviderErr != nil {
		return nil, m.getByProviderErr
	}

	return m.user, nil
}

func (m *mockUserauthRepository) Create(ctx context.Context, user *userauth.User) error {
	m.createdUser = user
	return nil
}

func (m *mockUserauthRepository) UpdateOAuthProfile(ctx context.Context, id, username string, email *string) error {
	m.updatedID = id
	m.updatedUsername = username
	m.updatedEmail = email
	return nil
}

func (m *mockUserauthRepository) Delete(ctx context.Context, id string) error {
	return nil
}

type mockUserauthStorage struct {
	createdUserID string
}

func (m *mockUserauthStorage) GetUserIdFromSessionId(context.Context, string) (string, error) {
	return "", nil
}

func (m *mockUserauthStorage) CreateSessionForUser(ctx context.Context, userID string) (string, error) {
	m.createdUserID = userID
	return "session-id", nil
}

func (m *mockUserauthStorage) DeleteSession(ctx context.Context, sessionID string) error {
	return nil
}

func (m *mockUserauthStorage) DeleteSessionsByUserId(ctx context.Context, userID string) error {
	return nil
}

type mockYandexClient struct {
	userInfo *yandexoauth.UserInfo
}

func (m *mockYandexClient) AuthorizeURL(state string) string {
	return ""
}

func (m *mockYandexClient) ExchangeCode(ctx context.Context, code string) (string, error) {
	return "access-token", nil
}

func (m *mockYandexClient) GetUserInfo(ctx context.Context, accessToken string) (*yandexoauth.UserInfo, error) {
	return m.userInfo, nil
}

func TestCreateSessionFromYandexCode_CreatesUserWithEmail(t *testing.T) {
	repo := &mockUserauthRepository{getByProviderErr: sql.ErrNoRows}
	storage := &mockUserauthStorage{}
	service := userauth.NewService(repo, storage, &mockYandexClient{
		userInfo: &yandexoauth.UserInfo{
			ID:           "provider-user-id",
			Login:        "alice",
			DefaultEmail: "alice@example.com",
		},
	}, "session", "/api/auth/yandex/callback", nil)

	sessionID, err := service.CreateSessionFromYandexCode(testLogContext(), "code")
	if err != nil {
		t.Fatalf("CreateSessionFromYandexCode() error = %v", err)
	}

	if sessionID != "session-id" {
		t.Fatalf("unexpected session id: %s", sessionID)
	}

	if repo.createdUser == nil {
		t.Fatalf("expected user to be created")
	}

	if repo.createdUser.Email == nil || *repo.createdUser.Email != "alice@example.com" {
		t.Fatalf("expected email to be saved, got %#v", repo.createdUser.Email)
	}

	if storage.createdUserID != repo.createdUser.ID {
		t.Fatalf("expected session for created user, got %s", storage.createdUserID)
	}
}

func TestCreateSessionFromYandexCode_UpdatesStoredEmail(t *testing.T) {
	oldEmail := "old@example.com"
	repo := &mockUserauthRepository{
		user: &userauth.User{
			ID:       "user-id",
			Username: "alice",
			Email:    &oldEmail,
			Provider: userauth.ProviderYandex,
		},
	}
	service := userauth.NewService(repo, &mockUserauthStorage{}, &mockYandexClient{
		userInfo: &yandexoauth.UserInfo{
			ID:           "provider-user-id",
			Login:        "alice",
			DefaultEmail: "alice@example.com",
		},
	}, "session", "/api/auth/yandex/callback", nil)

	if _, err := service.CreateSessionFromYandexCode(testLogContext(), "code"); err != nil {
		t.Fatalf("CreateSessionFromYandexCode() error = %v", err)
	}

	if repo.updatedID != "user-id" {
		t.Fatalf("expected user to be updated, got %s", repo.updatedID)
	}

	if repo.updatedUsername != "alice" {
		t.Fatalf("expected username to be preserved, got %s", repo.updatedUsername)
	}

	if repo.updatedEmail == nil || *repo.updatedEmail != "alice@example.com" {
		t.Fatalf("expected email to be updated, got %#v", repo.updatedEmail)
	}
}

func testLogContext() context.Context {
	return context.WithValue(context.Background(), log.WideEventKey, log.NewEvent("test"))
}
