package site

import (
	"net/url"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"hostamat/internal/uploadsession"
)

func BuildUploadSummary(domain string, files []uploadsession.File) *UploadSummary {
	uploadedFiles := make([]uploadsession.File, 0, len(files))
	skippedFiles := make([]uploadsession.File, 0)
	uploadedNames := make(map[string]struct{})

	for _, file := range files {
		if file.Uploaded {
			uploadedFiles = append(uploadedFiles, file)
			uploadedNames[normalizeUploadedFileName(file.Name)] = struct{}{}
			continue
		}

		skippedFiles = append(skippedFiles, file)
	}

	pages := make([]UploadedPage, 0)
	hasHomePage := false

	for _, file := range uploadedFiles {
		fileName := normalizeUploadedFileName(file.Name)
		if !isHTMLFile(fileName) {
			continue
		}

		if fileName == "index.html" {
			hasHomePage = true
		}

		pagePath := pagePathForHTMLFile(fileName, uploadedNames)
		pages = append(pages, UploadedPage{
			FileName: file.Name,
			Path:     pagePath,
			URL:      pageURL(domain, pagePath),
		})
	}

	slices.SortFunc(pages, func(a, b UploadedPage) int {
		return strings.Compare(a.Path, b.Path)
	})

	return &UploadSummary{
		UploadedFiles:  uploadedFiles,
		SkippedFiles:   skippedFiles,
		AvailablePages: pages,
		HasHomePage:    hasHomePage,
	}
}

func pagePathForHTMLFile(fileName string, uploadedNames map[string]struct{}) string {
	if fileName == "index.html" {
		return "/"
	}

	cleanPath := cleanPagePath(fileName)
	if cleanPathResolvesToFile(cleanPath, fileName, uploadedNames) {
		return "/" + cleanPath
	}

	return "/" + fileName
}

func cleanPagePath(fileName string) string {
	if strings.HasSuffix(fileName, "/index.html") {
		return strings.TrimSuffix(fileName, "/index.html")
	}

	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func cleanPathResolvesToFile(cleanPath, fileName string, uploadedNames map[string]struct{}) bool {
	for _, candidate := range []string{cleanPath, cleanPath + ".html", cleanPath + "/index.html"} {
		if _, ok := uploadedNames[candidate]; !ok {
			continue
		}

		return candidate == fileName
	}

	return false
}

func pageURL(domain, pagePath string) string {
	baseURL := "https://" + strings.TrimRight(domain, "/")
	if pagePath == "/" {
		return baseURL
	}

	return baseURL + (&url.URL{Path: pagePath}).EscapedPath()
}

func normalizeUploadedFileName(fileName string) string {
	fileName = strings.ReplaceAll(fileName, "\\", "/")
	fileName = strings.TrimSpace(fileName)
	fileName = strings.TrimPrefix(fileName, "/")
	fileName = path.Clean(fileName)
	if fileName == "." {
		return ""
	}

	return fileName
}

func isHTMLFile(fileName string) bool {
	return strings.EqualFold(filepath.Ext(fileName), ".html")
}
