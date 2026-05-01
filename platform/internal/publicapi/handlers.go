package publicapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"hostamat/internal/apikey"
	"hostamat/internal/customdomain"
	"hostamat/internal/site"
)

const maxZipUploadSize = 20 * 1024 * 1024

type siteService interface {
	Create(ctx context.Context, name string) (*site.Site, error)
	Get(ctx context.Context, id string) (*site.Site, error)
	GetList(ctx context.Context) ([]*site.Site, error)
	UpdateName(ctx context.Context, id, name string) error
	Delete(ctx context.Context, id string) error
	UploadZipArchive(ctx context.Context, siteID string, fileReader io.Reader, size int64) ([]site.UploadedFile, error)
}

type domainService interface {
	List(ctx context.Context) ([]customdomain.DomainRecord, error)
	Get(ctx context.Context, domainID string) (*customdomain.DomainRecord, error)
	Create(ctx context.Context, siteID, hostname string) (*customdomain.DomainRecord, error)
	UpdateSite(ctx context.Context, domainID, siteID string) (*customdomain.DomainRecord, error)
	Verify(ctx context.Context, domainID string) (*customdomain.DomainRecord, error)
	Delete(ctx context.Context, domainID string) error
}

type Handler struct {
	sites   siteService
	domains domainService
}

func NewHandler(sites siteService, domains domainService) *Handler {
	return &Handler{sites: sites, domains: domains}
}

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	authenticated := apikey.AuthenticatedKeyFromContext(r.Context())
	if authenticated == nil {
		RespondError(w, http.StatusUnauthorized, string(apikey.ErrorCodeMissingBearer), "missing bearer token")
		return
	}

	metadata := authenticated.Key.Metadata(authenticated.Key.UpdatedAt)
	respondJSON(w, http.StatusOK, meResponse{
		User: apiUser{
			ID:       authenticated.Owner.ID,
			Username: authenticated.Owner.Username,
		},
		APIKey: apiKey{
			ID:        authenticated.Key.ID,
			Name:      authenticated.Key.Name,
			KeyPrefix: authenticated.Key.KeyPrefix,
			Scopes:    metadata.Scopes,
			ExpiresAt: metadata.ExpiresAt,
		},
	})
}

func (h *Handler) listSites(w http.ResponseWriter, r *http.Request) {
	sites, err := h.sites.GetList(r.Context())
	if err != nil {
		ErrorWriter{}.WritePublicError(w, err)
		return
	}

	response := make([]siteResponse, 0, len(sites))
	for _, value := range sites {
		response = append(response, siteDTO(value))
	}

	respondJSON(w, http.StatusOK, sitesResponse{Sites: response})
}

func (h *Handler) createSite(w http.ResponseWriter, r *http.Request) {
	var request createSiteRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		RespondError(w, http.StatusBadRequest, "invalid_body", "Request body is invalid.")
		return
	}

	created, err := h.sites.Create(r.Context(), request.Name)
	if err != nil {
		ErrorWriter{}.WritePublicError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, siteEnvelopeResponse{Site: siteDTO(created)})
}

func (h *Handler) getSite(w http.ResponseWriter, r *http.Request) {
	value, err := h.sites.Get(r.Context(), r.PathValue("siteId"))
	if err != nil {
		ErrorWriter{}.WritePublicError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, siteEnvelopeResponse{Site: siteDTO(value)})
}

func (h *Handler) updateSite(w http.ResponseWriter, r *http.Request) {
	var request updateSiteRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		RespondError(w, http.StatusBadRequest, "invalid_body", "Request body is invalid.")
		return
	}

	if err := h.sites.UpdateName(r.Context(), r.PathValue("siteId"), request.Name); err != nil {
		ErrorWriter{}.WritePublicError(w, err)
		return
	}

	value, err := h.sites.Get(r.Context(), r.PathValue("siteId"))
	if err != nil {
		ErrorWriter{}.WritePublicError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, siteEnvelopeResponse{Site: siteDTO(value)})
}

func (h *Handler) deleteSite(w http.ResponseWriter, r *http.Request) {
	if err := h.sites.Delete(r.Context(), r.PathValue("siteId")); err != nil {
		ErrorWriter{}.WritePublicError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) uploadSite(w http.ResponseWriter, r *http.Request) {
	contentType := strings.ToLower(strings.TrimSpace(strings.Split(r.Header.Get("Content-Type"), ";")[0]))
	if contentType != "application/zip" {
		RespondError(w, http.StatusUnsupportedMediaType, "unsupported_media_type", "Upload body must be application/zip.")
		return
	}

	body := http.MaxBytesReader(w, r.Body, maxZipUploadSize+1)
	data, err := readAll(body)
	if err != nil {
		var maxBytesError *http.MaxBytesError
		if errors.As(err, &maxBytesError) {
			RespondError(w, http.StatusRequestEntityTooLarge, "upload_too_large", "Upload exceeds the maximum allowed size.")
			return
		}

		RespondError(w, http.StatusBadRequest, "invalid_upload", "Upload body is invalid.")
		return
	}

	if len(data) > maxZipUploadSize {
		RespondError(w, http.StatusRequestEntityTooLarge, "upload_too_large", "Upload exceeds the maximum allowed size.")
		return
	}

	files, err := h.sites.UploadZipArchive(r.Context(), r.PathValue("siteId"), bytes.NewReader(data), int64(len(data)))
	if err != nil {
		ErrorWriter{}.WritePublicError(w, err)
		return
	}

	summary := uploadSummary{}
	responseFiles := make([]uploadFileResponse, 0, len(files))
	for _, file := range files {
		status := "skipped"
		if file.Uploaded {
			summary.UploadedCount++
			status = "uploaded"
		} else {
			summary.SkippedCount++
		}

		summary.TotalBytes += int64(file.Size)
		responseFiles = append(responseFiles, uploadFileResponse{
			Path:   file.Name,
			Size:   int64(file.Size),
			Status: status,
		})
	}

	respondJSON(w, http.StatusOK, uploadResponse{Summary: summary, Files: responseFiles})
}

func (h *Handler) listDomains(w http.ResponseWriter, r *http.Request) {
	domains, err := h.domains.List(r.Context())
	if err != nil {
		ErrorWriter{}.WritePublicError(w, err)
		return
	}

	response := make([]domainResponse, 0, len(domains))
	for _, domain := range domains {
		if domain.Type == customdomain.DomainTypeCustom {
			response = append(response, domainDTO(domain))
		}
	}

	respondJSON(w, http.StatusOK, domainsResponse{Domains: response})
}

func (h *Handler) createDomain(w http.ResponseWriter, r *http.Request) {
	var request createDomainRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		RespondError(w, http.StatusBadRequest, "invalid_body", "Request body is invalid.")
		return
	}

	record, err := h.domains.Create(r.Context(), request.SiteID, request.Hostname)
	if err != nil {
		ErrorWriter{}.WritePublicError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, domainEnvelopeResponse{Domain: domainDTO(*record)})
}

func (h *Handler) getDomain(w http.ResponseWriter, r *http.Request) {
	record, err := h.domains.Get(r.Context(), r.PathValue("domainId"))
	if err != nil {
		ErrorWriter{}.WritePublicError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, domainEnvelopeResponse{Domain: domainDTO(*record)})
}

func (h *Handler) updateDomain(w http.ResponseWriter, r *http.Request) {
	var request updateDomainRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		RespondError(w, http.StatusBadRequest, "invalid_body", "Request body is invalid.")
		return
	}

	record, err := h.domains.UpdateSite(r.Context(), r.PathValue("domainId"), request.SiteID)
	if err != nil {
		ErrorWriter{}.WritePublicError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, domainEnvelopeResponse{Domain: domainDTO(*record)})
}

func (h *Handler) verifyDomain(w http.ResponseWriter, r *http.Request) {
	record, err := h.domains.Verify(r.Context(), r.PathValue("domainId"))
	if err != nil {
		ErrorWriter{}.WritePublicError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, domainEnvelopeResponse{Domain: domainDTO(*record)})
}

func (h *Handler) deleteDomain(w http.ResponseWriter, r *http.Request) {
	if err := h.domains.Delete(r.Context(), r.PathValue("domainId")); err != nil {
		ErrorWriter{}.WritePublicError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func respondJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func readAll(r io.Reader) ([]byte, error) {
	buf := bytes.Buffer{}
	if _, err := buf.ReadFrom(r); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
