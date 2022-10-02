package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Row struct {
	Id             int    `xml:"id"`
	Guid           string `xml:"guid"`
	IsActive       bool   `xml:"isActive"`
	Balance        string `xml:"balance"`
	PictureURL     string `xml:"picture"`
	Age            int    `xml:"age"`
	EyeColor       string `xml:"eyeColor"`
	FirstName      string `xml:"first_name"`
	LastName       string `xml:"last_name"`
	Gender         string `xml:"gender"`
	Company        string `xml:"company"`
	Email          string `xml:"email"`
	Phone          string `xml:"phone"`
	Address        string `xml:"address"`
	About          string `xml:"about"`
	Registered     string `xml:"registered"`
	FavouriteFruit string `xml:"favoriteFruit"`
}

type Rows struct {
	Version string `xml:"version,attr"`
	List    []Row  `xml:"row"`
}

type Query struct {
	Name  string `json:"name"`
	About string `json:"about"`
}

type Queries struct {
	QueryList []Query
}

type OrderField struct {
	Id   int    `json:"id"`
	Age  int    `json:"age"`
	Name string `json:"name"`
}

type Orders struct {
	OrderList []OrderField
}

type Handler struct {
	Name string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	SearchServer(w, r)
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("dataset.xml")
	if err != nil {
		panic(err)
	}

	rows := Rows{}
	err = xml.Unmarshal([]byte(data), &rows)
	if err != nil {
		panic(err)
	}

	query := r.FormValue("query")
	queryJSON := Queries{}
	if query != "" {
		err = json.Unmarshal([]byte(query), &queryJSON)
		if err != nil {
			panic(err)
		}

		fmt.Println("query is", query)
	}
	orderField := r.FormValue("order_field")
	orderFieldJSON := Orders{}
	if orderField != "" {
		err = json.Unmarshal([]byte(orderField), &orderFieldJSON)
		if err != nil {
			panic(err)
		}

		fmt.Println("orderField is", orderField)
	}
	orderBy := r.FormValue("order_by")
	if orderBy != "" {
		fmt.Println("order_by is", orderBy)
	}
	limit := r.FormValue("limit")
	if limit != "" {
		fmt.Println("limit is", limit)
	}
	offset := r.FormValue("offset")
	if offset != "" {
		fmt.Println("offset is", offset)
	}

	responseBody := make([]User, len(orderFieldJSON.OrderList))

	for i, order := range orderFieldJSON.OrderList {
		responseBody[i] = User{
			ID:     rows.List[order.Id].Id,
			Name:   rows.List[order.Id].FirstName + rows.List[order.Id].LastName,
			Age:    rows.List[order.Id].Age,
			About:  rows.List[order.Id].About,
			Gender: rows.List[order.Id].Gender,
		}
	}

	// fmt.Println("response:")
	// for i, resp := range responseBody {
	// 	fmt.Println(i, resp)
	// }

	body, err := json.Marshal(responseBody)
	if err != nil {
		panic(err)
	}

	// fmt.Println(string(body))

	w.Write(body)
}

func clientHit() {
	time.Sleep(5 * time.Second)

	client := &SearchClient{
		AccessToken: "hello",
		URL:         "http://localhost:8080/",
	}

	request := SearchRequest{
		Limit:      25,
		Offset:     0,
		Query:      "{\"querylist\":[{\"name\":\"Ivan\"}, {\"name\":\"Ivan2\"}]}",
		OrderField: "{\"orderlist\":[{\"id\":1},{\"id\":2}]}",
		OrderBy:    0,
	}
	resp, err := client.FindUsers(request)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("resp", resp)
}

func main() {

	mux := http.NewServeMux()

	testHandler := &Handler{Name: "test"}
	mux.Handle("/test/", testHandler)

	rootHandler := &Handler{Name: "root"}
	mux.Handle("/", rootHandler)

	server := http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("starting server at :8080")
	go clientHit()
	server.ListenAndServe()
}
