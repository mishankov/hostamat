package formcollection

import (
	"encoding/json"
	"html/template"
	"net/http"
)

type handlerError struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

func respondJSONError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(handlerError{
		ErrorCode:    code,
		ErrorMessage: message,
	})
}

func respondSubmissionPage(w http.ResponseWriter, status int, title, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(`<!doctype html>
<html lang="ru">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>` + htmlEscape(title) + `</title>
<style>
body{margin:0;font-family:system-ui,-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;background:#f8fafc;color:#111827}
main{min-height:100vh;display:grid;place-items:center;padding:24px}
section{max-width:420px;text-align:center}
h1{font-size:24px;line-height:1.25;margin:0 0 12px}
p{font-size:16px;line-height:1.5;margin:0;color:#4b5563}
</style>
</head>
<body><main><section><h1>` + htmlEscape(title) + `</h1><p>` + htmlEscape(message) + `</p></section></main></body>
</html>`))
}

func htmlEscape(value string) string {
	return template.HTMLEscapeString(value)
}
