package middleware

import (
	"fmt"
	"net/http"

	"redditclone/pkg/session"
)

var (
	noSessUrls = map[string]struct{}{
		"/": {},
	}
)

func Auth(sm *session.SessionsRedisManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, err := sm.Check(r)
		_, canbeWithouthSess := noSessUrls[r.URL.Path]
		if err != nil && !canbeWithouthSess {
			fmt.Println("no auth:", err.Error(), r.Method, r.URL.Path)
			http.Redirect(w, r, "/api/posts/", http.StatusFound)
			return
		}
		ctx := session.ContextWithSession(r.Context(), sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
