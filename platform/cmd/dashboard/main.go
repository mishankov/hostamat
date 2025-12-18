package main

import (
	"context"
	"embed"
	"encoding/json"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"time"

	"github.com/platforma-dev/platforma/application"
	"github.com/platforma-dev/platforma/database"
	"github.com/platforma-dev/platforma/httpserver"
	"github.com/platforma-dev/platforma/log"
)

const cookieName = "dashsession"

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

	db, err := database.New(os.Getenv("DB_CONN"))
	if err != nil {
		log.ErrorContext(ctx, "failed to create database", "error", err)
	}
	conn := db.Connection()
	defer conn.Close()

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
			Created  string `db:"created"`
		}

		type site struct {
			ID          string `db:"id"`
			Domain      string `db:"domain"`
			Name        string `db:"name"`
			Username    string `db:"username"`
			CreatedAt   string `db:"created_at"`
			UpdatedAt   string `db:"updated_at"`
			UploadCount int    `db:"upload_count"`
		}

		type dashboardData struct {
			Stats      []metric
			Users      []user
			UsersError string
			Sites      []site
			SitesError string
		}

		var data dashboardData
		data.Stats = make([]metric, 3)

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
		getMetric("us", "Uploads", 2, "select count(*) from uploadsessions")

		// Get last 10 registered users
		err = conn.SelectContext(ctx, &data.Users, `
				SELECT id, username, created
				FROM users
				WHERE username IS NOT NULL
				ORDER BY created DESC
				LIMIT 10
			`)
		if err != nil {
			data.UsersError = err.Error()
		}

		// Get last 10 created sites with last upload datetime
		err = conn.SelectContext(ctx, &data.Sites, `
				SELECT s.id, s.domain, s.name, u.username, s.created_at,
				       COALESCE(MAX(us.updated_at), s.updated_at) as updated_at,
				       COUNT(us.id) as upload_count
				FROM sites s
				LEFT JOIN users u ON s."user" = u.id
				LEFT JOIN uploadsessions us ON s.id = us.site
				GROUP BY s.id, s.domain, s.name, u.username, s.created_at, s.updated_at
				ORDER BY COALESCE(MAX(us.updated_at), s.updated_at) DESC
				LIMIT 10
			`)
		if err != nil {
			data.SitesError = err.Error()
		}

		log.InfoContext(ctx, "data", "stats", data)

		err = tmpls.ExecuteTemplate(w, "dashboard.html", data)
		if err != nil {
			log.ErrorContext(ctx, "failed to execute template", "error", err)
			w.WriteHeader(500)
		}
	})

	app.RegisterService("server", server)

	if err := app.Run(ctx); err != nil {
		log.ErrorContext(ctx, "applications run failed", "error", err)
	}
}
