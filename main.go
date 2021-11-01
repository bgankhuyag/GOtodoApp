package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "time"
    // "strconv"
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

func enableCors(w *http.ResponseWriter) {
  (*w).Header().Set("Access-Control-Allow-Origin", "*")
  (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
  (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func (h *BaseHandler) getItems(w http.ResponseWriter, r *http.Request){
    enableCors(&w)
    results, err := h.db.Query("select * from Applications ORDER BY created_at")
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
  enableCors(&w)
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
    response = JsonResponse{Type: "success", Data: []ToDo{ToDo{ID: int(lastId), Item: item, CreatedAt: currentTime, Completed: 0}}, Message: "The item has been inserted successfully!"}
  }
  json.NewEncoder(w).Encode(response)
}

func (h *BaseHandler) updateItem(w http.ResponseWriter, r *http.Request) {
  enableCors(&w)
  item := r.FormValue("item")
  params := mux.Vars(r)
  id := params["id"]
  var response JsonResponse
  if item == "" || id == "" {
    fmt.Printf("ID: %s", item)
    response = JsonResponse{Type: "error", Message: "ID or Item is missing"}
  } else {
    result, err := h.db.Exec("UPDATE `Applications` SET `item`=? WHERE `id`=?", item, id)
    checkErr(err)
    affected, err := result.RowsAffected()
    checkErr(err)
    fmt.Printf("Affected items: %d\n", affected)
    if (affected == 0) {
      response = JsonResponse{Type: "error", Message: "Item update unsuccessful!"}
    } else {
      response = JsonResponse{Type: "success", Message: "The item has been updated successfully!"}
    }
  }
  json.NewEncoder(w).Encode(response)
}

func (h *BaseHandler) deleteItem(w http.ResponseWriter, r *http.Request) {
  enableCors(&w)
  params := mux.Vars(r)
  id := params["id"]
  var response JsonResponse
  if id == "" {
    response = JsonResponse{Type: "error", Message: "ID is missing"}
  } else {
    result, err := h.db.Exec("DELETE FROM `Applications` WHERE `id`=?", id)
    checkErr(err)
    affected, err := result.RowsAffected()
    checkErr(err)
    fmt.Printf("Affected items: %d\n", affected)
    if (affected == 0) {
      response = JsonResponse{Type: "error", Message: "Item delete unsuccessful!"}
    } else {
      response = JsonResponse{Type: "success", Message: "The item has been deleted successfully!"}
    }
  }
  json.NewEncoder(w).Encode(response)
}

func (h *BaseHandler) updateCompleted(w http.ResponseWriter, r *http.Request) {
  enableCors(&w)
  completed := r.FormValue("completed")
  params := mux.Vars(r)
  id := params["id"]
  var response JsonResponse
  if id == "" {
    response = JsonResponse{Type: "error", Message: "ID is missing"}
  } else {
    result, err := h.db.Exec("UPDATE `Applications` SET `completed`=? WHERE `id`=?", completed, id)
    checkErr(err)
    affected, err := result.RowsAffected()
    checkErr(err)
    fmt.Printf("Affected items: %d\n", affected)
    response = JsonResponse{Type: "success", Message: "The item has been updated successfully!"}
  }
  json.NewEncoder(w).Encode(response)
}

func main() {
  db := setupDB()
  h := NewBaseHandler(db)
  r := mux.NewRouter()
  r.HandleFunc("/get_items", h.getItems).Methods("GET", "OPTIONS")
  r.HandleFunc("/new_item", h.newItem).Methods("POST", "OPTIONS")
  r.HandleFunc("/update_item/{id}", h.updateItem).Methods("PUT", "OPTIONS")
  r.HandleFunc("/update_completed/{id}", h.updateCompleted).Methods("PUT", "OPTIONS")
  r.HandleFunc("/delete_item/{id}", h.deleteItem).Methods("DELETE", "OPTIONS")
  log.Fatal(http.ListenAndServe(":8080", r))
  defer db.Close()
}
