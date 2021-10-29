package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "time"
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

func (h *BaseHandler) newItem(w http.ResponseWriter, r *http.Request) {
  item := r.FormValue("item")
  var response JsonResponse
  if item == "" {
    response = JsonResponse{Type: "error", Message: "Missing item!"}
  } else {
    currentTime := time.Now().Format("2006.01.02 15:04:05")
    result, err := h.db.Exec("INSERT INTO `Applications` (`item`, `created_at`, `completed`) VALUES (?, ?, ?)", item, currentTime, 0)
    checkErr(err)
    lastId, err := result.LastInsertId()
    checkErr(err)
    fmt.Printf("The last inserted item id: %d\n", lastId)
    response = JsonResponse{Type: "success", Message: "The item has been inserted successfully!"}
  }
  json.NewEncoder(w).Encode(response)
}

func (h *BaseHandler) updateItem(w http.ResponseWriter, r *http.Request) {
  item := r.FormValue("item")
  params := mux.Vars(r)
  id := params["id"]
  var response JsonResponse
  if item == "" || id == "" {
    response = JsonResponse{Type: "error", Message: "ID or item missing"}
  } else {
    result, err := h.db.Exec("UPDATE `Applications` SET `item`=? WHERE `id`=?", item, id)
    checkErr(err)
    affectedID, err := result.RowsAffected()
    checkErr(err)
    fmt.Printf("Affected item id: %d\n", affectedID)
    response = JsonResponse{Type: "success", Message: "The item has been updated successfully!"}
  }
  json.NewEncoder(w).Encode(response)
}

func main() {
  db := setupDB()
  h := NewBaseHandler(db)
  r := mux.NewRouter()
  r.HandleFunc("/get_items", h.getItems).Methods("GET")
  r.HandleFunc("/new_item", h.newItem).Methods("POST")
  r.HandleFunc("/update_item/{id}", h.updateItem).Methods("PUT")
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
