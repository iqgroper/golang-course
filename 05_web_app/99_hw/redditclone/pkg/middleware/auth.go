package middleware

import (
	"fmt"
	"net/http"
	"regexp"

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

func noAuthMethodsCheck(r *http.Request) bool {

	path := r.URL.Path
	if _, ok := noAuthUrls[path]; ok {
		return true
	}

	matchedPost, err := regexp.MatchString(`/api/post/[0-9]+`, path)
	if err != nil {
		fmt.Println("regexp err", err.Error())
	}

	matchedPosts, errP := regexp.MatchString(`/api/posts/[a-z]+`, path)
	if errP != nil {
		fmt.Println("regexp err", errP.Error())
	}

	matchedUser, errU := regexp.MatchString(`/api/user/[0-9a-z]+`, path)
	if errU != nil {
		fmt.Println("regexp err", errU.Error())
	}

	matchedVote, errV := regexp.MatchString(`/api/post/[0-9]+/upvote|/downvote`, path)
	if errV != nil {
		fmt.Println("regexp errV", errV.Error())
	}

	result := (matchedPost || matchedPosts || matchedUser) && !matchedVote && r.Method != "POST" && r.Method != "DELETE"
	fmt.Println("regexp check:", result)
	return result
}

func Auth(sm *session.SessionsManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("auth middle:", r.URL.Path, r.Method)
		if noAuthMethodsCheck(r) {
			next.ServeHTTP(w, r)
			return
		}
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
