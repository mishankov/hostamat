package userauth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
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
	UpdateUsername(ctx context.Context, id, username string) error
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
	if strings.TrimSpace(code) == "" {
		return "", fmt.Errorf("oauth code is required")
	}

	accessToken, err := s.yandexClient.ExchangeCode(ctx, code)
	if err != nil {
		return "", err
	}

	userInfo, err := s.yandexClient.GetUserInfo(ctx, accessToken)
	if err != nil {
		return "", err
	}

	user, err := s.repo.GetByProviderUserID(ctx, ProviderYandex, userInfo.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("failed to fetch existing oauth user: %w", err)
	}

	username := resolveUsername(userInfo)

	if errors.Is(err, sql.ErrNoRows) {
		providerUserID := userInfo.ID
		user = &User{
			ID:             uuid.NewString(),
			Username:       username,
			Password:       "",
			Salt:           "",
			Created:        time.Now(),
			Updated:        time.Now(),
			Status:         platformauth.StatusActive,
			Provider:       ProviderYandex,
			ProviderUserID: &providerUserID,
		}

		if err := s.repo.Create(ctx, user); err != nil {
			return "", fmt.Errorf("failed to create oauth user: %w", err)
		}
	} else if user.Username != username {
		if err := s.repo.UpdateUsername(ctx, user.ID, username); err != nil {
			return "", fmt.Errorf("failed to update oauth username: %w", err)
		}
		user.Username = username
	}

	sessionID, err := s.authStorage.CreateSessionForUser(ctx, user.ID)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	return sessionID, nil
}

func (s *Service) GetFromSession(ctx context.Context, sessionID string) (*platformauth.User, error) {
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

func (s *Service) DeleteSession(ctx context.Context, sessionID string) error {
	if err := s.authStorage.DeleteSession(ctx, sessionID); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

func (s *Service) DeleteUser(ctx context.Context) error {
	user := platformauth.UserFromContext(ctx)
	if user == nil {
		return platformauth.ErrUserNotFound
	}

	if err := s.authStorage.DeleteSessionsByUserId(ctx, user.ID); err != nil {
		return fmt.Errorf("failed to delete user sessions: %w", err)
	}

	if err := s.repo.Delete(ctx, user.ID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if s.cleanupEnqueuer != nil {
		job := platformauth.UserCleanupJob{
			UserID:    user.ID,
			DeletedAt: time.Now(),
		}
		if err := s.cleanupEnqueuer.Enqueue(ctx, job); err != nil {
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
