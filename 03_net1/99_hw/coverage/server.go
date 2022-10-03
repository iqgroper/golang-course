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
	ID             int    `xml:"id"`
	GUID           string `xml:"guid"`
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

var filename = "dataset.xml"
var correctFilename = "dataset.xml"

func SearchServer(w http.ResponseWriter, r *http.Request) {

	if r.Header["Accesstoken"][0] != "hello" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := fmt.Errorf("%s", err)
		fmt.Println(error)
		return
	}

	rows := Rows{}
	err1 := xml.Unmarshal(data, &rows)
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(fmt.Errorf("%s", err1))
		return
	}

	query := r.FormValue("query")
	queryJSON := Queries{}
	if query != "" {
		err = json.Unmarshal([]byte(query), &queryJSON)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(fmt.Errorf("%s", err))
			return
		}
		// fmt.Println("query is", query)
	}

	orderField := r.FormValue("order_field")
	if orderField == "" {
		orderField = "Name"
	}

	orderBy := r.FormValue("order_by")
	fmt.Println("orderby", orderBy)

	var responseBody []User
	if query == "" {
		responseBody = make([]User, len(rows.List))
		i := 0
		for _, user := range rows.List {
			responseBody[i] = User{
				ID:     user.ID,
				Name:   user.FirstName + user.LastName,
				Age:    user.Age,
				About:  user.About,
				Gender: user.Gender,
			}
			i++
		}

	} else {
		responseBody = make([]User, len(queryJSON.QueryList))
		i := 0
		for _, user := range rows.List {
			for _, order := range queryJSON.QueryList {
				if user.FirstName+user.LastName == order.Name || user.About == order.About {
					responseBody[i] = User{
						ID:     user.ID,
						Name:   user.FirstName + user.LastName,
						Age:    user.Age,
						About:  user.About,
						Gender: user.Gender,
					}
					i++
				}
			}
		}
	}

	if !(orderField == "Name" || orderField == "Age" || orderField == "Id") {
		w.WriteHeader(http.StatusBadRequest)
		error := SearchErrorResponse{
			Error: "OrderField invalid",
		}
		result, err2 := json.Marshal(error)
		if err2 != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(fmt.Errorf("%s", err2))
			return
		}
		_, err5 := w.Write(result)
		if err5 != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(fmt.Errorf("%s", err5))
			return
		}
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
		return
	}

	body, err := json.Marshal(responseBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(fmt.Errorf("%s", err))
		return
	}
	w.WriteHeader(http.StatusOK)

	fmt.Println(string(body))

	_, error := w.Write(body)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(fmt.Errorf("%s", err))
		return
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
		Query:      "{\"querylist\":[{\"Name\":\"GlennJordan\"},{\"Name\":\"RoseCarney\"},{\"Name\":\"OwenLynn\"}]}",
		OrderField: "sdafs",
		OrderBy:    1,
	}
	resp, err := client.FindUsers(request)
	if err != nil {
		fmt.Println(err)
		fmt.Println("resp", (resp))
		return
	}
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
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
