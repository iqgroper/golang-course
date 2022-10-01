package main

import (
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

type Handler struct {
	Name string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "Name:", h.Name, "URL:", r.URL.String())
	SearchServer(w, r)
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("dataset.xml")
	if err != nil {
		panic(err)
	}

	rows := new(Rows)
	err = xml.Unmarshal([]byte(data), &rows)
	if err != nil {
		panic(err)
	}

	query := r.FormValue("query")
	if query != "" {
		fmt.Println("query is", query)
	}

	orderField := r.FormValue("order_field")
	if orderField != "" {
		fmt.Println("orderField is", orderField)
	}
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
		Query:      "{\"name\":\"Ivan\"}",
		OrderField: "{\"id\":0}",
		OrderBy:    0,
	}
	resp, err := client.FindUsers(request)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp)
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
