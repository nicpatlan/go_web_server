package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/nicpatlan/go_web_server/internal/database"
)

type ApiConfig struct {
	fileserverHits int
	database       *database.DB
	jwtsecret      string
}

func main() {
	// load environment
	godotenv.Load()
	jwtsecret := os.Getenv("JWT_SECRET")

	// constants
	const filePattern = "/app/*"
	const fileStrip = "/app"
	const fileRoot = "."
	const healthzPattern = "GET /api/healthz"
	const metricPattern = "GET /admin/metrics"
	const resetPattern = "/api/reset"
	const postPattern = "POST /api/posts"
	const getPattern = "GET /api/posts"
	const getOnePattern = "GET /api/posts/{id}"
	const postUserPattern = "POST /api/users"
	const putUserPattern = "PUT /api/users"
	const loginUserPattern = "POST /api/login"
	const refreshPattern = "POST /api/refresh"
	const revokePattern = "POST /api/revoke"
	const port = "8080"
	const dbPath = "database.json"

	// database setup
	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Printf("Error loading database: %s", err)
		return
	}

	// create server mux handler and fileHits counter
	serveMux := http.NewServeMux()
	aCfg := ApiConfig{
		fileserverHits: 0,
		database:       db,
		jwtsecret:      jwtsecret,
	}

	// add handler
	serveMux.Handle(filePattern, aCfg.IncrFileHits(http.StripPrefix(fileStrip, http.FileServer(http.Dir(fileRoot)))))
	serveMux.Handle(healthzPattern, healthzHandler{})
	serveMux.Handle(metricPattern, aCfg.GetHitsHandler())
	serveMux.Handle(resetPattern, aCfg.GetResetHandler())
	serveMux.HandleFunc(postPattern, aCfg.PostHandlerFunc)
	serveMux.HandleFunc(getPattern, aCfg.GetPostsHandlerFunc)
	serveMux.HandleFunc(getOnePattern, aCfg.GetOnePostHandlerFunc)
	serveMux.HandleFunc(postUserPattern, aCfg.CreateUserHandlerFunc)
	serveMux.HandleFunc(loginUserPattern, aCfg.LoginUserHandlerFunc)
	serveMux.HandleFunc(putUserPattern, aCfg.UpdateUserHandlerFunc)
	serveMux.HandleFunc(refreshPattern, aCfg.RefreshTokenHandlerFunc)
	serveMux.HandleFunc(revokePattern, aCfg.RevokeTokenHandlerFunc)

	// create server on localhost port 8080
	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	// begin listening and responding to requests
	log.Printf("Running server from %s on port: %s", fileRoot, port)
	log.Fatal(server.ListenAndServe())
}
