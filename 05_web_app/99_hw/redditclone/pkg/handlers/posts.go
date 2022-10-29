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

const newPostTemplate = `
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
	tmpl, errParse := tmpl.Parse(newPostTemplate)
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

	tmpl := template.New("")
	tmpl, errParse := tmpl.Parse(postTemplate)
	if errParse != nil {
		fmt.Println("Error parsing in getPostByID method", errParse.Error())
	}
	var resp bytes.Buffer

	errExecution := tmpl.Execute(&resp, foundPost)
	if errExecution != nil {
		fmt.Println("Error executing template in getPostByID:", errExecution.Error())
		return
	}

	w.Write(resp.Bytes())
}

func (h *PostsHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	postIDStr, ok := vars["post_id"]
	if !ok {
		h.Logger.Println("no post_id in query in GetByID")
		http.Error(w, "no post_id in query in GetByID", http.StatusBadRequest)
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

	tmpl := template.New("")
	tmpl, errParse := tmpl.Parse(postTemplate)
	if errParse != nil {
		fmt.Println("Error parsing", errParse.Error())
	}
	var resp bytes.Buffer

	errExecution := tmpl.Execute(&resp, foundPost)
	if errExecution != nil {
		fmt.Println("Error executing template:", errExecution.Error())
		return
	}

	w.Write(resp.Bytes())
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

	tmpl := template.New("")
	tmpl, errParse := tmpl.Parse(postTemplate)
	if errParse != nil {
		fmt.Println("Error parsing", errParse.Error())
	}
	var resp bytes.Buffer

	errExecution := tmpl.Execute(&resp, foundPost)
	if errExecution != nil {
		fmt.Println("Error executing template:", errExecution.Error())
		return
	}

	w.Write(resp.Bytes())
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
	if err != nil {
		h.Logger.Println(err.Error())
		http.Error(w, "UpVote:"+err.Error(), http.StatusBadRequest)
		return
	}

	tmpl := template.New("")
	tmpl, errParse := tmpl.Parse(postTemplate)
	if errParse != nil {
		fmt.Println("Error parsing", errParse.Error())
	}
	var resp bytes.Buffer

	errExecution := tmpl.Execute(&resp, foundPost)
	if errExecution != nil {
		fmt.Println("Error executing template:", errExecution.Error())
		return
	}

	w.Write(resp.Bytes())

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
	if err != nil {
		h.Logger.Println(err.Error())
		http.Error(w, "DownVote:"+err.Error(), http.StatusBadRequest)
		return
	}

	tmpl := template.New("")
	tmpl, errParse := tmpl.Parse(postTemplate)
	if errParse != nil {
		fmt.Println("Error parsing", errParse.Error())
	}
	var resp bytes.Buffer

	errExecution := tmpl.Execute(&resp, foundPost)
	if errExecution != nil {
		fmt.Println("Error executing template:", errExecution.Error())
		return
	}

	w.Write(resp.Bytes())
}

func (h *PostsHandler) DeletePost(w http.ResponseWriter, r *http.Request) {

}

func (h *PostsHandler) GetByCategory(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// postIDStr, ok := vars["category_name"]
	// if !ok {
	// 	h.Logger.Println("no category_name in query")
	// 	http.Error(w, "no category_name in query", http.StatusBadRequest)
	// 	return
	// }

	// u64, err := strconv.ParseUint(postIDStr, 10, 32)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// postID := uint(u64)

	// foundPost, errFind := h.PostsRepo.GetByID(postID)
	// if errFind != nil {
	// 	h.Logger.Println("no such post in repo.getbyid:", errFind.Error())
	// 	http.Error(w, fmt.Sprintln("no such post in repo.getbyid:", errFind.Error()), http.StatusBadRequest)
	// 	return
	// }

	// tmpl := template.New("")
	// tmpl, errParse := tmpl.Parse(postTemplate)
	// if errParse != nil {
	// 	fmt.Println("Error parsing in getPostByID method", errParse.Error())
	// }
	// var resp bytes.Buffer

	// errExecution := tmpl.Execute(&resp, foundPost)
	// if errExecution != nil {
	// 	fmt.Println("Error executing template in getPostByID:", errExecution.Error())
	// 	return
	// }

	// w.Write(resp.Bytes())
}

func (h *PostsHandler) GetAllByUser(w http.ResponseWriter, r *http.Request) {

}
