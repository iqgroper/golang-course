package main

import (
	"net/http"
	"os"
	"redditclone/pkg/handlers"
	"redditclone/pkg/session"
	"redditclone/pkg/user"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {

	sm := session.NewSessionsManager()
	logger := log.WithFields(log.Fields{})

	userRepo := user.NewMemoryRepo()
	// postsRepo := posts.NewMemoryRepo()
	// commentsRepo := comments.NewMemoryRepo()

	userHandler := &handlers.UserHandler{
		UserRepo: userRepo,
		Logger:   logger,
		Sessions: sm,
	}

	// postsHandler := &handlers.ItemsHandler{
	// 	Tmpl:      templates,
	// 	Logger:    logger,
	// 	ItemsRepo: itemsRepo,
	// }

	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("../../"))
	http.Handle("/static/", fs)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../static/html")
	})

	r.HandleFunc("/api/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/api/login", userHandler.Login).Methods("POST")

	// http.Handle("/api/", http.StripPrefix("/api/", r))
	http.Handle("/api/", r)

	log.Print("Listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
