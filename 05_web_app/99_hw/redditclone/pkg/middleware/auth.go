package middleware

import (
	"fmt"
	"net/http"

	"redditclone/pkg/session"
)

var (
	noAuthUrls = map[string]struct{}{
		"/api/register": {},
		"/api/login":    {},
		"/api/posts/":   {},
	}
	noSessUrls = map[string]struct{}{
		"/": {},
	}
)

func Auth(sm *session.SessionsManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("auth middleware:", r.URL.Path)
		if _, ok := noAuthUrls[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}
		sess, err := sm.Check(r)
		_, canbeWithouthSess := noSessUrls[r.URL.Path]
		if err != nil && !canbeWithouthSess {
			fmt.Println("no auth:", err.Error())
			http.Redirect(w, r, "/api/posts/", http.StatusFound)
			return
		}
		ctx := session.ContextWithSession(r.Context(), sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
