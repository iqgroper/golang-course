package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"redditclone/pkg/user"

	"redditclone/pkg/comments"
	"redditclone/pkg/posts"

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

type NewPost struct {
	Score    int
	Views    uint
	Type     string
	Title    string
	Author   user.NewUser
	Category string
	Text     string
	Votes    []struct {
		user string
		vote int
	}
	Comments         []interface{}
	Created          string
	UpvotePercentage int
	ID               uint
}

func (h *PostsHandler) AddPost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error reading body")
		return
	}

	newPost := &NewPost{}
	errorUnmarshal := json.Unmarshal(body, newPost)
	if errorUnmarshal != nil {
		fmt.Println("error unmarshling:", errorUnmarshal.Error())
		return
	}
	fmt.Println(newPost)
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
