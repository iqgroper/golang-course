package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"

	"redditclone/pkg/comments"
	"redditclone/pkg/posts"
	"redditclone/pkg/session"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type PostsHandler struct {
	PostsRepo    posts.PostRepo
	CommentsRepo comments.CommentRepo
	Logger       *logrus.Entry
}

const postTemplate = `
{"score":{{.Score}},
"views":{{.Views}},
"type":"{{.Type}}",
"title":"{{.Title}}",
"author":{"username":"{{.Author.Username}}","id":"{{.Author.ID}}"},
"category":"{{.Category}}",
{{if (eq .Type "text")}}
"text":"{{.Text}}",
{{else}}
"url":"{{.URL}}",
{{end}}
"votes":[
	{{$first := 1}}
	{{range .VotesList}}
	{{if $first}}{{$first = 0}}{{else}},{{end}}
		{"user":"{{.User}}","vote":{{.Vote}}}
	{{end}}
],
"comments":[
	{{$first := 1}}
	{{range .Comments}}
	{{if $first}}{{$first = 0}}{{else}},{{end}}
	{	
		"author":{"username":"{{.Author.Login}}", "id":"{{.Author.ID}}"},
		"body":"{{.Body}}",
		"created":"{{.Created}}",
		"id":"{{.ID}}"
	}
	{{end}}
],
"created":"{{.CreatedDTTM}}",
"upvotePercentage":{{.UpvotePercentage}},
"id":"{{.ID}}"}
`

func (h *PostsHandler) AddPost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error reading body")
		return
	}

	newPost := &posts.NewPost{}
	errorUnmarshal := json.Unmarshal(body, newPost)
	if errorUnmarshal != nil {
		h.Logger.Println("error unmarshling new post:", errorUnmarshal.Error())
		http.Error(w, fmt.Sprintf(`error unmarshling new post: %s`, errorUnmarshal.Error()), http.StatusInternalServerError)
		return
	}

	sess, errSess := session.SessionFromContext(r.Context())
	if errSess != nil {
		h.Logger.Println("error getting session in AddPost:", errSess.Error())
		http.Error(w, fmt.Sprintf(`error getting session in AddPost: %s`, errSess.Error()), http.StatusInternalServerError)
		return
	}
	newPost.Author = *sess.User

	addedPost, errAdd := h.PostsRepo.Add(newPost)
	if errAdd != nil {
		h.Logger.Println("error adding post:", errAdd.Error())
		http.Error(w, fmt.Sprintf(`error adding post: %s`, errAdd.Error()), http.StatusInternalServerError)
		return
	}

	tmpl := template.New("")
	tmpl, errParse := tmpl.Parse(postTemplate)
	if errParse != nil {
		fmt.Println("Error parsing AddPost method", errParse.Error())
	}
	var resp bytes.Buffer

	errExecution := tmpl.Execute(&resp, addedPost)
	if errExecution != nil {
		fmt.Println("Error executing template in AddPost:", errExecution.Error())
		return
	}

	w.Write(resp.Bytes())
}

func getIDFromString(id string) uint {
	u64, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	return uint(u64)
}

func (h *PostsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	postIDStr, ok := vars["post_id"]
	if !ok {
		h.Logger.Println("no post_id in query in GetByID")
		http.Error(w, "no post_id in query in GetByID", http.StatusBadRequest)
		return
	}

	postID := getIDFromString(postIDStr)

	foundPost, errFind := h.PostsRepo.GetByID(postID)
	if errFind != nil {
		h.Logger.Println("no such post in repo.getbyid:", errFind.Error())
		http.Error(w, fmt.Sprintln("no such post in repo.getbyid:", errFind.Error()), http.StatusBadRequest)
		return
	}

	foundPost.Views += 1

	responseBody := SendPost(foundPost, "GetByID")

	w.Write([]byte(responseBody))
}

func SendPost(post *posts.Post, function string) string {
	tmpl := template.New("")
	tmpl, errParse := tmpl.Parse(postTemplate)
	if errParse != nil {
		fmt.Printf("Error parsing in %s: %s", function, errParse.Error())
	}
	var resp bytes.Buffer

	errExecution := tmpl.Execute(&resp, post)
	if errExecution != nil {
		fmt.Printf("Error executing template in %s: %s", function, errExecution.Error())
		return ""
	}
	return resp.String()
}

func (h *PostsHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	postIDStr, ok := vars["post_id"]
	if !ok {
		h.Logger.Println("no post_id in query in AddComment")
		http.Error(w, "no post_id in query in AddComment", http.StatusBadRequest)
		return
	}

	postID := getIDFromString(postIDStr)

	sess, errSess := session.SessionFromContext(r.Context())
	if errSess != nil {
		h.Logger.Println("error getting session in AddComment:", errSess.Error())
		http.Error(w, fmt.Sprintf(`error getting session in AddComment: %s`, errSess.Error()), http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	bodyRaw, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error reading body in AddComment")
		return
	}

	body := &struct{ Comment string }{}

	errorUnmarshal := json.Unmarshal(bodyRaw, body)
	if errorUnmarshal != nil {
		fmt.Println("error unmarshling in AddComment:", errorUnmarshal.Error())
		return
	}

	if len(body.Comment) == 0 {

		http.Error(w, `{"errors":[{"location":"body","param":"comment","msg":"is required"}]}`, http.StatusUnprocessableEntity)
		return
	}

	newComment, errAdd := h.CommentsRepo.Add(postID, body.Comment, sess.User)
	if errAdd != nil {
		h.Logger.Println("error adding in AddComment:", errAdd.Error())
		http.Error(w, fmt.Sprintf(`error addinc comment in AddComment: %s`, errAdd.Error()), http.StatusInternalServerError)
		return
	}

	foundPost, errFind := h.PostsRepo.GetByID(postID)
	if errFind != nil {
		h.Logger.Println("cant find post in AddComment:", errFind.Error())
		http.Error(w, fmt.Sprintln("cant find post in AddComment:", errFind.Error()), http.StatusBadRequest)
		return
	}

	foundPost.Comments = append(foundPost.Comments, newComment)

	responseBody := SendPost(foundPost, "AddComment")

	w.Write([]byte(responseBody))
}

func (h *PostsHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	postIDStr, ok := vars["post_id"]
	if !ok {
		h.Logger.Println("no post_id in query in DeleteComment")
		http.Error(w, "no post_id in query in DeleteComment", http.StatusBadRequest)
		return
	}
	commentIDStr, ok := vars["comment_id"]
	if !ok {
		h.Logger.Println("no comment_id in query in DeleteComment")
		http.Error(w, "no comment_id in query in DeleteComment", http.StatusBadRequest)
		return
	}

	postID := getIDFromString(postIDStr)
	commentID := getIDFromString(commentIDStr)

	ok, err := h.CommentsRepo.Delete(postID, commentID)
	if !ok || err != nil {
		h.Logger.Println("cant delete a comment in DeleteComment", err.Error())
		http.Error(w, "cant delete a comment in DeleteComment", http.StatusBadRequest)
		return
	}

	foundPost, errFind := h.PostsRepo.GetByID(postID)
	if errFind != nil {
		h.Logger.Println("cant find post in DeleteComment:", errFind.Error())
		http.Error(w, fmt.Sprintln("cant find post in DeleteComment:", errFind.Error()), http.StatusBadRequest)
		return
	}

	newComments, errGetAll := h.CommentsRepo.GetAll(foundPost.ID)
	if errGetAll != nil {
		h.Logger.Println("cant GetAll in DeleteComment:", errGetAll.Error())
		http.Error(w, fmt.Sprintln("cant GetAll in DeleteComment:", errGetAll.Error()), http.StatusInternalServerError)
		return
	}
	foundPost.Comments = newComments

	responseBody := SendPost(foundPost, "DeleteComment")

	w.Write([]byte(responseBody))
}

func (h *PostsHandler) UpVote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	postIDStr, ok := vars["post_id"]
	if !ok {
		h.Logger.Println("no post_id in query in UpVote")
		http.Error(w, "no post_id in query in UpVote", http.StatusBadRequest)
		return
	}

	postID := getIDFromString(postIDStr)

	sess, errSess := session.SessionFromContext(r.Context())
	if errSess != nil {
		h.Logger.Println("error getting session in UpVote:", errSess.Error())
		http.Error(w, fmt.Sprintf(`error getting session in UpVote: %s`, errSess.Error()), http.StatusInternalServerError)
		return
	}

	foundPost, err := h.PostsRepo.UpVote(postID, sess.User.Login)
	if err == posts.ErrNoCanDo {
		h.Logger.Println("no such post in UpVote:", err.Error())
		http.Redirect(w, r, fmt.Sprintf("/api/post/%s/unvote", postIDStr), http.StatusFound)
		// http.Error(w, fmt.Sprintf(`{"message":"%s"}`, err.Error()), http.StatusMethodNotAllowed)
		return
	}
	if err != nil {
		h.Logger.Println(err.Error())
		http.Error(w, "UpVote:"+err.Error(), http.StatusBadRequest)
		return
	}

	responseBody := SendPost(foundPost, "UpVote")

	w.Write([]byte(responseBody))
}

func (h *PostsHandler) DownVote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	postIDStr, ok := vars["post_id"]
	if !ok {
		h.Logger.Println("no post_id in query in DownVote")
		http.Error(w, "no post_id in query in DownVote", http.StatusBadRequest)
		return
	}

	postID := getIDFromString(postIDStr)

	sess, errSess := session.SessionFromContext(r.Context())
	if errSess != nil {
		h.Logger.Println("error getting session in DownVote:", errSess.Error())
		http.Error(w, fmt.Sprintf(`error getting session in DownVote: %s`, errSess.Error()), http.StatusInternalServerError)
		return
	}

	foundPost, err := h.PostsRepo.DownVote(postID, sess.User.Login)
	if err == posts.ErrNoCanDo {
		h.Logger.Println("in DownVote:", err.Error())
		http.Redirect(w, r, fmt.Sprintf("/api/post/%s/unvote", postIDStr), http.StatusFound)
		// http.Error(w, fmt.Sprintf(`{"message":"%s"}`, err.Error()), http.StatusMethodNotAllowed)
		return
	}
	if err != nil {
		h.Logger.Println(err.Error())
		http.Error(w, "DownVote:"+err.Error(), http.StatusBadRequest)
		return
	}

	responseBody := SendPost(foundPost, "DownVote")

	w.Write([]byte(responseBody))
}

func (h *PostsHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	postIDStr, ok := vars["post_id"]
	if !ok {
		h.Logger.Println("no post_id in query in DeletePost")
		http.Error(w, "no post_id in query in DeletePost", http.StatusBadRequest)
		return
	}

	postID := getIDFromString(postIDStr)

	_, errDel := h.PostsRepo.Delete(postID)
	if errDel != nil {
		h.Logger.Println("in DeletePost:", errDel.Error())
		http.Error(w, "in DeletePost: "+errDel.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.CommentsRepo.DeleteAllByPost(postID)
	if err != nil {
		h.Logger.Println("cant delete all comments in DeleteCommentByPost", err.Error())
		http.Error(w, "cant delete all comments in DeleteCommentByPost", http.StatusBadRequest)
		return
	}

	resp, errMar := json.Marshal(map[string]string{"message": "success"})
	if errMar != nil {
		h.Logger.Println("cant marshal in DeleteCommentByPost", err.Error())
		http.Error(w, "cant marshal in DeleteCommentByPost "+errMar.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}

const postsTemplate = `[{{$firstOuter := 1}}{{range .}}{{if $firstOuter}}{{$firstOuter = 0}}{{else}},{{end}}
	{"score":{{.Score}},
	"views":{{.Views}},
	"type":"{{.Type}}",
	"title":"{{.Title}}",
	"author":{"username":"{{.Author.Username}}","id":"{{.Author.ID}}"},
	"category":"{{.Category}}",
	{{if (eq .Type "text")}}
	"text":"{{.Text}}",
	{{else}}
	"url":"{{.URL}}",
	{{end}}
	"votes":[
		{{$first := 1}}
		{{range .VotesList}}
		{{if $first}}{{$first = 0}}{{else}},{{end}}
			{"user":"{{.User}}","vote":{{.Vote}}}
		{{end}}
	],
	"comments":[
		{{$first := 1}}
		{{range .Comments}}
		{{if $first}}{{$first = 0}}{{else}},{{end}}
		{	
			"author":{"username":"{{.Author.Login}}", "id":"{{.Author.ID}}"},
			"body":"{{.Body}}",
			"created":"{{.Created}}",
			"id":"{{.ID}}"
		}
		{{end}}
	],
	"created":"{{.CreatedDTTM}}",
	"upvotePercentage":{{.UpvotePercentage}},
	"id":"{{.ID}}"}
{{end}}]
`

func (h *PostsHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	postList, errGetting := h.PostsRepo.GetAll()
	if errGetting != nil {
		h.Logger.Println(errGetting)
		w.Write([]byte("[]"))
		return
	}

	tmpl := template.New("")
	tmpl, errParse := tmpl.Parse(postsTemplate)
	if errParse != nil {
		fmt.Printf("Error parsing: %s", errParse.Error())
	}
	var resp bytes.Buffer

	errExecution := tmpl.Execute(&resp, postList)
	if errExecution != nil {
		fmt.Printf("Error executing template: %s", errExecution.Error())
	}

	w.Write(resp.Bytes())
}

func (h *PostsHandler) GetByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryName, ok := vars["category_name"]
	if !ok {
		h.Logger.Println("no category_name in query")
		http.Error(w, "no category_name in query", http.StatusBadRequest)
		return
	}

	postList, errGetting := h.PostsRepo.GetAllByCategory(categoryName)
	if errGetting != nil {
		h.Logger.Println(errGetting)
		w.Write([]byte("[]"))
		return
	}

	tmpl := template.New("")
	tmpl, errParse := tmpl.Parse(postsTemplate)
	if errParse != nil {
		fmt.Printf("Error parsing: %s", errParse.Error())
	}
	var resp bytes.Buffer

	errExecution := tmpl.Execute(&resp, postList)
	if errExecution != nil {
		fmt.Printf("Error executing template: %s", errExecution.Error())
	}

	w.Write(resp.Bytes())
}

func (h *PostsHandler) GetAllByUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userLogin, ok := vars["user_login"]
	if !ok {
		h.Logger.Println("no user_login in query")
		http.Error(w, "no user_login in query", http.StatusBadRequest)
		return
	}

	postList, errGetting := h.PostsRepo.GetByUser(userLogin)
	if errGetting != nil {
		h.Logger.Println(errGetting)
		w.Write([]byte("[]"))
		return
	}

	tmpl := template.New("")
	tmpl, errParse := tmpl.Parse(postsTemplate)
	if errParse != nil {
		fmt.Printf("Error parsing: %s", errParse.Error())
	}
	var resp bytes.Buffer

	errExecution := tmpl.Execute(&resp, postList)
	if errExecution != nil {
		fmt.Printf("Error executing template: %s", errExecution.Error())
	}

	w.Write(resp.Bytes())
}

func (h *PostsHandler) UnVote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	postIDStr, ok := vars["post_id"]
	if !ok {
		h.Logger.Println("no post_id in query in UnVote")
		http.Error(w, "no post_id in query in UnVote", http.StatusBadRequest)
		return
	}

	postID := getIDFromString(postIDStr)

	sess, errSess := session.SessionFromContext(r.Context())
	if errSess != nil {
		h.Logger.Println("error getting session in UnVote:", errSess.Error())
		http.Error(w, fmt.Sprintf(`error getting session in UnVote: %s`, errSess.Error()), http.StatusInternalServerError)
		return
	}

	foundPost, errFind := h.PostsRepo.UnVote(postID, sess.User.Login)
	if errFind != nil {
		h.Logger.Println("no such post in UnVote:", errFind.Error())
		http.Error(w, fmt.Sprintln("no such post in UnVote:", errFind.Error()), http.StatusMethodNotAllowed)
		return
	}

	responseBody := SendPost(foundPost, "UnVote")

	w.Write([]byte(responseBody))
}
