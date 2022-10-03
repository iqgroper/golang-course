package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type ServerTestCase struct {
	URL         string
	AccessToken string
	Query       string
	OrderField  string
	OrderBy     int
	Response    string
	StatusCode  int
}

func TestSearchServer(t *testing.T) {
	cases := []ServerTestCase{
		{
			AccessToken: "hello",
			URL:         "http://localhost:8080/",
			Query:       "{\"querylist\":[{\"Name\":\"GlennJordan\"},{\"Name\":\"RoseCarney\"},{\"Name\":\"OwenLynn\"}]}",
			OrderField:  "Age",
			OrderBy:     1,
			Response:    `[{"ID":8,"Name":"GlennJordan","Age":29,"About":"Duis reprehenderit sit velit exercitation non aliqua magna quis ad excepteur anim. Eu cillum cupidatat sit magna cillum irure occaecat sunt officia officia deserunt irure. Cupidatat dolor cupidatat ipsum minim consequat Lorem adipisicing. Labore fugiat cupidatat nostrud voluptate ea eu pariatur non. Ipsum quis occaecat irure amet esse eu fugiat deserunt incididunt Lorem esse duis occaecat mollit.\n","Gender":"male"},{"ID":4,"Name":"OwenLynn","Age":30,"About":"Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n","Gender":"male"},{"ID":9,"Name":"RoseCarney","Age":36,"About":"Voluptate ipsum ad consequat elit ipsum tempor irure consectetur amet. Et veniam sunt in sunt ipsum non elit ullamco est est eu. Exercitation ipsum do deserunt do eu adipisicing id deserunt duis nulla ullamco eu. Ad duis voluptate amet quis commodo nostrud occaecat minim occaecat commodo. Irure sint incididunt est cupidatat laborum in duis enim nulla duis ut in ut. Cupidatat ex incididunt do ullamco do laboris eiusmod quis nostrud excepteur quis ea.\n","Gender":"female"}]`,
			StatusCode:  http.StatusOK,
		},
	}
	for caseNum, item := range cases {

		requestBody := SearchRequest{
			Limit:      25,
			Offset:     0,
			Query:      item.Query,
			OrderField: item.OrderField,
			OrderBy:    item.OrderBy,
		}

		body, err := json.Marshal(requestBody)
		if err != nil {
			fmt.Println(fmt.Errorf("%s", err))
		}

		requestURL := item.URL + "?order_field=" + item.OrderField + "&query=" + item.Query + "&order_by=" + strconv.Itoa(item.OrderBy)
		req := httptest.NewRequest("GET", requestURL, bytes.NewReader(body))
		req.Header.Add("Accesstoken", item.AccessToken)

		w := httptest.NewRecorder()

		SearchServer(w, req)

		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				caseNum, w.Code, item.StatusCode)
		}

		resp := w.Result()

		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("[%d] failed to read body: %v", caseNum, err)
		}

		bodyStr := string(responseBody)

		if bodyStr != item.Response {
			t.Errorf("[%d] wrong Response: got %+v, expected %+v",
				caseNum, bodyStr, item.Response)
		}
	}
}

type ClientTestCase struct {
	Request     SearchRequest
	AccessToken string
	Response    SearchResponse
}

func TestFindUsers(t *testing.T) {
	cases := []ClientTestCase{
		{
			Request: SearchRequest{
				Limit:      25,
				Offset:     0,
				Query:      "{\"querylist\":[{\"Name\":\"GlennJordan\"},{\"Name\":\"OwenLynn\"}]}",
				OrderField: "Age",
				OrderBy:    1,
			},
			AccessToken: "hello",
			Response: SearchResponse{
				Users: []User{
					{
						ID:     0,
						Name:   "GlennJordan",
						Age:    12,
						About:  "",
						Gender: "female",
					},
					{
						ID:     0,
						Name:   "OwenLynn",
						Age:    12,
						About:  "",
						Gender: "male",
					},
				},
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	for caseNum, item := range cases {

		client := &SearchClient{
			AccessToken: item.AccessToken,
			URL:         ts.URL,
		}

		result, _ := client.FindUsers(item.Request)

		fmt.Println(caseNum, result)

		// if err != nil && !item.IsError {
		// 	t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		// }
		// if err == nil && item.IsError {
		// 	t.Errorf("[%d] expected error, got nil", caseNum)
		// }
		// if !reflect.DeepEqual(item.Result, result) {
		// 	t.Errorf("[%d] wrong result, expected %#v, got %#v", caseNum, item.Result, result)
		// }
	}
	ts.Close()
}
