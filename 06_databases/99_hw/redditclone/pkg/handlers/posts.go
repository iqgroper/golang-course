package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"redditclone/pkg/posts"
	"redditclone/pkg/session"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type PostsHandler struct {
	PostsRepo    posts.PostRepo
	CommentsRepo posts.CommentRepo
	Logger       *logrus.Entry
}

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

	resp, errEncoding := json.Marshal(addedPost)
	if errEncoding != nil {
		h.Logger.Println("error marshalling post:", errEncoding.Error())
		http.Error(w, fmt.Sprintf(`error marshalling post: %s`, errEncoding.Error()), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

func (h *PostsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	postID, ok := vars["post_id"]
	if !ok {
		h.Logger.Println("no post_id in query in GetByID")
		http.Error(w, "no post_id in query in GetByID", http.StatusBadRequest)
		return
	}

	foundPost, errFind := h.PostsRepo.GetByID(postID)
	if errFind != nil {
		h.Logger.Println("no such post in repo.getbyid:", errFind.Error())
		http.Error(w, fmt.Sprintln("no such post in repo.getbyid:", errFind.Error()), http.StatusBadRequest)
		return
	}

	foundPost.Views += 1

	resp, errEncoding := json.Marshal(foundPost)
	if errEncoding != nil {
		h.Logger.Println("error marshalling post:", errEncoding.Error())
		http.Error(w, fmt.Sprintf(`error marshalling post: %s`, errEncoding.Error()), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

func (h *PostsHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	postID, ok := vars["post_id"]
	if !ok {
		h.Logger.Println("no post_id in query in AddComment")
		http.Error(w, "no post_id in query in AddComment", http.StatusBadRequest)
		return
	}

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

	_, errAdd := h.CommentsRepo.Add(postID, body.Comment, sess.User)
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

	resp, errEncoding := json.Marshal(foundPost)
	if errEncoding != nil {
		h.Logger.Println("error marshalling post:", errEncoding.Error())
		http.Error(w, fmt.Sprintf(`error marshalling post: %s`, errEncoding.Error()), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

func (h *PostsHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	postID, ok := vars["post_id"]
	if !ok {
		h.Logger.Println("no post_id in query in DeleteComment")
		http.Error(w, "no post_id in query in DeleteComment", http.StatusBadRequest)
		return
	}
	commentID, ok := vars["comment_id"]
	if !ok {
		h.Logger.Println("no comment_id in query in DeleteComment")
		http.Error(w, "no comment_id in query in DeleteComment", http.StatusBadRequest)
		return
	}

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

	resp, errEncoding := json.Marshal(foundPost)
	if errEncoding != nil {
		h.Logger.Println("error marshalling post:", errEncoding.Error())
		http.Error(w, fmt.Sprintf(`error marshalling post: %s`, errEncoding.Error()), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

func (h *PostsHandler) UpVote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	postID, ok := vars["post_id"]
	if !ok {
		h.Logger.Println("no post_id in query in UpVote")
		http.Error(w, "no post_id in query in UpVote", http.StatusBadRequest)
		return
	}

	sess, errSess := session.SessionFromContext(r.Context())
	if errSess != nil {
		h.Logger.Println("error getting session in UpVote:", errSess.Error())
		http.Error(w, fmt.Sprintf(`error getting session in UpVote: %s`, errSess.Error()), http.StatusInternalServerError)
		return
	}

	foundPost, err := h.PostsRepo.UpVote(postID, sess.User.Login)
	if err == posts.ErrNoCanDo {
		h.Logger.Println("no such post in UpVote:", err.Error(), " redirected to unvote")
		http.Redirect(w, r, fmt.Sprintf("/api/post/%s/unvote", postID), http.StatusFound)
		return
	}
	if err != nil {
		h.Logger.Println(err.Error())
		http.Error(w, "UpVote:"+err.Error(), http.StatusBadRequest)
		return
	}

	resp, errEncoding := json.Marshal(foundPost)
	if errEncoding != nil {
		h.Logger.Println("error marshalling post:", errEncoding.Error())
		http.Error(w, fmt.Sprintf(`error marshalling post: %s`, errEncoding.Error()), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

func (h *PostsHandler) DownVote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	postID, ok := vars["post_id"]
	if !ok {
		h.Logger.Println("no post_id in query in DownVote")
		http.Error(w, "no post_id in query in DownVote", http.StatusBadRequest)
		return
	}

	sess, errSess := session.SessionFromContext(r.Context())
	if errSess != nil {
		h.Logger.Println("error getting session in DownVote:", errSess.Error())
		http.Error(w, fmt.Sprintf(`error getting session in DownVote: %s`, errSess.Error()), http.StatusInternalServerError)
		return
	}

	foundPost, err := h.PostsRepo.DownVote(postID, sess.User.Login)
	if err == posts.ErrNoCanDo {
		h.Logger.Println("in DownVote:", err.Error(), " redirected to unvote")
		http.Redirect(w, r, fmt.Sprintf("/api/post/%s/unvote", postID), http.StatusFound)
		return
	}
	if err != nil {
		h.Logger.Println(err.Error())
		http.Error(w, "DownVote:"+err.Error(), http.StatusBadRequest)
		return
	}

	resp, errEncoding := json.Marshal(foundPost)
	if errEncoding != nil {
		h.Logger.Println("error marshalling post:", errEncoding.Error())
		http.Error(w, fmt.Sprintf(`error marshalling post: %s`, errEncoding.Error()), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

func (h *PostsHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	postID, ok := vars["post_id"]
	if !ok {
		h.Logger.Println("no post_id in query in DeletePost")
		http.Error(w, "no post_id in query in DeletePost", http.StatusBadRequest)
		return
	}

	_, errDel := h.PostsRepo.Delete(postID)
	if errDel != nil {
		h.Logger.Println("in DeletePost:", errDel.Error())
		http.Error(w, "in DeletePost: "+errDel.Error(), http.StatusBadRequest)
		return
	}

	resp, errMar := json.Marshal(map[string]string{"message": "success"})
	if errMar != nil {
		h.Logger.Println("cant marshal in DeleteCommentByPost", errMar.Error())
		http.Error(w, "cant marshal in DeleteCommentByPost "+errMar.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}

func (h *PostsHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	postList, errGetting := h.PostsRepo.GetAll()
	if errGetting != nil {
		h.Logger.Println(errGetting)
		w.Write([]byte("[]"))
		return
	}

	resp, errEncoding := json.Marshal(postList)
	if errEncoding != nil {
		h.Logger.Println("error marshalling post:", errEncoding.Error())
		http.Error(w, fmt.Sprintf(`error marshalling post: %s`, errEncoding.Error()), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
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

	resp, errEncoding := json.Marshal(postList)
	if errEncoding != nil {
		h.Logger.Println("error marshalling post:", errEncoding.Error())
		http.Error(w, fmt.Sprintf(`error marshalling post: %s`, errEncoding.Error()), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
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

	resp, errEncoding := json.Marshal(postList)
	if errEncoding != nil {
		h.Logger.Println("error marshalling post:", errEncoding.Error())
		http.Error(w, fmt.Sprintf(`error marshalling post: %s`, errEncoding.Error()), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

func (h *PostsHandler) UnVote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	postID, ok := vars["post_id"]
	if !ok {
		h.Logger.Println("no post_id in query in UnVote")
		http.Error(w, "no post_id in query in UnVote", http.StatusBadRequest)
		return
	}

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

	resp, errEncoding := json.Marshal(foundPost)
	if errEncoding != nil {
		h.Logger.Println("error marshalling post:", errEncoding.Error())
		http.Error(w, fmt.Sprintf(`error marshalling post: %s`, errEncoding.Error()), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}
