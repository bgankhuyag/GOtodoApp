package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    _ "github.com/jinzhu/gorm/dialects/mysql"
)

type ToDo struct {
    ID   int `json:"id"`
    Item string `json:"item"`
    CreatedAt   string `json:"created_at"`
    Completed int `json:"completed"`
}

type JsonResponse struct {
    Type    string `json:"type"`
    Data    []ToDo `json:"data"`
    Message string `json:"message"`
}

type BaseHandler struct {
  db *sql.DB
}

// NewBaseHandler returns a new BaseHandler
func NewBaseHandler(db *sql.DB) *BaseHandler {
	return &BaseHandler{
		db: db,
	}
}

// Function for handling errors
func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

func setupDB() *sql.DB {
    db, err := sql.Open("mysql", "sql5447318:Ia1zrrKczy@tcp(sql5.freesqldatabase.com:3306)/sql5447318")
    checkErr(err)
    return db
}

func (h *BaseHandler) getItems(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!")
    // db := setupDB()
    results, err := h.db.Query("select * from Applications")
  	checkErr(err)
  	defer results.Close()
    var list []ToDo
  	for results.Next() {
  		var (
  			id    int
  			item  string
  			created_at string
        completed int
  		)
  		err = results.Scan(&id, &item, &created_at, &completed)
  		checkErr(err)
      list = append(list, ToDo{ID: id, Item: item, CreatedAt: created_at, Completed: completed})
  	}
    var response = JsonResponse{Type: "success", Data: list}
    json.NewEncoder(w).Encode(response)
}

func main() {
  db := setupDB()
  h := NewBaseHandler(db)
  r := mux.NewRouter()
  r.HandleFunc("/getItems", h.getItems).Methods("GET")
  log.Fatal(http.ListenAndServe(":8080", r))
  defer db.Close()

	// var (
	// 	id    int
	// 	name  string
	// 	price int
	// )
	// err = db.QueryRow("Select * from product where id = 1").Scan(&id, &name, &price)
	// if err != nil {
	// 	log.Fatal("Unable to parse row:", err)
	// }
	// fmt.Printf("ID: %d, Name: '%s', Price: %d\n", id, name, price)
	// products := []struct {
	// 	name  string
	// 	price int
	// }{
	// 	{"Light", 10},
	// 	{"Mic", 30},
	// 	{"Router", 90},
	// }
	// stmt, err := db.Prepare("INSERT INTO product (name, price) VALUES (?, ?)")
	// defer stmt.Close()
	// if err != nil {
	// 	log.Fatal("Unable to prepare statement:", err)
	// }
	// for _, product := range products {
	// 	_, err = stmt.Exec(product.name, product.price)
	// 	if err != nil {
	// 		log.Fatal("Unable to execute statement:", err)
	// 	}
	// }
}
