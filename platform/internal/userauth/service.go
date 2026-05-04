package userauth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"hostamat/internal/adapters/yandexoauth"

	"github.com/google/uuid"
	platformauth "github.com/platforma-dev/platforma/auth"
	"github.com/platforma-dev/platforma/log"
)

const oauthStateCookieName = "hostamat_oauth_state"

type repository interface {
	Get(ctx context.Context, id string) (*User, error)
	GetByProviderUserID(ctx context.Context, provider Provider, providerUserID string) (*User, error)
	Create(ctx context.Context, user *User) error
	UpdateOAuthProfile(ctx context.Context, id, username string, email *string) error
	Delete(ctx context.Context, id string) error
}

type authStorage interface {
	GetUserIdFromSessionId(context.Context, string) (string, error)
	CreateSessionForUser(context.Context, string) (string, error)
	DeleteSession(ctx context.Context, sessionID string) error
	DeleteSessionsByUserId(ctx context.Context, userID string) error
}

type yandexClient interface {
	AuthorizeURL(state string) string
	ExchangeCode(ctx context.Context, code string) (string, error)
	GetUserInfo(ctx context.Context, accessToken string) (*yandexoauth.UserInfo, error)
}

type cleanupEnqueuer interface {
	Enqueue(ctx context.Context, job platformauth.UserCleanupJob) error
}

type Service struct {
	repo              repository
	authStorage       authStorage
	yandexClient      yandexClient
	sessionCookieName string
	oauthStatePath    string
	cleanupEnqueuer   cleanupEnqueuer
}

func NewService(repo repository, authStorage authStorage, yandexClient yandexClient, sessionCookieName, redirectURI string, cleanupEnqueuer cleanupEnqueuer) *Service {
	return &Service{
		repo:              repo,
		authStorage:       authStorage,
		yandexClient:      yandexClient,
		sessionCookieName: sessionCookieName,
		oauthStatePath:    resolveOAuthStatePath(redirectURI),
		cleanupEnqueuer:   cleanupEnqueuer,
	}
}

func (s *Service) CookieName() string {
	return s.sessionCookieName
}

func (s *Service) OAuthStateCookieName() string {
	return oauthStateCookieName
}

func (s *Service) OAuthStateCookiePath() string {
	return s.oauthStatePath
}

func (s *Service) StartYandexAuth() (string, string, error) {
	state, err := generateOAuthState()
	if err != nil {
		return "", "", err
	}

	return s.yandexClient.AuthorizeURL(state), state, nil
}

func (s *Service) CreateSessionFromYandexCode(ctx context.Context, code string) (string, error) {
	ev := log.EventFromContext(ctx)

	ev.AddStep(slog.LevelInfo, "userauth.CreateSessionFromYandexCode")
	ev.AddAttrs(map[string]any{
		"auth.provider": "yandex",
	})

	if strings.TrimSpace(code) == "" {

		ev.AddStep(slog.LevelWarn, "userauth.CreateSessionFromYandexCode.missing_code")

		return "", fmt.Errorf("oauth code is required")
	}

	ev.AddStep(slog.LevelInfo, "userauth.CreateSessionFromYandexCode.exchange_code")

	accessToken, err := s.yandexClient.ExchangeCode(ctx, code)
	if err != nil {

		ev.AddError(err)

		return "", err
	}

	ev.AddStep(slog.LevelInfo, "userauth.CreateSessionFromYandexCode.fetch_userinfo")

	userInfo, err := s.yandexClient.GetUserInfo(ctx, accessToken)
	if err != nil {

		ev.AddError(err)

		return "", err
	}

	ev.AddStep(slog.LevelInfo, "userauth.CreateSessionFromYandexCode.load_user")

	user, err := s.repo.GetByProviderUserID(ctx, ProviderYandex, userInfo.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {

		ev.AddError(err)

		return "", fmt.Errorf("failed to fetch existing oauth user: %w", err)
	}

	username := resolveUsername(userInfo)
	email := resolveEmail(userInfo)

	if errors.Is(err, sql.ErrNoRows) {
		providerUserID := userInfo.ID
		user = &User{
			ID:             uuid.NewString(),
			Username:       username,
			Email:          email,
			Password:       "",
			Salt:           "",
			Created:        time.Now(),
			Updated:        time.Now(),
			Status:         platformauth.StatusActive,
			Provider:       ProviderYandex,
			ProviderUserID: &providerUserID,
		}

		ev.AddStep(slog.LevelInfo, "userauth.CreateSessionFromYandexCode.create_user")

		if err := s.repo.Create(ctx, user); err != nil {

			ev.AddError(err)

			return "", fmt.Errorf("failed to create oauth user: %w", err)
		}
	} else if user.Username != username || !equalStringPointers(user.Email, email) {

		ev.AddStep(slog.LevelInfo, "userauth.CreateSessionFromYandexCode.update_oauth_profile")

		if err := s.repo.UpdateOAuthProfile(ctx, user.ID, username, email); err != nil {

			ev.AddError(err)

			return "", fmt.Errorf("failed to update oauth profile: %w", err)
		}
		user.Username = username
		user.Email = email
	}

	ev.AddAttrs(map[string]any{
		"user.id":       user.ID,
		"user.provider": string(user.Provider),
	})
	ev.AddStep(slog.LevelInfo, "userauth.CreateSessionFromYandexCode.create_session")

	sessionID, err := s.authStorage.CreateSessionForUser(ctx, user.ID)
	if err != nil {

		ev.AddError(err)

		return "", fmt.Errorf("failed to create session: %w", err)
	}

	return sessionID, nil
}

func (s *Service) GetFromSession(ctx context.Context, sessionID string) (*platformauth.User, error) {
	ev := log.EventFromContext(ctx)

	ev.AddStep(slog.LevelInfo, "userauth.GetFromSession")
	ev.AddAttrs(map[string]any{
		"auth.provider": "yandex",
	})

	userID, err := s.authStorage.GetUserIdFromSessionId(ctx, sessionID)
	if err != nil {

		ev.AddError(err)

		return nil, fmt.Errorf("failed to get user id from session: %w", err)
	}

	if userID == "" {

		ev.AddStep(slog.LevelInfo, "userauth.GetFromSession.user_not_found")

		return nil, platformauth.ErrUserNotFound
	}

	ev.AddAttrs(map[string]any{
		"user.id": userID,
	})

	user, err := s.repo.Get(ctx, userID)
	if errors.Is(err, sql.ErrNoRows) {

		ev.AddStep(slog.LevelInfo, "userauth.GetFromSession.user_not_found")

		return nil, platformauth.ErrUserNotFound
	}
	if err != nil {

		ev.AddError(err)

		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user.Provider != ProviderYandex {

		ev.AddStep(slog.LevelWarn, "userauth.GetFromSession.invalid_provider")

		return nil, platformauth.ErrUserNotFound
	}

	ev.AddAttrs(map[string]any{
		"user.provider": string(user.Provider),
	})

	return &platformauth.User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Salt:     user.Salt,
		Created:  user.Created,
		Updated:  user.Updated,
		Status:   user.Status,
	}, nil
}

func (s *Service) GetProfileFromSession(ctx context.Context, sessionID string) (*User, error) {
	userID, err := s.authStorage.GetUserIdFromSessionId(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user id from session: %w", err)
	}
	if userID == "" {
		return nil, platformauth.ErrUserNotFound
	}

	user, err := s.repo.Get(ctx, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, platformauth.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user.Provider != ProviderYandex {
		return nil, platformauth.ErrUserNotFound
	}

	return user, nil
}

func (s *Service) DeleteSession(ctx context.Context, sessionID string) error {
	ev := log.EventFromContext(ctx)

	ev.AddStep(slog.LevelInfo, "userauth.DeleteSession")

	if err := s.authStorage.DeleteSession(ctx, sessionID); err != nil {

		ev.AddError(err)

		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

func (s *Service) DeleteUser(ctx context.Context) error {
	ev := log.EventFromContext(ctx)

	ev.AddStep(slog.LevelInfo, "userauth.DeleteUser")
	ev.AddAttrs(map[string]any{
		"auth.provider": "yandex",
	})

	user := platformauth.UserFromContext(ctx)
	if user == nil {

		ev.AddStep(slog.LevelWarn, "userauth.DeleteUser.user_not_found")

		return platformauth.ErrUserNotFound
	}

	ev.AddAttrs(map[string]any{
		"user.id": user.ID,
	})

	if err := s.authStorage.DeleteSessionsByUserId(ctx, user.ID); err != nil {

		ev.AddError(err)

		return fmt.Errorf("failed to delete user sessions: %w", err)
	}

	if err := s.repo.Delete(ctx, user.ID); err != nil {

		ev.AddError(err)

		return fmt.Errorf("failed to delete user: %w", err)
	}

	if s.cleanupEnqueuer != nil {
		job := platformauth.UserCleanupJob{
			UserID:    user.ID,
			DeletedAt: time.Now(),
		}
		if err := s.cleanupEnqueuer.Enqueue(ctx, job); err != nil {

			ev.AddError(err)

			log.ErrorContext(ctx, "failed to enqueue cleanup job", "error", err, "user_id", user.ID)
		}
	}

	return nil
}

func generateOAuthState() (string, error) {
	stateBytes := make([]byte, 32)
	if _, err := rand.Read(stateBytes); err != nil {
		return "", fmt.Errorf("failed to generate oauth state: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(stateBytes), nil
}

func resolveOAuthStatePath(redirectURI string) string {
	parsedURI, err := url.Parse(redirectURI)
	if err != nil || parsedURI.Path == "" {
		return "/api/auth/yandex/callback"
	}

	return parsedURI.Path
}

func resolveUsername(userInfo *yandexoauth.UserInfo) string {
	if userInfo.Login != "" {
		return userInfo.Login
	}

	if userInfo.DefaultEmail != "" {
		return userInfo.DefaultEmail
	}

	return userInfo.ID
}

func resolveEmail(userInfo *yandexoauth.UserInfo) *string {
	email := strings.TrimSpace(userInfo.DefaultEmail)
	if email == "" {
		return nil
	}

	return &email
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}

func equalStringPointers(left, right *string) bool {
	if left == nil || right == nil {
		return left == right
	}

	return *left == *right
}
