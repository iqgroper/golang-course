package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"redditclone/pkg/handlers"
	"redditclone/pkg/posts"
	"reflect"
	"testing"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/golang/mock/gomock"
)

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

func TestAddPostHandler(t *testing.T) {

	// ctrl := gomock.NewController(t)
	// defer ctrl.Finish()

	// mockPostRepo := posts.NewMockPostRepo(ctrl)
	// mockCommentRepo := posts.NewMockCommentRepo(ctrl)
	// service := &handlers.PostsHandler{
	// 	PostsRepo:    mockPostRepo,
	// 	CommentsRepo: mockCommentRepo,
	// 	Logger:       log.WithFields(log.Fields{}),
	// }

	// resultPost := &posts.Post{
	// 	ID:               "0",
	// 	Title:            "title",
	// 	Score:            1,
	// 	VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}},
	// 	Category:         "category",
	// 	Comments:         make([]*posts.Comment, 0, 10),
	// 	CreatedDTTM:      time.Now().UTC(),
	// 	Text:             "text",
	// 	Type:             "text",
	// 	UpvotePercentage: 100,
	// 	Views:            0,
	// 	Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
	// }

	// mockPostRepo.EXPECT().Add(resultPost).Return(resultPost, nil)

	// newPost := posts.NewPost{
	// 	Type:     "text",
	// 	Title:    "title",
	// 	Text:     "text",
	// 	Category: "category",
	// }

	// reqBody, _ := json.Marshal(newPost)

	// r := httptest.NewRequest("POST", "/api/posts", bytes.NewReader(reqBody))
	// w := httptest.NewRecorder()

	// service.AddPost(w, r)

	// resp := w.Result()
	// if resp.StatusCode != 200 {
	// 	t.Errorf("wrong status code, got: %d, expected 200", resp.StatusCode)
	// 	return
	// }

	// body, _ := io.ReadAll(resp.Body)

	// recievedPost := &posts.Post{}

	// err := json.Unmarshal(body, recievedPost)
	// if err != nil {
	// 	t.Errorf("Error unmarshalling resp body: %s", err.Error())
	// 	return
	// }

	// if !reflect.DeepEqual(resultPost, recievedPost) {
	// 	t.Errorf("Wrong result\nExpected: %v\nRecieved: %v", resultPost, recievedPost)
	// 	return
	// }
	// resp.Body.Close()
}
