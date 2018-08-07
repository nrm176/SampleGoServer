package main

import (
	"encoding/json"
		"github.com/gorilla/mux"
	"net/http"
	"log"
)

// our main function
//func main() {
//	router := mux.NewRouter()
//	log.Fatal(http.ListenAndServe(":8000", router))
//}

type Country struct {
	Name string              `json:"name"`
	Prefectures []Prefecture `json:"prefectures"`
}

type Prefecture struct {
	Name string    `json:"name"`
	Capital string `json:"capital"`
	Population int `json:"population"`
}

var country []Country


func GetCountry(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(country)
}

func main() {
	tokyo := Prefecture{Name:"東京都", Capital:"東京", Population:13482040}
	saitama := Prefecture{Name:"埼玉県", Capital:"さいたま市", Population:7249287}
	kanagawa := Prefecture{Name:"神奈川県", Capital:"横浜市", Population:9116252}
	//japan := Country{
	//	Name:"日本",
	//	Prefectures:[]Prefecture{tokyo, saitama, kanagawa},
	//}

	country = append(country, Country{
		Name:"日本",
		Prefectures:[]Prefecture{tokyo, saitama, kanagawa},
	})
	//jsonBytes, err := json.Marshal(country)
	//if err != nil {
	//	fmt.Println("JSON Marshal error:", err)
	//	return
	//}
	router := mux.NewRouter()
	router.HandleFunc("/country", GetCountry).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}