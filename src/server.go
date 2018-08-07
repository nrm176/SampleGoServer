package main

import (
	"database/sql"
	"log"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"reflect"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	)

var db *sql.DB

func load_Env() (string) {
	var db_config string
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db_config = os.Getenv("DATA_SOURCE")
	println(db_config)

	return db_config
}

func initConnection() {
	var db_config = load_Env()
	var err error
	db, err = sql.Open("postgres", db_config)

	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Successfully connected!")
}

func query(code string, dt string) ([]byte, error) {
	var objects []map[string]interface{}

	rows, _ := db.Query("select * from moving_average($1,$2, 25)", code, dt)

	for rows.Next() {
		columns, err := rows.ColumnTypes()
		if err != nil {
			return nil, err
		}

		values := make([]interface{}, len(columns))
		object := map[string]interface{}{}
		for i, column := range columns {
			object[column.Name()] = reflect.New(column.ScanType()).Interface()
			values[i] = object[column.Name()]
		}

		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		objects = append(objects, object)
	}

	return json.MarshalIndent(objects, "", "\t")
}

func GetSMA(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]

	date := r.FormValue("date")

	log.Printf("quering %s on %s", code, date)
	res, err := query(code, date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
func main() {
	initConnection()
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/sma/{code:[0-9]+}", GetSMA).Queries("date", "{date}").Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))

}
