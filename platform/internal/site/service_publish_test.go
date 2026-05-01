package site

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"hostamat/internal/file"
	"hostamat/internal/subdomain"
	"hostamat/internal/uploadsession"

	"github.com/platforma-dev/platforma/auth"
	"github.com/platforma-dev/platforma/log"
)

func TestCreateUploadSessionDoesNotDeleteLiveFiles(t *testing.T) {
	t.Parallel()

	repo := newSiteRepoMock()
	files := &siteFileServiceMock{}
	uploads := newUploadSessionServiceMock()
	service := NewService(repo, nil, nil, files, uploads)

	us, err := service.CreateUploadSession(testContext(), "site-1")
	if err != nil {
		t.Fatalf("CreateUploadSession() error = %v", err)
	}

	if us == nil {
		t.Fatal("CreateUploadSession() session = nil")
	}

	if files.deleteAsyncCalls != 0 {
		t.Fatalf("DeleteFromFolderAsync() calls = %d, want 0", files.deleteAsyncCalls)
	}

	if repo.activeRelease != "" {
		t.Fatalf("active release = %q, want empty", repo.activeRelease)
	}
}

func TestUploadFileWritesToSessionReleasePrefix(t *testing.T) {
	t.Parallel()

	repo := newSiteRepoMock()
	files := &siteFileServiceMock{
		uploaded: []file.File{{Name: "index.html", Size: 12}},
	}
	uploads := newUploadSessionServiceMock()
	uploads.sessions["session-1"] = &uploadsession.UploadSession{
		Id:     "session-1",
		Site:   "site-1",
		User:   "user-1",
		Status: uploadsession.UploadSessionStatusOpen,
	}

	service := NewService(repo, nil, nil, files, uploads)

	err := service.UploadFile(testContext(), "session-1", "index.html", bytes.NewBufferString("hello"), 5)
	if err != nil {
		t.Fatalf("UploadFile() error = %v", err)
	}

	wantPath := "site-1/releases/session-1"
	if files.uploadPath != wantPath {
		t.Fatalf("UploadFile() path = %q, want %q", files.uploadPath, wantPath)
	}

	if len(uploads.files["session-1"]) != 1 {
		t.Fatalf("recorded files = %d, want 1", len(uploads.files["session-1"]))
	}
}

func TestUploadFileZipUsesRemainingUncompressedSizeLimit(t *testing.T) {
	t.Parallel()

	repo := newSiteRepoMock()
	files := &siteFileServiceMock{
		uploaded: []file.File{{Name: "index.html", Size: 12}},
	}
	uploads := newUploadSessionServiceMock()
	uploads.sessions["session-1"] = &uploadsession.UploadSession{
		Id:     "session-1",
		Site:   "site-1",
		User:   "user-1",
		Status: uploadsession.UploadSessionStatusOpen,
	}
	uploads.files["session-1"] = []uploadsession.File{{Name: "existing.css", Uploaded: true, Size: 7}}

	service := NewService(repo, nil, nil, files, uploads)

	err := service.UploadFile(testContext(), "session-1", "archive.zip", bytes.NewReader([]byte("zip")), 3)
	if err != nil {
		t.Fatalf("UploadFile() error = %v", err)
	}

	if files.zipUploadPath != "site-1/releases/session-1" {
		t.Fatalf("UploadFile() zip path = %q, want site-1/releases/session-1", files.zipUploadPath)
	}

	wantMaxSize := int64(maxFilesSize - 7)
	if files.zipUploadMaxSize != wantMaxSize {
		t.Fatalf("UploadFile() zip max size = %d, want %d", files.zipUploadMaxSize, wantMaxSize)
	}
}

func TestCloseUploadSessionPublishesActiveRelease(t *testing.T) {
	t.Parallel()

	repo := newSiteRepoMock()
	files := &siteFileServiceMock{}
	uploads := newUploadSessionServiceMock()
	uploads.sessions["session-1"] = &uploadsession.UploadSession{
		Id:     "session-1",
		Site:   "site-1",
		User:   "user-1",
		Status: uploadsession.UploadSessionStatusOpen,
	}
	uploads.files["session-1"] = []uploadsession.File{{Name: "index.html", Uploaded: true, Size: 12}}

	service := NewService(repo, nil, nil, files, uploads)

	returnedFiles, err := service.CloseUploadSession(testContext(), "session-1")
	if err != nil {
		t.Fatalf("CloseUploadSession() error = %v", err)
	}

	if repo.activeRelease != "session-1" {
		t.Fatalf("active release = %q, want %q", repo.activeRelease, "session-1")
	}

	if uploads.sessions["session-1"].Status != uploadsession.UploadSessionStatusPublished {
		t.Fatalf("session status = %q, want %q", uploads.sessions["session-1"].Status, uploadsession.UploadSessionStatusPublished)
	}

	if len(returnedFiles) != 1 {
		t.Fatalf("returned files = %d, want 1", len(returnedFiles))
	}
}

func TestUploadZipArchiveFailureLeavesActiveReleaseUnchanged(t *testing.T) {
	t.Parallel()

	repo := newSiteRepoMock()
	repo.activeRelease = "old-release"
	files := &siteFileServiceMock{
		uploadErr: errors.New("invalid zip"),
	}
	uploads := newUploadSessionServiceMock()
	service := NewService(repo, nil, nil, files, uploads)

	_, err := service.UploadZipArchive(testContext(), "site-1", bytes.NewReader([]byte("bad zip")), 7)
	if err == nil {
		t.Fatal("UploadZipArchive() error = nil, want error")
	}

	if repo.activeRelease != "old-release" {
		t.Fatalf("active release = %q, want old-release", repo.activeRelease)
	}

	if len(uploads.sessions) != 1 {
		t.Fatalf("sessions = %d, want 1", len(uploads.sessions))
	}

	for _, us := range uploads.sessions {
		if us.Status != uploadsession.UploadSessionStatusFailed {
			t.Fatalf("session status = %q, want %q", us.Status, uploadsession.UploadSessionStatusFailed)
		}
	}

	if files.deleteAsyncCalls != 1 {
		t.Fatalf("DeleteFromFolderAsync() calls = %d, want 1", files.deleteAsyncCalls)
	}
}

func TestUploadZipArchiveUsesUncompressedSizeLimit(t *testing.T) {
	t.Parallel()

	repo := newSiteRepoMock()
	files := &siteFileServiceMock{
		uploaded: []file.File{{Name: "index.html", Size: 12}},
	}
	uploads := newUploadSessionServiceMock()
	service := NewService(repo, nil, nil, files, uploads)

	_, err := service.UploadZipArchive(testContext(), "site-1", bytes.NewReader([]byte("zip")), 3)
	if err != nil {
		t.Fatalf("UploadZipArchive() error = %v", err)
	}

	if files.zipUploadPath != "site-1/releases/session-1" {
		t.Fatalf("UploadZipArchive() path = %q, want site-1/releases/session-1", files.zipUploadPath)
	}

	if files.zipUploadMaxSize != maxFilesSize {
		t.Fatalf("UploadZipArchive() max size = %d, want %d", files.zipUploadMaxSize, maxFilesSize)
	}
}

func TestUploadZipArchiveMapsExpandedSizeError(t *testing.T) {
	t.Parallel()

	repo := newSiteRepoMock()
	repo.activeRelease = "old-release"
	files := &siteFileServiceMock{
		zipUploadErr: file.ErrUploadSizeExceeded,
	}
	uploads := newUploadSessionServiceMock()
	service := NewService(repo, nil, nil, files, uploads)

	_, err := service.UploadZipArchive(testContext(), "site-1", bytes.NewReader([]byte("zip")), 3)
	if !errors.Is(err, errUploadSize) {
		t.Fatalf("UploadZipArchive() error = %v, want %v", err, errUploadSize)
	}

	if repo.activeRelease != "old-release" {
		t.Fatalf("active release = %q, want old-release", repo.activeRelease)
	}
}

func TestCloseUploadSessionRejectsFailedSession(t *testing.T) {
	t.Parallel()

	repo := newSiteRepoMock()
	repo.activeRelease = "old-release"
	files := &siteFileServiceMock{}
	uploads := newUploadSessionServiceMock()
	uploads.sessions["session-1"] = &uploadsession.UploadSession{
		Id:     "session-1",
		Site:   "site-1",
		User:   "user-1",
		Status: uploadsession.UploadSessionStatusFailed,
	}

	service := NewService(repo, nil, nil, files, uploads)

	_, err := service.CloseUploadSession(testContext(), "session-1")
	if err == nil {
		t.Fatal("CloseUploadSession() error = nil, want error")
	}

	if repo.activeRelease != "old-release" {
		t.Fatalf("active release = %q, want old-release", repo.activeRelease)
	}
}

func TestGetFileReadsFromActiveRelease(t *testing.T) {
	t.Parallel()

	repo := newSiteRepoMock()
	release := "release-1"
	repo.site.ActiveReleaseId = &release
	files := &siteFileServiceMock{
		reader:       bytes.NewBufferString("ok"),
		resolvedPath: "site-1/releases/release-1/about.html",
	}
	uploads := newUploadSessionServiceMock()
	service := NewService(repo, nil, nil, files, uploads)

	_, resolvedPath, err := service.GetFile(testContext(), "site.test", "/about")
	if err != nil {
		t.Fatalf("GetFile() error = %v", err)
	}

	wantRequestPath := "site-1/releases/release-1/about"
	if files.getPath != wantRequestPath {
		t.Fatalf("GetFile() requested path = %q, want %q", files.getPath, wantRequestPath)
	}

	if resolvedPath != "about.html" {
		t.Fatalf("GetFile() resolvedPath = %q, want %q", resolvedPath, "about.html")
	}
}

type siteRepoMock struct {
	site          *Site
	activeRelease string
}

func newSiteRepoMock() *siteRepoMock {
	now := time.Now()
	return &siteRepoMock{
		site: &Site{
			Id:        "site-1",
			Name:      "Site",
			User:      "user-1",
			Subdomain: "subdomain-1",
			Domain:    "site.test",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

func (m *siteRepoMock) Get(ctx context.Context, id string) (*Site, error) {
	return m.site, nil
}

func (m *siteRepoMock) GetByDomain(ctx context.Context, domain string) (*Site, error) {
	return m.site, nil
}

func (m *siteRepoMock) GetList(ctx context.Context, userId string) ([]*Site, error) {
	return []*Site{m.site}, nil
}

func (m *siteRepoMock) Create(ctx context.Context, site *Site) error {
	m.site = site
	return nil
}

func (m *siteRepoMock) Update(ctx context.Context, site *Site) error {
	m.site = site
	return nil
}

func (m *siteRepoMock) UpdateActiveRelease(ctx context.Context, siteId, releaseId string) error {
	m.activeRelease = releaseId
	m.site.ActiveReleaseId = &releaseId
	return nil
}

func (m *siteRepoMock) Delete(ctx context.Context, id string) error {
	return nil
}

type siteFileServiceMock struct {
	uploadPath       string
	zipUploadPath    string
	zipUploadMaxSize int64
	uploadErr        error
	zipUploadErr     error
	uploaded         []file.File
	skipped          []file.File
	deleteAsyncCalls int
	getPath          string
	reader           io.Reader
	resolvedPath     string
}

func (m *siteFileServiceMock) GetFile(ctx context.Context, path string) (io.Reader, string, error) {
	m.getPath = path
	return m.reader, m.resolvedPath, nil
}

func (m *siteFileServiceMock) UploadFile(ctx context.Context, path, filename string, reader io.Reader, size int64) ([]file.File, []file.File, error) {
	m.uploadPath = path
	if m.uploadErr != nil {
		return nil, nil, m.uploadErr
	}
	return m.uploaded, m.skipped, nil
}

func (m *siteFileServiceMock) UploadZipArchive(ctx context.Context, path string, reader io.Reader, size, maxUncompressedSize int64) ([]file.File, []file.File, error) {
	m.zipUploadPath = path
	m.zipUploadMaxSize = maxUncompressedSize
	if m.zipUploadErr != nil {
		return nil, nil, m.zipUploadErr
	}
	if m.uploadErr != nil {
		return nil, nil, m.uploadErr
	}
	return m.uploaded, m.skipped, nil
}

func (m *siteFileServiceMock) DeleteFromFolderAsync(ctx context.Context, folder string) error {
	m.deleteAsyncCalls++
	return nil
}

type uploadSessionServiceMock struct {
	sessions map[string]*uploadsession.UploadSession
	files    map[string][]uploadsession.File
}

func newUploadSessionServiceMock() *uploadSessionServiceMock {
	return &uploadSessionServiceMock{
		sessions: make(map[string]*uploadsession.UploadSession),
		files:    make(map[string][]uploadsession.File),
	}
}

func (m *uploadSessionServiceMock) Create(ctx context.Context, siteId string) (*uploadsession.UploadSession, error) {
	id := "session-1"
	if _, ok := m.sessions[id]; ok {
		id = "session-2"
	}

	us := &uploadsession.UploadSession{
		Id:        id,
		User:      "user-1",
		Site:      siteId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    uploadsession.UploadSessionStatusOpen,
	}
	m.sessions[id] = us
	return us, nil
}

func (m *uploadSessionServiceMock) Get(ctx context.Context, id string) (*uploadsession.UploadSession, error) {
	us, ok := m.sessions[id]
	if !ok {
		return nil, errors.New("upload session not found")
	}
	return us, nil
}

func (m *uploadSessionServiceMock) GetSize(ctx context.Context, uploadSessionId string) (int, error) {
	var size int
	for _, file := range m.files[uploadSessionId] {
		size += file.Size
	}
	return size, nil
}

func (m *uploadSessionServiceMock) CreateFile(ctx context.Context, uploadSessionId, filename string, size int, uploaded bool) error {
	m.files[uploadSessionId] = append(m.files[uploadSessionId], uploadsession.File{
		Name:     filename,
		Uploaded: uploaded,
		Size:     size,
	})
	return nil
}

func (m *uploadSessionServiceMock) GetFiles(ctx context.Context, uploadSessionId string) ([]uploadsession.File, error) {
	return m.files[uploadSessionId], nil
}

func (m *uploadSessionServiceMock) Close(ctx context.Context, uploadSessionId string) error {
	us, err := m.Get(ctx, uploadSessionId)
	if err != nil {
		return err
	}
	us.Status = uploadsession.UploadSessionStatusPublished
	delete(m.files, uploadSessionId)
	return nil
}

func (m *uploadSessionServiceMock) SetStatus(ctx context.Context, uploadSessionId string, status uploadsession.UploadSessionStatus) error {
	us, err := m.Get(ctx, uploadSessionId)
	if err != nil {
		return err
	}
	us.Status = status
	return nil
}

func testContext() context.Context {
	user := &auth.User{ID: "user-1"}
	ctx := context.WithValue(context.Background(), log.WideEventKey, log.NewEvent("test"))
	ctx = context.WithValue(ctx, auth.UserContextKey, user)
	ctx = context.WithValue(ctx, log.UserIDKey, user.ID)
	return ctx
}

var _ subdomainService = (*subdomainServiceMock)(nil)

type subdomainServiceMock struct{}

func (m *subdomainServiceMock) Create(ctx context.Context) (*subdomain.Subdomain, error) {
	return nil, nil
}

func (m *subdomainServiceMock) Delete(ctx context.Context, id string) error {
	return nil
}
