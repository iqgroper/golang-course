package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
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

type Handler struct {
	Name string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	SearchServer(w, r)
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	//обработать ситуации когда нулевые значения query orderfield
	//перед каждой паникой отправлять итернал еррор

	if r.Header["Accesstoken"][0] != "hello" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	data, err := ioutil.ReadFile("dataset.xml")
	if err != nil {
		panic(err) //internal server error: error w/ filename
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
	if orderField != "" {
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

	responseBody := make([]User, len(queryJSON.QueryList))
	i := 0
	for _, user := range rows.List {
		for _, order := range queryJSON.QueryList {
			if user.FirstName+user.LastName == order.Name || user.About == order.About {
				responseBody[i] = User{
					ID:     user.Id,
					Name:   user.FirstName + user.LastName,
					Age:    user.Age,
					About:  user.About,
					Gender: user.Gender,
				}
				i++
			}
		}
	}

	if !(orderBy == "Name" || orderBy == "Age" || orderBy == "Id") {
		w.WriteHeader(http.StatusBadRequest)
		error := SearchErrorResponse{
			Error: "OrderField invalid",
		}
		result, err := json.Marshal(error)
		if err != nil {
			panic(err)
		}
		w.Write(result)
		return
	}

	switch orderBy {
	case "0":
		break
	case "1":
		sort.Slice(responseBody, func(i, j int) bool {
			switch orderField {
			case "Id":
				return responseBody[i].ID < responseBody[j].ID
			case "Age":
				return responseBody[i].Age < responseBody[j].Age
			case "Name":
				return responseBody[i].Name < responseBody[j].Name
			default:
				return false
			}
		})
	case "-1":
		sort.Slice(responseBody, func(i, j int) bool {
			switch orderField {
			case "Id":
				return responseBody[i].ID > responseBody[j].ID
			case "Age":
				return responseBody[i].Age > responseBody[j].Age
			case "Name":
				return responseBody[i].Name > responseBody[j].Name
			default:
				return false
			}
		})
	default:
		w.WriteHeader(http.StatusBadRequest)
		error := SearchErrorResponse{
			Error: "OrderBy invalid",
		}
		result, err := json.Marshal(error)
		if err != nil {
			panic(err)
		}
		w.Write(result)
		return
	}

	body, err := json.Marshal(responseBody)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
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
		Query:      "{\"querylist\":[{\"name\":\"GlennJordan\"}, {\"name\":\"RoseCarney\"}, {\"name\":\"OwenLynn\"}]}",
		OrderField: "Agge",
		OrderBy:    -1,
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
