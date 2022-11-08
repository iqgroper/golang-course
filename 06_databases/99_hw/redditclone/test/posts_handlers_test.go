package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"redditclone/pkg/handlers"
	"redditclone/pkg/posts"
	"redditclone/pkg/session"
	"redditclone/pkg/user"
	"reflect"
	"testing"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"
)

func TestAddPostHandler(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := posts.NewMockPostRepo(ctrl)
	mockCommentRepo := posts.NewMockCommentRepo(ctrl)

	service := &handlers.PostsHandler{
		PostsRepo:    mockPostRepo,
		CommentsRepo: mockCommentRepo,
		Logger:       log.WithFields(log.Fields{}),
	}

	resultPost := &posts.Post{
		ID:               "0",
		Title:            "title",
		Score:            1,
		VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}},
		Category:         "category",
		Comments:         make([]*posts.Comment, 0, 10),
		CreatedDTTM:      time.Now().UTC(),
		Text:             "text",
		Type:             "text",
		UpvotePercentage: 100,
		Views:            0,
		Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
	}

	newUser := &user.User{ID: "0", Login: "login", Password: "asdfasdf"}

	newPost := &posts.NewPost{
		Type:     "text",
		Title:    "title",
		Text:     "text",
		Category: "category",
		Author:   *newUser,
	}

	mockPostRepo.EXPECT().Add(newPost).Return(resultPost, nil)

	reqBody, _ := json.Marshal(newPost)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/posts", bytes.NewReader(reqBody))

	sess := &session.Session{
		ID:   "session id",
		User: newUser,
	}
	ctx := session.ContextWithSession(r.Context(), sess)

	service.AddPost(w, r.WithContext(ctx))

	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("wrong status code, got: %d, expected 200", resp.StatusCode)
		return
	}

	body, _ := io.ReadAll(resp.Body)

	recievedPost := &posts.Post{}

	err := json.Unmarshal(body, recievedPost)
	if err != nil {
		t.Errorf("Error unmarshalling resp body: %s", err.Error())
		return
	}

	if !reflect.DeepEqual(resultPost, recievedPost) {
		t.Errorf("Wrong result\nExpected: %v\nRecieved: %v", resultPost, recievedPost)
		return
	}
	resp.Body.Close()

	//error in add
	mockPostRepo.EXPECT().Add(newPost).Return(nil, fmt.Errorf("error in add"))

	reqBody, _ = json.Marshal(newPost)

	w = httptest.NewRecorder()
	r = httptest.NewRequest("POST", "/api/posts", bytes.NewReader(reqBody))

	service.AddPost(w, r.WithContext(ctx))

	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("wrong status code, got: %d, expected 200", resp.StatusCode)
		return
	}
}

func TestAddCommentHandler(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := posts.NewMockPostRepo(ctrl)
	mockCommentRepo := posts.NewMockCommentRepo(ctrl)

	service := &handlers.PostsHandler{
		PostsRepo:    mockPostRepo,
		CommentsRepo: mockCommentRepo,
		Logger:       log.WithFields(log.Fields{}),
	}

	newUser := &user.User{ID: "0", Login: "login", Password: "asdfasdf"}

	resultPost := &posts.Post{
		ID:        "0",
		Title:     "title",
		Score:     1,
		VotesList: []posts.VoteStruct{{User: "author.login", Vote: 1}},
		Category:  "category",
		Comments: []*posts.Comment{
			{
				ID:      "0",
				Body:    "newcomment",
				PostID:  "0",
				Created: time.Now().UTC(),
				Author:  newUser,
			},
		},
		CreatedDTTM:      time.Now().UTC(),
		Text:             "text",
		Type:             "text",
		UpvotePercentage: 100,
		Views:            0,
		Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
	}

	sess := &session.Session{
		ID:   "session id",
		User: newUser,
	}
	reqBody, _ := json.Marshal(struct{ Comment string }{"newcomment"})

	//no post id in query
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/posts/", bytes.NewReader(reqBody))
	ctxr := session.ContextWithSession(req.Context(), sess)

	service.AddComment(w, req.WithContext(ctxr))

	resp := w.Result()
	if resp.StatusCode != 400 {
		t.Errorf("wrong status code, got: %d, expected 400", resp.StatusCode)
		return
	}

	//SUCCESSFUL REQ
	mockCommentRepo.EXPECT().Add("0", "newcomment", newUser).Return(nil, nil)
	mockPostRepo.EXPECT().GetByID("0").Return(resultPost, nil)

	w = httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/posts/0", bytes.NewReader(reqBody))
	vars := map[string]string{
		"post_id": "0",
	}
	r = mux.SetURLVars(r, vars)
	ctx := session.ContextWithSession(r.Context(), sess)

	service.AddComment(w, r.WithContext(ctx))

	resp = w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("wrong status code, got: %d, expected 200", resp.StatusCode)
		return
	}

	body, _ := io.ReadAll(resp.Body)

	recievedPost := &posts.Post{}

	err := json.Unmarshal(body, recievedPost)
	if err != nil {
		t.Errorf("Error unmarshalling resp body: %s", err.Error())
		return
	}

	//deep equal не подходит потому что там ссылки и они на память
	require.Equal(t, resultPost.Comments[0].Body, recievedPost.Comments[0].Body)
	require.Equal(t, resultPost.Comments[0].Author.ID, recievedPost.Comments[0].Author.ID)
	require.Equal(t, resultPost.Comments[0].Author.Login, recievedPost.Comments[0].Author.Login)
	require.Equal(t, resultPost.Comments[0].Body, recievedPost.Comments[0].Body)
	require.Equal(t, resultPost.Comments[0].ID, recievedPost.Comments[0].ID)
	require.Equal(t, resultPost.Comments[0].Author.ID, recievedPost.Comments[0].Author.ID)
	require.Equal(t, resultPost.Comments[0].Author.Login, recievedPost.Comments[0].Author.Login)
	require.Equal(t, resultPost.Comments[0].Body, recievedPost.Comments[0].Body)
	require.Equal(t, resultPost.Comments[0].ID, recievedPost.Comments[0].ID)
	resp.Body.Close()

	//blank comment
	reqBody, _ = json.Marshal(struct{ Comment string }{""})

	w = httptest.NewRecorder()
	r = httptest.NewRequest("POST", "/api/posts/0", bytes.NewReader(reqBody))
	vars = map[string]string{
		"post_id": "0",
	}
	r = mux.SetURLVars(r, vars)

	service.AddComment(w, r.WithContext(ctx))

	resp = w.Result()
	if resp.StatusCode != 422 {
		t.Errorf("wrong status code, got: %d, expected 200", resp.StatusCode)
		return
	}

	//error in Add
	mockCommentRepo.EXPECT().Add("0", "newcomment", newUser).Return(nil, fmt.Errorf("add error"))

	reqBody, _ = json.Marshal(struct{ Comment string }{"newcomment"})

	w = httptest.NewRecorder()
	r = httptest.NewRequest("POST", "/api/posts/0", bytes.NewReader(reqBody))
	vars = map[string]string{
		"post_id": "0",
	}
	r = mux.SetURLVars(r, vars)

	service.AddComment(w, r.WithContext(ctx))

	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("wrong status code, got: %d, expected 200", resp.StatusCode)
		return
	}

	//error in GetByID
	mockCommentRepo.EXPECT().Add("0", "newcomment", newUser).Return(nil, nil)
	mockPostRepo.EXPECT().GetByID("0").Return(nil, fmt.Errorf("error in getbyid"))

	reqBody, _ = json.Marshal(struct{ Comment string }{"newcomment"})

	w = httptest.NewRecorder()
	r = httptest.NewRequest("POST", "/api/posts/0", bytes.NewReader(reqBody))
	vars = map[string]string{
		"post_id": "0",
	}
	r = mux.SetURLVars(r, vars)

	service.AddComment(w, r.WithContext(ctx))

	resp = w.Result()
	if resp.StatusCode != 400 {
		t.Errorf("wrong status code, got: %d, expected 200", resp.StatusCode)
		return
	}
}

func TestDeleteCommentHandler(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := posts.NewMockPostRepo(ctrl)
	mockCommentRepo := posts.NewMockCommentRepo(ctrl)

	service := &handlers.PostsHandler{
		PostsRepo:    mockPostRepo,
		CommentsRepo: mockCommentRepo,
		Logger:       log.WithFields(log.Fields{}),
	}

	newUser := &user.User{ID: "0", Login: "login", Password: "asdfasdf"}

	resultPost := &posts.Post{
		ID:        "0",
		Title:     "title",
		Score:     1,
		VotesList: []posts.VoteStruct{{User: "author.login", Vote: 1}},
		Category:  "category",
		Comments: []*posts.Comment{
			{
				ID:      "0",
				Body:    "newcomment",
				PostID:  "0",
				Created: time.Now().UTC(),
				Author:  newUser,
			},
		},
		CreatedDTTM:      time.Now().UTC(),
		Text:             "text",
		Type:             "text",
		UpvotePercentage: 100,
		Views:            0,
		Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
	}
	expectedPost := &posts.Post{
		ID:               "0",
		Title:            "title",
		Score:            1,
		VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}},
		Category:         "category",
		Comments:         []*posts.Comment{},
		CreatedDTTM:      time.Now().UTC(),
		Text:             "text",
		Type:             "text",
		UpvotePercentage: 100,
		Views:            0,
		Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
	}

	sess := &session.Session{
		ID:   "session id",
		User: newUser,
	}
	//no post id in query
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/api/posts/", nil)
	ctxr := session.ContextWithSession(req.Context(), sess)

	service.DeleteComment(w, req.WithContext(ctxr))

	resp := w.Result()
	if resp.StatusCode != 400 {
		t.Errorf("wrong status code, got: %d, expected 400", resp.StatusCode)
		return
	}

	//no comment id in query
	w = httptest.NewRecorder()
	req = httptest.NewRequest("DELETE", "/api/posts/0/", nil)
	ctxr = session.ContextWithSession(req.Context(), sess)

	service.DeleteComment(w, req.WithContext(ctxr))

	resp = w.Result()
	if resp.StatusCode != 400 {
		t.Errorf("wrong status code, got: %d, expected 400", resp.StatusCode)
		return
	}

	//successful deletion
	mockCommentRepo.EXPECT().Delete("0", "0").Return(true, nil)
	mockPostRepo.EXPECT().GetByID("0").Return(resultPost, nil)
	mockCommentRepo.EXPECT().GetAll("0").Return([]*posts.Comment{}, nil)

	w = httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/posts/0/0", nil)
	vars := map[string]string{
		"post_id":    "0",
		"comment_id": "0",
	}
	r = mux.SetURLVars(r, vars)
	ctx := session.ContextWithSession(r.Context(), sess)

	service.DeleteComment(w, r.WithContext(ctx))

	resp = w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("wrong status code, got: %d, expected 200", resp.StatusCode)
		return
	}

	body, _ := io.ReadAll(resp.Body)

	recievedPost := &posts.Post{}

	err := json.Unmarshal(body, recievedPost)
	if err != nil {
		t.Errorf("Error unmarshalling resp body: %s", err.Error())
		return
	}

	//deepEqual не подходит потому что там ссылки и они на память
	require.Equal(t, len(expectedPost.Comments), len(recievedPost.Comments))
	resp.Body.Close()

	// //error in Delete
	mockCommentRepo.EXPECT().Delete("0", "0").Return(false, fmt.Errorf("error in delet function"))

	w = httptest.NewRecorder()
	r = httptest.NewRequest("POST", "/api/posts/0/0", nil)
	vars = map[string]string{
		"post_id":    "0",
		"comment_id": "0",
	}
	r = mux.SetURLVars(r, vars)

	service.DeleteComment(w, r.WithContext(ctx))

	resp = w.Result()
	if resp.StatusCode != 400 {
		t.Errorf("wrong status code, got: %d, expected 400", resp.StatusCode)
		return
	}

	// //error in GetByID
	mockCommentRepo.EXPECT().Delete("0", "0").Return(true, nil)
	mockPostRepo.EXPECT().GetByID("0").Return(nil, fmt.Errorf("error in getByID"))

	w = httptest.NewRecorder()
	r = httptest.NewRequest("POST", "/api/posts/0/0", nil)
	vars = map[string]string{
		"post_id":    "0",
		"comment_id": "0",
	}
	r = mux.SetURLVars(r, vars)

	service.DeleteComment(w, r.WithContext(ctx))

	resp = w.Result()
	if resp.StatusCode != 400 {
		t.Errorf("wrong status code, got: %d, expected 400", resp.StatusCode)
		return
	}

	//error in GetAll
	mockCommentRepo.EXPECT().Delete("0", "0").Return(true, nil)
	mockPostRepo.EXPECT().GetByID("0").Return(resultPost, nil)
	mockCommentRepo.EXPECT().GetAll("0").Return(nil, fmt.Errorf("error in GetAll"))

	w = httptest.NewRecorder()
	r = httptest.NewRequest("POST", "/api/posts/0/0", nil)
	vars = map[string]string{
		"post_id":    "0",
		"comment_id": "0",
	}
	r = mux.SetURLVars(r, vars)

	service.DeleteComment(w, r.WithContext(ctx))

	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("wrong status code, got: %d, expected 400", resp.StatusCode)
		return
	}
}

func TestGetPostByIDHandler(t *testing.T) {

	// мы передаём t сюда, это надо чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	mockPostRepo := posts.NewMockPostRepo(ctrl)
	mockCommentRepo := posts.NewMockCommentRepo(ctrl)
	service := &handlers.PostsHandler{
		PostsRepo:    mockPostRepo,
		CommentsRepo: mockCommentRepo,
		Logger:       log.WithFields(log.Fields{}),
	}

	resultPost := &posts.Post{
		ID:               "0",
		Title:            "title",
		Score:            1,
		VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}},
		Category:         "category",
		Comments:         make([]*posts.Comment, 0, 10),
		CreatedDTTM:      time.Now().UTC(),
		Text:             "text",
		Type:             "text",
		UpvotePercentage: 100,
		Views:            0,
		Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
	}

	mockPostRepo.EXPECT().GetByID("0").Return(resultPost, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/post/0", nil)
	vars := map[string]string{
		"post_id": "0",
	}
	r = mux.SetURLVars(r, vars)

	service.GetByID(w, r)

	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("wrong status code, got: %d, expected 200", resp.StatusCode)
		return
	}

	body, _ := io.ReadAll(resp.Body)

	recievedPost := &posts.Post{}

	err := json.Unmarshal(body, &recievedPost)
	if err != nil {
		t.Errorf("Error unmarshalling resp body: %s", err.Error())
		return
	}

	if !reflect.DeepEqual(resultPost, recievedPost) {
		t.Errorf("Wrong result\nExpected: %v\nRecieved: %v", resultPost, recievedPost)
		return
	}
	resp.Body.Close()

	//no params in query
	w = httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/post/", nil)

	service.GetByID(w, req)

	resp = w.Result()
	if resp.StatusCode != 400 {
		t.Errorf("wrong status code, got: %d, expected 400", resp.StatusCode)
		return
	}

	//getting wrong post_id
	mockPostRepo.EXPECT().GetByID("10").Return(nil, posts.ErrNoPost)

	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "/api/post/10", nil)
	vars = map[string]string{
		"post_id": "10",
	}
	r = mux.SetURLVars(r, vars)

	service.GetByID(w, r)

	resp = w.Result()
	if resp.StatusCode != 400 {
		t.Errorf("wrong status code, got: %d, expected 400", resp.StatusCode)
		return
	}
}

func TestGetByCategoryHandler(t *testing.T) {

	// мы передаём t сюда, это надо чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	mockPostRepo := posts.NewMockPostRepo(ctrl)
	mockCommentRepo := posts.NewMockCommentRepo(ctrl)
	service := &handlers.PostsHandler{
		PostsRepo:    mockPostRepo,
		CommentsRepo: mockCommentRepo,
		Logger:       log.WithFields(log.Fields{}),
	}

	resultPosts := []*posts.Post{
		{
			ID:               "0",
			Title:            "title",
			Score:            1,
			VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}},
			Category:         "funny",
			Comments:         make([]*posts.Comment, 0, 10),
			CreatedDTTM:      time.Now().UTC(),
			Text:             "text",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            0,
			Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
		},
		{
			ID:               "1",
			Title:            "title1",
			Score:            1,
			VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}},
			Category:         "funny",
			Comments:         make([]*posts.Comment, 0, 10),
			CreatedDTTM:      time.Now().UTC(),
			Text:             "text1",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            0,
			Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
		},
	}

	mockPostRepo.EXPECT().GetAllByCategory("funny").Return(resultPosts, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/posts/funny", nil)
	vars := map[string]string{
		"category_name": "funny",
	}
	r = mux.SetURLVars(r, vars)

	service.GetByCategory(w, r)

	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("wrong status code, got: %d, expected 200", resp.StatusCode)
		return
	}

	body, _ := io.ReadAll(resp.Body)

	recievedPosts := make([]*posts.Post, 0, 2)

	err := json.Unmarshal(body, &recievedPosts)
	if err != nil {
		t.Errorf("Error unmarshalling resp body: %s", err.Error())
		return
	}

	for idx := range resultPosts {
		if !reflect.DeepEqual(resultPosts[idx], recievedPosts[idx]) {
			t.Errorf("Wrong result\nExpected: %v\nRecieved: %v", resultPosts[idx], recievedPosts[idx])
			return
		}
	}
	resp.Body.Close()

	//nothing to return
	mockPostRepo.EXPECT().GetAllByCategory("funny").Return(nil, fmt.Errorf("no such post found"))

	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "/api/posts/funny", nil)
	vars = map[string]string{
		"category_name": "funny",
	}
	r = mux.SetURLVars(r, vars)

	service.GetByCategory(w, r)

	resp = w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("wrong status code, got: %d, expected 200", resp.StatusCode)
		return
	}

	body, _ = io.ReadAll(resp.Body)

	if string(body) != "[]" {
		t.Errorf("Wrong result\nExpected: %v\nRecieved: %v", resultPosts, recievedPosts)
		return
	}
	resp.Body.Close()

	//no params in query

	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "/api/posts/", nil)

	service.GetByCategory(w, r)

	resp = w.Result()
	if resp.StatusCode != 400 {
		t.Errorf("wrong status code, got: %d, expected 400", resp.StatusCode)
		return
	}
}

func TestGetAllByUserHandler(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := posts.NewMockPostRepo(ctrl)
	mockCommentRepo := posts.NewMockCommentRepo(ctrl)
	service := &handlers.PostsHandler{
		PostsRepo:    mockPostRepo,
		CommentsRepo: mockCommentRepo,
		Logger:       log.WithFields(log.Fields{}),
	}

	resultPosts := []*posts.Post{
		{
			ID:               "0",
			Title:            "title",
			Score:            1,
			VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}},
			Category:         "funny",
			Comments:         make([]*posts.Comment, 0, 10),
			CreatedDTTM:      time.Now().UTC(),
			Text:             "text",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            0,
			Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
		},
		{
			ID:               "1",
			Title:            "title1",
			Score:            1,
			VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}},
			Category:         "funny",
			Comments:         make([]*posts.Comment, 0, 10),
			CreatedDTTM:      time.Now().UTC(),
			Text:             "text1",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            0,
			Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
		},
	}

	mockPostRepo.EXPECT().GetByUser("author.login").Return(resultPosts, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/user/", nil)
	vars := map[string]string{
		"user_login": "author.login",
	}
	r = mux.SetURLVars(r, vars)

	service.GetAllByUser(w, r)

	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("wrong status code, got: %d, expected 200", resp.StatusCode)
		return
	}

	body, _ := io.ReadAll(resp.Body)

	recievedPosts := make([]*posts.Post, 0, 2)

	err := json.Unmarshal(body, &recievedPosts)
	if err != nil {
		t.Errorf("Error unmarshalling resp body: %s", err.Error())
		return
	}

	for idx := range resultPosts {
		if !reflect.DeepEqual(resultPosts[idx], recievedPosts[idx]) {
			t.Errorf("Wrong result\nExpected: %v\nRecieved: %v", resultPosts[idx], recievedPosts[idx])
			return
		}
	}
	resp.Body.Close()

	//nothing to return
	mockPostRepo.EXPECT().GetByUser("author.login").Return(nil, posts.ErrNoPost)

	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "/api/user/", nil)
	vars = map[string]string{
		"user_login": "author.login",
	}
	r = mux.SetURLVars(r, vars)

	service.GetAllByUser(w, r)

	resp = w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("wrong status code, got: %d, expected 200", resp.StatusCode)
		return
	}

	body, _ = io.ReadAll(resp.Body)

	if string(body) != "[]" {
		t.Errorf("Wrong result\nExpected: %v\nRecieved: %v", resultPosts, recievedPosts)
		return
	}
	resp.Body.Close()

	//no params in query

	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "/api/user/", nil)

	service.GetAllByUser(w, r)

	resp = w.Result()
	if resp.StatusCode != 400 {
		t.Errorf("wrong status code, got: %d, expected 400", resp.StatusCode)
		return
	}
}

func TestGetAllPostsHandler(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := posts.NewMockPostRepo(ctrl)
	mockCommentRepo := posts.NewMockCommentRepo(ctrl)
	service := &handlers.PostsHandler{
		PostsRepo:    mockPostRepo,
		CommentsRepo: mockCommentRepo,
		Logger:       log.WithFields(log.Fields{}),
	}

	resultPosts := []*posts.Post{
		{
			ID:               "0",
			Title:            "title",
			Score:            1,
			VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}},
			Category:         "category",
			Comments:         make([]*posts.Comment, 0, 10),
			CreatedDTTM:      time.Now().UTC(),
			Text:             "text",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            0,
			Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
		},
		{
			ID:               "1",
			Title:            "title1",
			Score:            1,
			VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}},
			Category:         "category1",
			Comments:         make([]*posts.Comment, 0, 10),
			CreatedDTTM:      time.Now().UTC(),
			Text:             "text1",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            0,
			Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
		},
	}

	mockPostRepo.EXPECT().GetAll().Return(resultPosts, nil)

	r := httptest.NewRequest("GET", "/api/posts/", nil)
	w := httptest.NewRecorder()

	service.GetAll(w, r)

	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("wrong status code, got: %d, expected 200", resp.StatusCode)
		return
	}

	body, _ := io.ReadAll(resp.Body)

	recievedPosts := make([]*posts.Post, 0, 2)

	err := json.Unmarshal(body, &recievedPosts)
	if err != nil {
		t.Errorf("Error unmarshalling resp body: %s", err.Error())
		return
	}

	for idx := range resultPosts {
		if !reflect.DeepEqual(resultPosts[idx], recievedPosts[idx]) {
			t.Errorf("Wrong result\nExpected: %v\nRecieved: %v", resultPosts[idx], recievedPosts[idx])
			return
		}
	}
	resp.Body.Close()

	//nothing to return
	mockPostRepo.EXPECT().GetAll().Return(nil, fmt.Errorf("no such post found"))

	r = httptest.NewRequest("GET", "/api/posts/", nil)
	w = httptest.NewRecorder()

	service.GetAll(w, r)

	resp = w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("wrong status code, got: %d, expected 200", resp.StatusCode)
		return
	}

	body, _ = io.ReadAll(resp.Body)

	if string(body) != "[]" {
		t.Errorf("Wrong result\nExpected: %v\nRecieved: %v", resultPosts, recievedPosts)
		return
	}
	resp.Body.Close()
}
