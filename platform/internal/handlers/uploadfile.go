package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"hostamat/internal/adapters/caddy"
	"hostamat/internal/adapters/config"
	"hostamat/internal/adapters/timeweb"
	"log"
	"net/http"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type logger interface {
	DebugContext(ctx context.Context, msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}

type UploadFileHandler struct {
	logger     logger
	timewebAPI *timeweb.TimewebAPI
	caddyAPI   *caddy.CaddyAPI
	baseDomain string
}

func NewUploadFileHandler(logger logger, timewebAPI *timeweb.TimewebAPI, caddyAPI *caddy.CaddyAPI, baseDomain string) *UploadFileHandler {
	return &UploadFileHandler{logger: logger, timewebAPI: timewebAPI, caddyAPI: caddyAPI, baseDomain: baseDomain}
}

func (ufh *UploadFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseMultipartForm(10 << 20) // 10 MB

	file, handler, err := r.FormFile("file")
	if err != nil {
		ufh.logger.ErrorContext(ctx, "error retrieving file", "error", err)
		http.Error(w, "Error retrieving file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fmt.Fprintf(w, "Received file: %s\n", handler.Filename)
	fmt.Fprintf(w, "File size: %d\n", handler.Size)

	subdomain := randomBase64String(5)

	err = ufh.timewebAPI.AddSubdomain(ctx, ufh.baseDomain, subdomain)
	if err != nil {
		ufh.logger.ErrorContext(ctx, "error in timeweb add subdomain", "error", err)
	}

	err = ufh.caddyAPI.AddSubdomain(ctx, subdomain)
	if err != nil {
		ufh.logger.ErrorContext(ctx, "error in caddy add subdomain", "error", err)
	}

	fmt.Fprintf(w, "File saved successfully.")

	c := config.New()

	log.Println("Minio config", c.S3)

	minioClient, err := minio.New(c.S3.BaseURL, &minio.Options{
		Creds:  credentials.NewStaticV4(c.S3.Login, c.S3.Password, ""),
		Secure: false,
	})

	if err != nil {
		log.Fatal("Error creating minio client", err)
	}

	err = minioClient.MakeBucket(ctx, c.S3.Bucket, minio.MakeBucketOptions{})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, c.S3.Bucket)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", c.S3.Bucket)
		} else {
			log.Fatalln("Error creating bucket", err)
		}
	} else {
		log.Printf("Successfully created %s\n", c.S3.Bucket)
	}

	info, err := minioClient.PutObject(ctx, c.S3.Bucket, subdomain+"/"+handler.Filename, file, handler.Size, minio.PutObjectOptions{})
	if err != nil {
		log.Fatal("Error saving file to minio", err)
	}

	ufh.logger.DebugContext(ctx, "Upload file info", "info", info)
}

func randomBase64String(length int) string {
	bytes := make([]byte, length+1)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal(err)
	}

	result := base64.URLEncoding.EncodeToString(bytes)

	return result[:len(result)-1]
}
