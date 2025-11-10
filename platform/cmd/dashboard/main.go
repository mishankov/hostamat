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

	"github.com/mishankov/platforma/application"
	"github.com/mishankov/platforma/database"
	"github.com/mishankov/platforma/httpserver"
	"github.com/mishankov/platforma/log"
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
	server := httpserver.New("4444", 3*time.Second)
	db, err := database.New(os.Getenv("DB_CONN"))
	if err != nil {
		log.ErrorContext(ctx, "failed to create database", "error", err)
	}
	defer db.Connection().Close()

	tmpls, err := template.ParseFS(Templates(), "login.html", "dashboard.html")
	if err != nil {
		log.ErrorContext(ctx, "failed to parse templates", "error", err)
	}

	authToken := os.Getenv("TOKEN")

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

	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

		err = tmpls.ExecuteTemplate(w, "dashboard.html", nil)
		if err != nil {
			log.ErrorContext(ctx, "failed to execute template", "error", err)
		}
	})

	app.RegisterService("server", server)
	app.RegisterDatabase("default", db)

	if err := app.Run(ctx); err != nil {
		log.ErrorContext(ctx, "applications run failed", "error", err)
	}
}
