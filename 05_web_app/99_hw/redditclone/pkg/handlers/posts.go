package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

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
