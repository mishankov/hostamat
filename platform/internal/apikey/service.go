package apikey

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/platforma-dev/platforma/auth"
	platformauth "github.com/platforma-dev/platforma/auth"
)

type repository interface {
	Create(ctx context.Context, key *APIKey) error
	Get(ctx context.Context, id string) (*APIKey, error)
	GetForUser(ctx context.Context, userID, id string) (*APIKey, error)
	ListVisibleForUser(ctx context.Context, userID string) ([]*APIKey, error)
	CountActiveForUser(ctx context.Context, userID string, now time.Time) (int, error)
	UpdateName(ctx context.Context, key *APIKey) error
	Revoke(ctx context.Context, id string, now time.Time) error
	UpdateLastUsed(ctx context.Context, id string, now, threshold time.Time) error
	GetOwner(ctx context.Context, userID string) (*Owner, error)
}

type Service struct {
	repo   repository
	pepper string
	now    func() time.Time
}

func NewService(repo repository, pepper string) *Service {
	return &Service{
		repo:   repo,
		pepper: pepper,
		now:    time.Now,
	}
}

func (s *Service) Create(ctx context.Context, request CreateRequest) (*CreateResponse, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, ErrOwnerInactive
	}

	if err := validateName(request.Name); err != nil {
		return nil, err
	}

	if err := validateScopes(request.Scopes); err != nil {
		return nil, err
	}

	now := s.now()
	if request.ExpiresAt != nil && !request.ExpiresAt.After(now) {
		return nil, ErrInvalidExpiration
	}

	count, err := s.repo.CountActiveForUser(ctx, user.ID, now)
	if err != nil {
		return nil, err
	}

	if count >= maxActiveKeys {
		return nil, ErrActiveKeyLimit
	}

	id := uuid.NewString()
	secret, err := generateSecret()
	if err != nil {
		return nil, err
	}

	fullSecret := formatSecret(id, secret)
	key := &APIKey{
		ID:           id,
		UserID:       user.ID,
		Name:         strings.TrimSpace(request.Name),
		KeyPrefix:    keyPrefix(id),
		SecretDigest: s.digest(id, secret),
		Scopes:       stringsFromScopes(request.Scopes),
		ExpiresAt:    timePtrToNull(request.ExpiresAt),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.Create(ctx, key); err != nil {
		return nil, err
	}

	return &CreateResponse{
		Key:    key.Metadata(now),
		Secret: fullSecret,
	}, nil
}

func (s *Service) List(ctx context.Context) ([]PublicKeyMetadata, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, ErrOwnerInactive
	}

	keys, err := s.repo.ListVisibleForUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	now := s.now()
	response := make([]PublicKeyMetadata, 0, len(keys))
	for _, key := range keys {
		response = append(response, key.Metadata(now))
	}

	return response, nil
}

func (s *Service) UpdateName(ctx context.Context, id, name string) (*PublicKeyMetadata, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, ErrOwnerInactive
	}

	if err := validateName(name); err != nil {
		return nil, err
	}

	key, err := s.repo.GetForUser(ctx, user.ID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUnknownKey
		}

		return nil, err
	}

	now := s.now()
	key.Name = strings.TrimSpace(name)
	key.UpdatedAt = now

	if err := s.repo.UpdateName(ctx, key); err != nil {
		return nil, err
	}

	metadata := key.Metadata(now)
	return &metadata, nil
}

func (s *Service) Revoke(ctx context.Context, id string) error {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return ErrOwnerInactive
	}

	if _, err := s.repo.GetForUser(ctx, user.ID, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUnknownKey
		}

		return err
	}

	return s.repo.Revoke(ctx, id, s.now())
}

func (s *Service) Authenticate(ctx context.Context, bearer string) (*AuthenticatedKey, error) {
	id, secret, err := parseSecret(bearer)
	if err != nil {
		return nil, err
	}

	key, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUnknownKey
		}

		return nil, err
	}

	now := s.now()
	if key.RevokedAt.Valid {
		return nil, ErrRevokedKey
	}

	if key.ExpiresAt.Valid && !key.ExpiresAt.Time.After(now) {
		return nil, ErrExpiredKey
	}

	expectedDigest := s.digest(id, secret)
	if subtle.ConstantTimeCompare([]byte(expectedDigest), []byte(key.SecretDigest)) != 1 {
		return nil, ErrUnknownKey
	}

	owner, err := s.repo.GetOwner(ctx, key.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrOwnerInactive
		}

		return nil, err
	}

	if owner.Status != platformauth.StatusActive {
		return nil, ErrOwnerInactive
	}

	threshold := now.Add(-lastUsedUpdateAfter)
	if err := s.repo.UpdateLastUsed(ctx, key.ID, now, threshold); err != nil {
		return nil, err
	}
	key.LastUsedAt = sql.NullTime{Time: now, Valid: true}

	return &AuthenticatedKey{Key: key, Owner: owner}, nil
}

func (s *Service) digest(id, secret string) string {
	mac := hmac.New(sha256.New, []byte(s.pepper))
	_, _ = mac.Write([]byte(id))
	_, _ = mac.Write([]byte(":"))
	_, _ = mac.Write([]byte(secret))

	return hex.EncodeToString(mac.Sum(nil))
}

func validateName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrMissingName
	}

	if len(name) > maxNameLength {
		return ErrNameTooLong
	}

	return nil
}

func validateScopes(scopes []Scope) error {
	if len(scopes) == 0 {
		return ErrMissingScopes
	}

	seen := make(map[Scope]struct{}, len(scopes))
	for _, scope := range scopes {
		if _, ok := allowedScopes()[scope]; !ok {
			return fmt.Errorf("%w: %s", ErrInvalidScope, scope)
		}

		if _, ok := seen[scope]; ok {
			continue
		}
		seen[scope] = struct{}{}
	}

	return nil
}

func allowedScopes() map[Scope]struct{} {
	return map[Scope]struct{}{
		ScopeSitesRead:    {},
		ScopeSitesWrite:   {},
		ScopeSitesUpload:  {},
		ScopeDomainsRead:  {},
		ScopeDomainsWrite: {},
	}
}

func generateSecret() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func keyPrefix(id string) string {
	return "hmat_" + id
}

func formatSecret(id, secret string) string {
	return keyPrefix(id) + "_" + secret
}

func parseSecret(value string) (string, string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", "", ErrMissingBearer
	}

	parts := strings.SplitN(value, "_", 3)
	if len(parts) != 3 || parts[0] != "hmat" {
		return "", "", ErrMalformedBearer
	}

	if _, err := uuid.Parse(parts[1]); err != nil {
		return "", "", ErrMalformedBearer
	}

	if parts[2] == "" {
		return "", "", ErrMalformedBearer
	}

	return parts[1], parts[2], nil
}

func timePtrToNull(value *time.Time) sql.NullTime {
	if value == nil {
		return sql.NullTime{}
	}

	return sql.NullTime{Time: *value, Valid: true}
}
