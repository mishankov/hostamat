package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"hostamat/internal/adapters/caddy"
	"hostamat/internal/adapters/config"
	"hostamat/internal/adapters/s3"
	"hostamat/internal/httpclient"

	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
	"github.com/platforma-dev/platforma/application"
	"github.com/platforma-dev/platforma/database"
	"github.com/platforma-dev/platforma/httpserver"
	"github.com/platforma-dev/platforma/log"
)

const cookieName = "dashsession"

type WipeCounts struct {
	Users          int
	Sessions       int
	Sites          int
	CustomDomains  int
	Subdomains     int
	UploadSessions int
	Files          int
	S3Objects      int
	CaddyHosts     int
}

//go:embed templates
var templates embed.FS

func Templates() fs.FS {
	serverRoot, err := fs.Sub(templates, "templates")
	if err != nil {
		panic(err)
	}

	return serverRoot
}

func main() {
	ctx := context.Background()
	app := application.New()
	authToken := os.Getenv("TOKEN")

	conf := config.New()
	log.InfoContext(ctx, "config", "config", conf)

	db, err := database.New(os.Getenv("DB_CONN"))
	if err != nil {
		log.ErrorContext(ctx, "failed to create database", "error", err)
	}
	conn := db.Connection()
	defer conn.Close()

	httpClient := httpclient.New(10 * time.Second)
	caddyAPI := caddy.New(httpClient, conf.Caddy.BaseURL, "srv0")

	s3Client, err := s3.New(conf.S3.BaseURL, conf.S3.Login, conf.S3.Password, false)
	if err != nil {
		log.ErrorContext(ctx, "failed to create S3 client", "error", err)
	}

	tmpls, err := template.ParseFS(Templates(), "login.html", "dashboard.html")
	if err != nil {
		log.ErrorContext(ctx, "failed to parse templates", "error", err)
	}

	server := httpserver.New("4444", 3*time.Second)
	server.Use(httpserver.NewRecoverMiddleware())

	server.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			err := tmpls.ExecuteTemplate(w, "login.html", nil)
			if err != nil {
				log.ErrorContext(ctx, "failed to execute template", "error", err)
			}
		case http.MethodPost:
			type req struct {
				Token string `json:"token"`
			}

			var data req
			json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				log.ErrorContext(ctx, "failed to decode request", "error", err)
				w.WriteHeader(500)
			}

			if data.Token == authToken {
				http.SetCookie(w, &http.Cookie{
					Name:     cookieName,
					Value:    authToken,
					Path:     "/",
					HttpOnly: true})
				w.WriteHeader(200)
			}

			w.WriteHeader(http.StatusUnauthorized)
		}

	})

	server.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cookie, err := r.Cookie(cookieName)
		if err != nil {
			log.InfoContext(ctx, "failed to get cookie", "error", err)
			http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
			return
		}

		if cookie.Value != authToken {
			log.InfoContext(ctx, "invalid cookie", "cookie", cookie.Value)
			http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
			return
		}

		type metric struct {
			Id     string   `json:"id"`
			Name   string   `json:"name"`
			Values []string `json:"values"`
			Error  string   `json:"error"`
		}

		type user struct {
			ID       string `db:"id"`
			Username string `db:"username"`
			Email    string `db:"email"`
			Created  string `db:"created"`
		}

		type site struct {
			ID           string `db:"id"`
			Domain       string `db:"domain"`
			Name         string `db:"name"`
			Username     string `db:"username"`
			Email        string `db:"email"`
			CreatedAt    string `db:"created_at"`
			UpdatedAt    string `db:"updated_at"`
			UploadCount  int    `db:"upload_count"`
			WarningCount int    `db:"warning_count"`
			WarningCodes string `db:"warning_codes"`
		}

		type customDomain struct {
			ID                    string `db:"id"`
			Hostname              string `db:"hostname"`
			Status                string `db:"status"`
			SiteName              string `db:"site_name"`
			SiteDomain            string `db:"site_domain"`
			Username              string `db:"username"`
			Email                 string `db:"email"`
			LastVerificationError string `db:"last_verification_error"`
			LastCheckedAt         string `db:"last_checked_at"`
			ActivatedAt           string `db:"activated_at"`
			CreatedAt             string `db:"created_at"`
			UpdatedAt             string `db:"updated_at"`
		}

		type caddyHost struct {
			Host       string
			Kind       string
			Status     string
			SiteName   string
			SiteDomain string
			Username   string
			Email      string
		}

		type dashboardData struct {
			Stats              []metric
			Users              []user
			UsersError         string
			Sites              []site
			SitesError         string
			CustomDomains      []customDomain
			CustomDomainsError string
			CaddyHosts         []caddyHost
			CaddyHostsError    string
			CaddyShorthandId   string
			WipeCounts         WipeCounts
		}

		var data dashboardData
		data.CaddyShorthandId = conf.Caddy.ShorthandId
		data.Stats = make([]metric, 4)

		getMetric := func(id, name string, pos int, query string) {
			data.Stats[pos].Id = id
			data.Stats[pos].Name = name
			var dest = []string{""}
			err := conn.SelectContext(ctx, &dest, query)
			data.Stats[pos].Values = dest
			if err != nil {
				data.Stats[pos].Error = err.Error()
			}
		}

		getMetric("users", "Users", 0, "select count(*) from users")
		getMetric("sites", "Sites", 1, "select count(*) from sites")
		getMetric("custom_domains", "Custom Domains", 2, "select count(*) from custom_domains")
		getMetric("us", "Uploads", 3, "select count(*) from uploadsessions")

		data.WipeCounts.Users, _ = strconv.Atoi(data.Stats[0].Values[0])
		data.WipeCounts.Sites, _ = strconv.Atoi(data.Stats[1].Values[0])
		data.WipeCounts.CustomDomains, _ = strconv.Atoi(data.Stats[2].Values[0])
		data.WipeCounts.UploadSessions, _ = strconv.Atoi(data.Stats[3].Values[0])

		var sessionCount int
		conn.GetContext(ctx, &sessionCount, "SELECT COUNT(*) FROM sessions")
		data.WipeCounts.Sessions = sessionCount

		var subdomainCount int
		conn.GetContext(ctx, &subdomainCount, "SELECT COUNT(*) FROM subdomains")
		data.WipeCounts.Subdomains = subdomainCount

		var fileCount int
		conn.GetContext(ctx, &fileCount, "SELECT COUNT(*) FROM files")
		data.WipeCounts.Files = fileCount

		hosts, err := caddyAPI.GetHostsFromShorthandId(ctx, conf.Caddy.ShorthandId)
		if err != nil {
			data.WipeCounts.CaddyHosts = 0
			data.CaddyHostsError = err.Error()
		} else {
			slices.Sort(hosts)
			data.WipeCounts.CaddyHosts = len(hosts)
			data.CaddyHosts = make([]caddyHost, 0, len(hosts))
			for _, host := range hosts {
				kind := "configured only"
				if strings.HasSuffix(host, "."+conf.PublicBaseDomain) {
					kind = "site subdomain"
				}

				data.CaddyHosts = append(data.CaddyHosts, caddyHost{
					Host: host,
					Kind: kind,
				})
			}

			if len(hosts) > 0 {
				type caddyHostContext struct {
					Host       string `db:"host"`
					Kind       string `db:"kind"`
					Status     string `db:"status"`
					SiteName   string `db:"site_name"`
					SiteDomain string `db:"site_domain"`
					Username   string `db:"username"`
					Email      string `db:"email"`
				}

				query, args, err := sqlx.In(`
					SELECT s.domain as host,
					       'site subdomain' as kind,
					       '' as status,
					       s.name as site_name,
					       s.domain as site_domain,
					       COALESCE(u.username, '') as username,
					       COALESCE(u.email, '') as email
					FROM sites s
					LEFT JOIN users u ON s."user" = u.id
					WHERE s.domain IN (?)
					UNION ALL
					SELECT cd.hostname as host,
					       'custom domain' as kind,
					       cd.status as status,
					       COALESCE(s.name, '') as site_name,
					       COALESCE(s.domain, '') as site_domain,
					       COALESCE(u.username, '') as username,
					       COALESCE(u.email, '') as email
					FROM custom_domains cd
					LEFT JOIN sites s ON cd.site = s.id
					LEFT JOIN users u ON cd."user" = u.id
					WHERE cd.hostname IN (?)
				`, hosts, hosts)
				if err != nil {
					data.CaddyHostsError = err.Error()
				} else {
					query = conn.Rebind(query)
					var contexts []caddyHostContext
					err = conn.SelectContext(ctx, &contexts, query, args...)
					if err != nil {
						data.CaddyHostsError = err.Error()
					} else {
						contextByHost := make(map[string]caddyHostContext, len(contexts))
						for _, hostContext := range contexts {
							contextByHost[hostContext.Host] = hostContext
						}

						for i, host := range data.CaddyHosts {
							hostContext, ok := contextByHost[host.Host]
							if !ok {
								continue
							}

							data.CaddyHosts[i].Kind = hostContext.Kind
							data.CaddyHosts[i].Status = hostContext.Status
							data.CaddyHosts[i].SiteName = hostContext.SiteName
							data.CaddyHosts[i].SiteDomain = hostContext.SiteDomain
							data.CaddyHosts[i].Username = hostContext.Username
							data.CaddyHosts[i].Email = hostContext.Email
						}
					}
				}
			}
		}

		objectsCh := s3Client.ListObjects(ctx, conf.S3.Bucket, minio.ListObjectsOptions{Recursive: true})
		s3Count := 0
		for obj := range objectsCh {
			if obj.Err != nil {
				break
			}
			s3Count++
		}
		data.WipeCounts.S3Objects = s3Count

		// Get last 100 registered users
		err = conn.SelectContext(ctx, &data.Users, `
				SELECT id, username, COALESCE(email, '') as email, created
				FROM users
				WHERE username IS NOT NULL
				ORDER BY created DESC
				LIMIT 100
			`)
		if err != nil {
			data.UsersError = err.Error()
		}

		// Get last 100 created sites with last upload datetime
		err = conn.SelectContext(ctx, &data.Sites, `
				SELECT s.id, s.domain, s.name, COALESCE(u.username, '') as username,
				       COALESCE(u.email, '') as email, s.created_at,
				       COALESCE(MAX(us.updated_at), s.updated_at) as updated_at,
				       COUNT(DISTINCT us.id) as upload_count,
				       COUNT(DISTINCT sw.code) as warning_count,
				       COALESCE(string_agg(DISTINCT sw.code, ', '), '') as warning_codes
				FROM sites s
				LEFT JOIN users u ON s."user" = u.id
				LEFT JOIN uploadsessions us ON s.id = us.site
				LEFT JOIN site_warnings sw ON s.id = sw.site
				GROUP BY s.id, s.domain, s.name, u.username, u.email, s.created_at, s.updated_at
				ORDER BY COALESCE(MAX(us.updated_at), s.updated_at) DESC
				LIMIT 100
			`)
		if err != nil {
			data.SitesError = err.Error()
		}

		// Get last 100 custom domains with site and owner context
		err = conn.SelectContext(ctx, &data.CustomDomains, `
				SELECT cd.id, cd.hostname, cd.status,
				       COALESCE(s.name, '') as site_name,
				       COALESCE(s.domain, '') as site_domain,
				       COALESCE(u.username, '') as username,
				       COALESCE(u.email, '') as email,
				       cd.last_verification_error,
				       COALESCE(cd.last_checked_at::text, '') as last_checked_at,
				       COALESCE(cd.activated_at::text, '') as activated_at,
				       cd.created_at::text as created_at,
				       cd.updated_at::text as updated_at
				FROM custom_domains cd
				LEFT JOIN sites s ON cd.site = s.id
				LEFT JOIN users u ON cd."user" = u.id
				ORDER BY cd.updated_at DESC
				LIMIT 100
			`)
		if err != nil {
			data.CustomDomainsError = err.Error()
		}

		log.InfoContext(ctx, "data", "stats", data)

		err = tmpls.ExecuteTemplate(w, "dashboard.html", data)
		if err != nil {
			log.ErrorContext(ctx, "failed to execute template", "error", err)
			w.WriteHeader(500)
		}
	})

	server.HandleFunc("/api/wipe", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		ctx := r.Context()

		cookie, err := r.Cookie(cookieName)
		if err != nil || cookie.Value != authToken {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
			return
		}

		var req struct {
			Confirmation string `json:"confirmation"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
			return
		}

		if req.Confirmation != "DELETE" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid confirmation"})
			return
		}

		log.InfoContext(ctx, "Starting wipe operation")

		deleted, err := wipeAllData(ctx, conn, s3Client, caddyAPI, conf.S3.Bucket, conf.Caddy.ShorthandId)
		if err != nil {
			log.ErrorContext(ctx, "Wipe failed", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		log.InfoContext(ctx, "Wipe completed successfully", "deleted", deleted)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"deleted": deleted,
			"message": "All data wiped successfully",
		})
	})

	app.RegisterService("server", server)

	if err := app.Run(ctx); err != nil {
		log.ErrorContext(ctx, "applications run failed", "error", err)
	}
}

func wipeAllData(ctx context.Context, conn *sqlx.DB, s3Client *s3.S3, caddyAPI *caddy.CaddyAPI, s3Bucket, caddyShorthand string) (*WipeCounts, error) {
	deleted := &WipeCounts{}

	log.InfoContext(ctx, "Step 1/9: Clearing Caddy subdomain hosts...")
	hosts, err := caddyAPI.GetHostsFromShorthandId(ctx, caddyShorthand)
	if err != nil {
		return nil, fmt.Errorf("failed to get Caddy hosts: %w", err)
	}
	deleted.CaddyHosts = len(hosts)

	err = caddyAPI.ClearSubdomains(ctx, caddyShorthand)
	if err != nil {
		return nil, fmt.Errorf("failed to clear Caddy hosts: %w", err)
	}
	log.InfoContext(ctx, "Cleared Caddy hosts", "count", deleted.CaddyHosts)

	log.InfoContext(ctx, "Step 2/9: Cleaning S3 bucket...")
	objectsCh := s3Client.ListObjects(ctx, s3Bucket, minio.ListObjectsOptions{Recursive: true})
	s3Count := 0
	for obj := range objectsCh {
		if obj.Err != nil {
			return nil, fmt.Errorf("failed to list S3 objects: %w", obj.Err)
		}
		s3Count++
	}
	deleted.S3Objects = s3Count

	err = s3Client.CleanBucket(ctx, s3Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to clean S3 bucket: %w", err)
	}
	log.InfoContext(ctx, "Cleaned S3 bucket", "count", deleted.S3Objects)

	tables := []struct {
		name  string
		field *int
		step  string
	}{
		{"files", &deleted.Files, "3/9"},
		{"uploadsessions", &deleted.UploadSessions, "4/9"},
		{"custom_domains", &deleted.CustomDomains, "5/9"},
		{"sessions", &deleted.Sessions, "6/9"},
		{"sites", &deleted.Sites, "7/9"},
		{"subdomains", &deleted.Subdomains, "8/9"},
		{"users", &deleted.Users, "9/9"},
	}

	for _, table := range tables {
		log.InfoContext(ctx, fmt.Sprintf("Step %s: Deleting %s from database...", table.step, table.name))
		result, err := conn.Exec(fmt.Sprintf("DELETE FROM %s", table.name))
		if err != nil {
			return nil, fmt.Errorf("failed to delete from %s: %w", table.name, err)
		}
		rows, _ := result.RowsAffected()
		*table.field = int(rows)
		log.InfoContext(ctx, "Deleted table", "table", table.name, "count", rows)
	}

	log.InfoContext(ctx, "Wipe completed successfully")
	return deleted, nil
}
