package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"gotth/app"
	"gotth/app/assets"
	"gotth/app/docs"
	_ "github.com/a-h/templ"
)

func main() {
	mux := http.NewServeMux()

	// Serve built tailwind css output directly from memory
	mux.Handle("GET /globals.css", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		w.Write(assets.CSS)
	}))

	// Serve public assets automatically from public/ folder at root /
	mux.HandleFunc("GET /{file}", func(w http.ResponseWriter, r *http.Request) {
		filePath := filepath.Join("public", r.PathValue("file"))
		if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
			http.ServeFile(w, r, filePath)
			return
		}
		http.NotFound(w, r)
	})

	// Route handlers matching Next page routing
	mux.HandleFunc("GET /{$}", app.PageHandler)
	mux.HandleFunc("POST /chart", app.ChartHandler)
	mux.HandleFunc("GET /wisdom", app.WisdomHandler)

	// Docs Page Routes (supporting boosted slug mapping)
	mux.HandleFunc("GET /docs", docs.PageHandler)
	mux.HandleFunc("GET /docs/{slug}", docs.PageHandler)

	port := ":3000"
	fmt.Printf("Server Component Factory™ running at http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
