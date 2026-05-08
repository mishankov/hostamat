package formcollection_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
	"testing"
	"time"

	"hostamat/internal/formcollection"

	"github.com/platforma-dev/platforma/auth"
)

func TestRecordDiscoveryFindsSupportedAndUnsupportedForms(t *testing.T) {
	t.Parallel()

	repo := &repositoryMock{ownsSite: true}
	service := formcollection.NewService(repo, &fileServiceMock{})

	summary, err := service.RecordDiscovery(context.Background(), "site-1", []formcollection.UploadedHTMLFile{
		{
			FileName: "index.html",
			PagePath: "/",
			Body: []byte(`<!doctype html>
				<form aria-label="Contact" action="https://forms.example/old">
					<input name="email" type="email">
					<textarea name="message"></textarea>
					<input type="file" name="attachment">
				</form>
				<form method="post"><input name="name"><button>Send</button></form>
				<form method="put"><input name="ignored"></form>`),
		},
	})
	if err != nil {
		t.Fatalf("RecordDiscovery() error = %v", err)
	}

	if summary.TotalForms != 3 || summary.SupportedForms != 1 || summary.UnsupportedForms != 2 {
		t.Fatalf("summary counts = %#v", summary)
	}

	if len(repo.forms) != 3 {
		t.Fatalf("stored forms = %d, want 3", len(repo.forms))
	}
	if repo.forms[0].UnsupportedReason != "file_input" {
		t.Fatalf("first form unsupported reason = %q, want file_input", repo.forms[0].UnsupportedReason)
	}
	if repo.forms[1].Label != "Send" {
		t.Fatalf("second form label = %q, want Send", repo.forms[1].Label)
	}
	if repo.forms[2].UnsupportedReason != "unsupported_method" {
		t.Fatalf("third form unsupported reason = %q, want unsupported_method", repo.forms[2].UnsupportedReason)
	}
}

func TestEnableCollectionRewritesOnlySelectedFormsAndOverwritesAction(t *testing.T) {
	t.Parallel()

	source := `<!doctype html>
		<form id="selected" action="https://third-party.example/lead"><input name="email"></form>
		<form id="other" action="https://third-party.example/keep"><input name="name"></form>`
	repo := &repositoryMock{ownsSite: true}
	files := &fileServiceMock{
		files: map[string][]byte{
			"site-1/index.html": []byte(source),
		},
	}
	service := formcollection.NewService(repo, files)

	summary, err := service.RecordDiscovery(authContext("user-1"), "site-1", []formcollection.UploadedHTMLFile{
		{FileName: "index.html", PagePath: "/", Body: []byte(source)},
	})
	if err != nil {
		t.Fatalf("RecordDiscovery() error = %v", err)
	}

	enabled, err := service.EnableCollection(authContext("user-1"), "site-1", summary.DiscoveryID, []string{summary.Forms[0].ID})
	if err != nil {
		t.Fatalf("EnableCollection() error = %v", err)
	}
	if enabled.SupportedForms != 1 || !repo.active[summary.Forms[0].ID] {
		t.Fatalf("enabled summary = %#v, active = %#v", enabled, repo.active)
	}

	rewritten := string(files.files["site-1/index.html"])
	if !strings.Contains(rewritten, `id="selected"`) || !strings.Contains(rewritten, `action="/forms/submit/`) {
		t.Fatalf("rewritten html = %s, want selected form Hostamat action", rewritten)
	}
	if !strings.Contains(rewritten, `method="post"`) {
		t.Fatalf("rewritten html = %s, want post method", rewritten)
	}
	if strings.Contains(rewritten, "third-party.example/lead") {
		t.Fatalf("rewritten html = %s, selected old action was not overwritten", rewritten)
	}
	if !strings.Contains(rewritten, `id="other" action="https://third-party.example/keep"`) {
		t.Fatalf("rewritten html = %s, unselected form action changed", rewritten)
	}
}

func TestSubmitStoresOnlyJSONFieldsAndTimestamps(t *testing.T) {
	t.Parallel()

	repo := &repositoryMock{
		activeForm: &formcollection.DiscoveredForm{
			ID:       "form-1",
			SiteID:   "site-1",
			PagePath: "/contact",
			Active:   true,
		},
	}
	service := formcollection.NewService(repo, &fileServiceMock{})

	submission, err := service.Submit(context.Background(), "token-1", url.Values{
		"name":   {"Ada"},
		"topics": {"sales", "support"},
		"empty":  {""},
	})
	if err != nil {
		t.Fatalf("Submit() error = %v", err)
	}
	if submission.SiteID != "site-1" || submission.FormID != "form-1" || submission.PagePath != "/contact" {
		t.Fatalf("submission context = %#v", submission)
	}
	if len(repo.submissions) != 1 {
		t.Fatalf("stored submissions = %d, want 1", len(repo.submissions))
	}

	var fields map[string]any
	if err := json.Unmarshal(repo.submissions[0].Fields, &fields); err != nil {
		t.Fatalf("stored fields are not JSON: %v", err)
	}
	if fields["name"] != "Ada" || fields["empty"] != "" {
		t.Fatalf("stored scalar fields = %#v", fields)
	}
	if got := fields["topics"].([]any); len(got) != 2 || got[0] != "sales" || got[1] != "support" {
		t.Fatalf("stored duplicate fields = %#v", fields["topics"])
	}
}

func TestSubmitRejectsLimitsAndInactiveForms(t *testing.T) {
	t.Parallel()

	t.Run("too many fields", func(t *testing.T) {
		t.Parallel()
		repo := &repositoryMock{activeForm: &formcollection.DiscoveredForm{ID: "form-1", SiteID: "site-1", Active: true}}
		service := formcollection.NewService(repo, &fileServiceMock{})
		values := url.Values{}
		for i := 0; i <= formcollection.MaxSubmissionFields; i++ {
			values.Set(fmt.Sprintf("field_%d", i), "x")
		}
		if _, err := service.Submit(context.Background(), "token-1", values); err == nil {
			t.Fatal("Submit() error = nil, want field-count error")
		}
	})

	t.Run("inactive form", func(t *testing.T) {
		t.Parallel()
		repo := &repositoryMock{activeForm: &formcollection.DiscoveredForm{ID: "form-1", SiteID: "site-1", Active: false}}
		service := formcollection.NewService(repo, &fileServiceMock{})
		if _, err := service.Submit(context.Background(), "token-1", url.Values{"name": {"Ada"}}); err == nil {
			t.Fatal("Submit() error = nil, want inactive form error")
		}
	})
}

type repositoryMock struct {
	ownsSite          bool
	latestDiscoveryID string
	forms             []formcollection.DiscoveredForm
	active            map[string]bool
	tokens            map[string]string
	activeForm        *formcollection.DiscoveredForm
	submissions       []formcollection.Submission
}

func (m *repositoryMock) UserOwnsSite(ctx context.Context, userID, siteID string) (bool, error) {
	return m.ownsSite, nil
}

func (m *repositoryMock) CreateDiscovery(ctx context.Context, discovery *formcollection.Discovery, forms []formcollection.DiscoveredForm) error {
	m.latestDiscoveryID = discovery.ID
	m.forms = append(m.forms, forms...)
	return nil
}

func (m *repositoryMock) DisableActiveBySite(ctx context.Context, siteID string) error {
	for i := range m.forms {
		if m.forms[i].SiteID == siteID {
			m.forms[i].Active = false
		}
	}
	return nil
}

func (m *repositoryMock) GetLatestDiscoveryID(ctx context.Context, siteID string) (string, error) {
	return m.latestDiscoveryID, nil
}

func (m *repositoryMock) ListDefinitions(ctx context.Context, userID, siteID string) ([]formcollection.DiscoveredForm, error) {
	return m.forms, nil
}

func (m *repositoryMock) GetDefinitionsForEnable(ctx context.Context, userID, siteID, discoveryID string, ids []string) ([]formcollection.DiscoveredForm, error) {
	selected := make([]formcollection.DiscoveredForm, 0, len(ids))
	want := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		want[id] = struct{}{}
	}
	for _, form := range m.forms {
		if _, ok := want[form.ID]; ok && form.SiteID == siteID && form.DiscoveryID == discoveryID {
			selected = append(selected, form)
		}
	}
	return selected, nil
}

func (m *repositoryMock) ActivateDefinition(ctx context.Context, id, token string) error {
	if m.active == nil {
		m.active = make(map[string]bool)
	}
	if m.tokens == nil {
		m.tokens = make(map[string]string)
	}
	m.active[id] = true
	m.tokens[id] = token
	return nil
}

func (m *repositoryMock) GetActiveByToken(ctx context.Context, token string) (*formcollection.DiscoveredForm, error) {
	return m.activeForm, nil
}

func (m *repositoryMock) CreateSubmission(ctx context.Context, submission *formcollection.Submission) error {
	submission.CreatedAt = time.Now()
	m.submissions = append(m.submissions, *submission)
	return nil
}

func (m *repositoryMock) ListSubmissions(ctx context.Context, userID string, filters formcollection.SubmissionFilters) ([]formcollection.Submission, error) {
	return m.submissions, nil
}

func (m *repositoryMock) DeleteSubmission(ctx context.Context, userID, submissionID string) (bool, error) {
	return true, nil
}

func (m *repositoryMock) DeleteBySite(ctx context.Context, siteID string) error {
	return nil
}

func (m *repositoryMock) DeleteByUser(ctx context.Context, userID string) error {
	return nil
}

type fileServiceMock struct {
	files map[string][]byte
}

func (m *fileServiceMock) GetFile(ctx context.Context, path string) (io.Reader, string, error) {
	return bytes.NewReader(m.files[path]), path, nil
}

func (m *fileServiceMock) PutFile(ctx context.Context, folder, filename string, file io.Reader, size int64) error {
	if m.files == nil {
		m.files = make(map[string][]byte)
	}
	body, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	m.files[folder+"/"+filename] = body
	return nil
}

func authContext(userID string) context.Context {
	return context.WithValue(context.Background(), auth.UserContextKey, &auth.User{ID: userID})
}
