package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"gotth/app"
	"gotth/app/assets"
	"gotth/app/docs"
	_ "github.com/a-h/templ"
)

func main() {
	mux := http.NewServeMux()

	// Setup virtual filesystem for embedded public assets
	publicFS, err := fs.Sub(assets.Public, "public")
	if err != nil {
		log.Fatal(err)
	}
	fileServer := http.FileServer(http.FS(publicFS))

	// Serve built tailwind css output directly from memory
	mux.Handle("GET /globals.css", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		w.Write(assets.CSS)
	}))

	// Serve public assets automatically from memory
	mux.HandleFunc("GET /{file}", func(w http.ResponseWriter, r *http.Request) {
		fileName := r.PathValue("file")
		if file, err := publicFS.Open(fileName); err == nil {
			file.Close()
			fileServer.ServeHTTP(w, r)
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}
	addr := ":" + port
	fmt.Printf("Server Component Factory™ running at http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
