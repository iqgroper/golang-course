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

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {

	sm := session.NewSessionsManager()
	logger := log.WithFields(log.Fields{})

	userRepo := user.NewMysqlRepo()
	postsRepo := posts.NewPostsMongoRepository()
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

	fs := http.FileServer(http.Dir("../../"))
	http.Handle("/static/", fs)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../static/html")
	})

	authRouter := mux.NewRouter()
	noAuthRouter := mux.NewRouter()

	noAuthRouter.HandleFunc("/api/register", userHandler.Register).Methods("POST")
	noAuthRouter.HandleFunc("/api/login", userHandler.Login).Methods("POST")
	noAuthRouter.HandleFunc("/api/posts/", postsHandler.GetAll).Methods("GET")
	noAuthRouter.HandleFunc("/api/post/{post_id}", postsHandler.GetByID).Methods("GET")
	noAuthRouter.HandleFunc("/api/posts/{category_name}", postsHandler.GetByCategory).Methods("GET")
	noAuthRouter.HandleFunc("/api/user/{user_login}", postsHandler.GetAllByUser).Methods("GET")

	authRouter.HandleFunc("/api/posts", postsHandler.AddPost).Methods("POST")
	authRouter.HandleFunc("/api/post/{post_id}", postsHandler.AddComment).Methods("POST")
	authRouter.HandleFunc("/api/post/{post_id}/{comment_id}", postsHandler.DeleteComment).Methods("DELETE")
	authRouter.HandleFunc("/api/post/{post_id}/upvote", postsHandler.UpVote).Methods("GET")
	authRouter.HandleFunc("/api/post/{post_id}/downvote", postsHandler.DownVote).Methods("GET")
	authRouter.HandleFunc("/api/post/{post_id}/unvote", postsHandler.UnVote).Methods("GET")
	authRouter.HandleFunc("/api/post/{post_id}", postsHandler.DeletePost).Methods("DELETE")

	authMux := middleware.Auth(sm, authRouter)

	noAuthRouter.PathPrefix("/api/").Handler(authMux)

	siteMux := middleware.AccessLog(logger, noAuthRouter)
	// siteMux = middleware.Panic(siteMux)

	http.Handle("/api/", siteMux)

	port := ":8080"
	log.Printf("Listening on %s", port)
	errServer := http.ListenAndServe(port, nil)
	if errServer != nil {
		log.Fatal(errServer)
	}
}
