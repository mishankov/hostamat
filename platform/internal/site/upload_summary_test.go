package site_test

import (
	"slices"
	"testing"

	"hostamat/internal/site"
	"hostamat/internal/uploadsession"
)

func TestBuildUploadSummaryDiscoversPagesUsingServeFilePrecedence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		files       []uploadsession.File
		wantHome    bool
		wantPages   []site.UploadedPage
		wantSkipped []string
	}{
		{
			name: "root index maps to site root without trailing slash in url",
			files: []uploadsession.File{
				{Name: "index.html", Uploaded: true, Size: 12},
			},
			wantHome: true,
			wantPages: []site.UploadedPage{
				{FileName: "index.html", Path: "/", URL: "https://example.hostamat.test"},
			},
		},
		{
			name: "simple html maps to clean extensionless path",
			files: []uploadsession.File{
				{Name: "about.html", Uploaded: true, Size: 12},
			},
			wantPages: []site.UploadedPage{
				{FileName: "about.html", Path: "/about", URL: "https://example.hostamat.test/about"},
			},
		},
		{
			name: "nested index maps to clean folder path when it owns the route",
			files: []uploadsession.File{
				{Name: "about/index.html", Uploaded: true, Size: 12},
			},
			wantPages: []site.UploadedPage{
				{FileName: "about/index.html", Path: "/about", URL: "https://example.hostamat.test/about"},
			},
		},
		{
			name: "html conflict follows clean route resolution order",
			files: []uploadsession.File{
				{Name: "about.html", Uploaded: true, Size: 12},
				{Name: "about/index.html", Uploaded: true, Size: 12},
			},
			wantPages: []site.UploadedPage{
				{FileName: "about.html", Path: "/about", URL: "https://example.hostamat.test/about"},
				{FileName: "about/index.html", Path: "/about/index.html", URL: "https://example.hostamat.test/about/index.html"},
			},
		},
		{
			name: "uppercase html extension is discovered with exact path",
			files: []uploadsession.File{
				{Name: "Docs.HTML", Uploaded: true, Size: 12},
			},
			wantPages: []site.UploadedPage{
				{FileName: "Docs.HTML", Path: "/Docs.HTML", URL: "https://example.hostamat.test/Docs.HTML"},
			},
		},
		{
			name: "url path is escaped while display path remains readable",
			files: []uploadsession.File{
				{Name: "my page.html", Uploaded: true, Size: 12},
			},
			wantPages: []site.UploadedPage{
				{FileName: "my page.html", Path: "/my page", URL: "https://example.hostamat.test/my%20page"},
			},
		},
		{
			name: "skipped html and assets are not pages",
			files: []uploadsession.File{
				{Name: "index.html", Uploaded: false, Size: 0},
				{Name: "style.css", Uploaded: true, Size: 12},
				{Name: "script.js", Uploaded: true, Size: 12},
			},
			wantSkipped: []string{"index.html"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			summary := site.BuildUploadSummary("example.hostamat.test", tt.files)

			if summary.HasHomePage != tt.wantHome {
				t.Fatalf("HasHomePage = %v, want %v", summary.HasHomePage, tt.wantHome)
			}

			if !slices.Equal(summary.AvailablePages, tt.wantPages) {
				t.Fatalf("AvailablePages = %#v, want %#v", summary.AvailablePages, tt.wantPages)
			}

			gotSkipped := make([]string, 0, len(summary.SkippedFiles))
			for _, file := range summary.SkippedFiles {
				gotSkipped = append(gotSkipped, file.Name)
			}

			if !slices.Equal(gotSkipped, tt.wantSkipped) {
				t.Fatalf("SkippedFiles = %v, want %v", gotSkipped, tt.wantSkipped)
			}
		})
	}
}
