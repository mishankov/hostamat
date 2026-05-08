package formcollection

import (
	"bytes"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/platforma-dev/platforma/auth"
)

type repository interface {
	UserOwnsSite(ctx context.Context, userID, siteID string) (bool, error)
	CreateDiscovery(ctx context.Context, discovery *Discovery, forms []DiscoveredForm) error
	DisableActiveBySite(ctx context.Context, siteID string) error
	GetLatestDiscoveryID(ctx context.Context, siteID string) (string, error)
	ListDefinitions(ctx context.Context, userID, siteID string) ([]DiscoveredForm, error)
	GetDefinitionsForEnable(ctx context.Context, userID, siteID, discoveryID string, ids []string) ([]DiscoveredForm, error)
	ActivateDefinition(ctx context.Context, id, token string) error
	GetActiveByToken(ctx context.Context, token string) (*DiscoveredForm, error)
	CreateSubmission(ctx context.Context, submission *Submission) error
	ListSubmissions(ctx context.Context, userID string, filters SubmissionFilters) ([]Submission, error)
	DeleteSubmission(ctx context.Context, userID, submissionID string) (bool, error)
	DeleteBySite(ctx context.Context, siteID string) error
	DeleteByUser(ctx context.Context, userID string) error
}

type fileService interface {
	GetFile(ctx context.Context, path string) (io.Reader, string, error)
	PutFile(ctx context.Context, folder, filename string, file io.Reader, size int64) error
}

type Service struct {
	repo        repository
	fileService fileService
	now         func() time.Time
}

func NewService(repo repository, fileService fileService) *Service {
	return &Service{
		repo:        repo,
		fileService: fileService,
		now:         time.Now,
	}
}

func (s *Service) RecordDiscovery(ctx context.Context, siteID string, files []UploadedHTMLFile) (*DiscoverySummary, error) {
	now := s.now()
	discovery := &Discovery{
		ID:        uuid.NewString(),
		SiteID:    siteID,
		CreatedAt: now,
	}

	forms := make([]DiscoveredForm, 0)
	for _, file := range files {
		discovered, err := discoverForms(file)
		if err != nil {
			return nil, err
		}

		for i := range discovered {
			discovered[i].ID = uuid.NewString()
			discovered[i].DiscoveryID = discovery.ID
			discovered[i].SiteID = siteID
			discovered[i].CreatedAt = now
			discovered[i].UpdatedAt = now
			forms = append(forms, discovered[i])
		}
	}

	if err := s.repo.CreateDiscovery(ctx, discovery, forms); err != nil {
		return nil, err
	}

	return buildSummary(discovery.ID, forms), nil
}

func (s *Service) DisableActiveBySite(ctx context.Context, siteID string) error {
	return s.repo.DisableActiveBySite(ctx, siteID)
}

func (s *Service) EnableCollection(ctx context.Context, siteID, discoveryID string, formIDs []string) (*DiscoverySummary, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, errAccessDenied
	}
	if len(formIDs) == 0 {
		return nil, errNoFormsSelected
	}

	owns, err := s.repo.UserOwnsSite(ctx, user.ID, siteID)
	if err != nil {
		return nil, err
	}
	if !owns {
		return nil, errAccessDenied
	}

	latestDiscoveryID, err := s.repo.GetLatestDiscoveryID(ctx, siteID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errFormNotFound
		}
		return nil, err
	}
	if latestDiscoveryID != discoveryID {
		return nil, errFormNotFound
	}

	forms, err := s.repo.GetDefinitionsForEnable(ctx, user.ID, siteID, discoveryID, formIDs)
	if err != nil {
		return nil, err
	}
	if len(forms) != len(uniqueStrings(formIDs)) {
		return nil, errFormNotFound
	}

	for _, form := range forms {
		if !form.Supported {
			return nil, errFormNotSupported
		}
	}

	tokensByFormID := make(map[string]string, len(forms))
	targetsByFile := make(map[string][]rewriteTarget)
	pageByFile := make(map[string]string)
	for _, form := range forms {
		token, err := newPublicToken()
		if err != nil {
			return nil, err
		}
		tokensByFormID[form.ID] = token
		targetsByFile[form.FileName] = append(targetsByFile[form.FileName], rewriteTarget{
			FormKey: form.FormKey,
			Action:  "/forms/submit/" + token,
		})
		pageByFile[form.FileName] = form.PagePath
	}

	for fileName, targets := range targetsByFile {
		reader, _, err := s.fileService.GetFile(ctx, siteID+"/"+fileName)
		if err != nil {
			return nil, err
		}
		body, err := readAll(reader)
		if err != nil {
			return nil, err
		}
		rewritten, err := rewriteForms(pageByFile[fileName], body, targets)
		if err != nil {
			return nil, err
		}
		if err := s.fileService.PutFile(ctx, siteID, fileName, bytes.NewReader(rewritten), int64(len(rewritten))); err != nil {
			return nil, err
		}
	}

	for i := range forms {
		if err := s.repo.ActivateDefinition(ctx, forms[i].ID, tokensByFormID[forms[i].ID]); err != nil {
			return nil, err
		}
		forms[i].Active = true
		forms[i].PublicToken = tokensByFormID[forms[i].ID]
	}

	return buildSummary(discoveryID, forms), nil
}

func (s *Service) ListDefinitions(ctx context.Context, siteID string) ([]DiscoveredForm, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, errAccessDenied
	}

	return s.repo.ListDefinitions(ctx, user.ID, siteID)
}

func (s *Service) Submit(ctx context.Context, token string, values url.Values) (*Submission, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, errFormNotFound
	}

	form, err := s.repo.GetActiveByToken(ctx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errFormNotFound
		}
		return nil, err
	}
	if !form.Active {
		return nil, errFormInactive
	}

	fields, err := submissionFields(values)
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(fields)
	if err != nil {
		return nil, err
	}

	submission := &Submission{
		ID:        uuid.NewString(),
		SiteID:    form.SiteID,
		FormID:    form.ID,
		PagePath:  form.PagePath,
		Fields:    payload,
		CreatedAt: s.now(),
	}

	if err := s.repo.CreateSubmission(ctx, submission); err != nil {
		return nil, err
	}

	return submission, nil
}

func (s *Service) ListSubmissions(ctx context.Context, filters SubmissionFilters) ([]Submission, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, errAccessDenied
	}

	return s.repo.ListSubmissions(ctx, user.ID, filters)
}

func (s *Service) DeleteSubmission(ctx context.Context, submissionID string) error {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return errAccessDenied
	}

	deleted, err := s.repo.DeleteSubmission(ctx, user.ID, submissionID)
	if err != nil {
		return err
	}
	if !deleted {
		return errFormNotFound
	}

	return nil
}

func (s *Service) DeleteBySite(ctx context.Context, siteID string) error {
	return s.repo.DeleteBySite(ctx, siteID)
}

func (s *Service) DeleteByUser(ctx context.Context, userID string) error {
	return s.repo.DeleteByUser(ctx, userID)
}

func buildSummary(discoveryID string, forms []DiscoveredForm) *DiscoverySummary {
	summary := &DiscoverySummary{
		DiscoveryID: discoveryID,
		TotalForms:  len(forms),
		Forms:       forms,
	}
	for _, form := range forms {
		if form.Supported {
			summary.SupportedForms++
		} else {
			summary.UnsupportedForms++
		}
	}
	return summary
}

func submissionFields(values url.Values) (map[string]any, error) {
	if len(values) > MaxSubmissionFields {
		return nil, errTooManyFields
	}

	fields := make(map[string]any, len(values))
	for name, value := range values {
		if strings.TrimSpace(name) == "" {
			continue
		}
		if len(value) == 1 {
			if len([]rune(value[0])) > MaxSubmissionValueLength {
				return nil, errFieldValueTooLong
			}
			fields[name] = value[0]
			continue
		}

		copied := make([]string, len(value))
		for i, v := range value {
			if len([]rune(v)) > MaxSubmissionValueLength {
				return nil, errFieldValueTooLong
			}
			copied[i] = v
		}
		fields[name] = copied
	}

	return fields, nil
}

func newPublicToken() (string, error) {
	var raw [32]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw[:]), nil
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	unique := make([]string, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		unique = append(unique, value)
	}
	return unique
}

func isAccessDenied(err error) bool {
	return errors.Is(err, errAccessDenied)
}
