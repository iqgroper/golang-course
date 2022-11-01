package main

import (
	"net/http"
	"os"
	"redditclone/pkg/comments"
	"redditclone/pkg/handlers"
	"redditclone/pkg/middleware"
	"redditclone/pkg/posts"
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
	postsRepo := posts.NewMemoryRepo()
	commentsRepo := comments.NewMemoryRepo()

	userHandler := &handlers.UserHandler{
		UserRepo: userRepo,
		Logger:   logger,
		Sessions: sm,
	}

	postsHandler := &handlers.PostsHandler{
		PostsRepo:    postsRepo,
		CommentsRepo: commentsRepo,
		Logger:       logger,
	}

	r := mux.NewRouter()
	mux := middleware.Auth(sm, r)
	// mux = middleware.AccessLog(logger, mux)
	// mux = middleware.Panic(mux)

	fs := http.FileServer(http.Dir("../../"))
	http.Handle("/static/", fs)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../static/html")
	})

	r.HandleFunc("/api/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/api/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/api/posts/", postsHandler.GetAll).Methods("GET")
	r.HandleFunc("/api/posts", postsHandler.AddPost).Methods("POST")
	r.HandleFunc("/api/post/{post_id}", postsHandler.GetByID).Methods("GET")
	r.HandleFunc("/api/post/{post_id}", postsHandler.AddComment).Methods("POST")
	r.HandleFunc("/api/post/{post_id}/{comment_id}", postsHandler.DeleteComment).Methods("DELETE")
	r.HandleFunc("/api/post/{post_id}/upvote", postsHandler.UpVote).Methods("GET")
	r.HandleFunc("/api/post/{post_id}/downvote", postsHandler.DownVote).Methods("GET")
	r.HandleFunc("/api/post/{post_id}/unvote", postsHandler.UnVote).Methods("GET")
	r.HandleFunc("/api/posts/{category_name}", postsHandler.GetByCategory).Methods("GET")
	r.HandleFunc("/api/user/{user_login}", postsHandler.GetAllByUser).Methods("GET")
	r.HandleFunc("/api/post/{post_id}", postsHandler.DeletePost).Methods("DELETE")
	r.HandleFunc("/api/user/{user_login}", postsHandler.GetAllByUser).Methods("GET")

	// http.Handle("/api/", http.StripPrefix("/api/", r))

	http.Handle("/api/", mux)

	port := ":8084"
	log.Printf("Listening on %s", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
