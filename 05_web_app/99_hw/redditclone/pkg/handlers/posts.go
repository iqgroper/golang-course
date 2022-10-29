package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"redditclone/pkg/comments"
	"redditclone/pkg/posts"
	"redditclone/pkg/session"

	"github.com/sirupsen/logrus"
)

type PostsHandler struct {
	PostsRepo    posts.PostRepo
	CommentsRepo comments.CommentRepo
	Logger       *logrus.Entry
}

type ReturnStruct struct {
	index   int
	element posts.Post
}

const bodyTemplate = `[{{range .}}
"{{.index}}":{{.elem}}
{{end}}]`

func (h *PostsHandler) GetAll(w http.ResponseWriter, r *http.Request) { //elem behavior
	elems, err := h.PostsRepo.GetAll()
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}

	respSlice := make([]ReturnStruct, 0, 100)
	for idx, elem := range elems {
		resp := ReturnStruct{
			index:   idx,
			element: *elem,
		}
		respSlice = append(respSlice, resp)
	}

	tmpl := template.New("")
	tmpl, errParse := tmpl.Parse(bodyTemplate)
	if errParse != nil {
		fmt.Println("Error parsing New method", errParse.Error())
	}
	var resp bytes.Buffer

	errExecution := tmpl.Execute(&resp, respSlice)
	if errExecution != nil {
		fmt.Println("Error executing template in GetAll posts:", errExecution.Error())
		return
	}

	fmt.Println(resp.String())

}

const postTemplate = `
{"score":{{.Score}},
"views":{{.Views}},
"type":"{{.Text}}",
"title":"{{.Title}}",
"author":{"username":"{{.Author.Username}}","id":"{{.Author.ID}}"},
"category":"funny",
"text":"asdfasdfasdfa",
"votes":[{"user":"{{.Author.Username}}","vote":1}],
"comments":[],
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
		h.Logger.Println("error getting session in Add:", errSess.Error())
		http.Error(w, fmt.Sprintf(`error getting session in Add: %s`, errSess.Error()), http.StatusInternalServerError)
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

func (h *PostsHandler) GetByCategory(w http.ResponseWriter, r *http.Request) {

}

func (h *PostsHandler) GetByID(w http.ResponseWriter, r *http.Request) {

}

func (h *PostsHandler) AddComment(w http.ResponseWriter, r *http.Request) {

}

func (h *PostsHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {

}

func (h *PostsHandler) UpVote(w http.ResponseWriter, r *http.Request) {

}

func (h *PostsHandler) DownVote(w http.ResponseWriter, r *http.Request) {

}

func (h *PostsHandler) DeletePost(w http.ResponseWriter, r *http.Request) {

}

func (h *PostsHandler) GetAllByUser(w http.ResponseWriter, r *http.Request) {

}
