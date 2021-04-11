package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"

	// _ infront of import means it is merely injecting side-effects
	// In layman's terms basically mean create package-level variable and execute
	// init function of the package
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db, _ = gorm.Open("mysql", "root:root@/ajadb?charset=utf8&parseTime=True&loc=Local")

type ItemActivityModel struct {
	Id         int       `gorm:"primary_key"`
	Date       time.Time `json:"done_date"`
	Order      int
	TodoItemId int `gorm:"index:todo_idx"`
}

type TodoItemModel struct {
	Id          int `gorm:"primary_key"`
	Description string
	Completed   bool
	DoneDates   []ItemActivityModel
}

func getItemByID(Id int) bool {
	todo := &TodoItemModel{}
	result := db.First(&todo, Id)
	if result.Error != nil {
		log.Warn("Activity not found in the database")
		return false
	}
	return true
}

func addDoneDate(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	err := getItemByID(id)
	if err == false {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"updated": false, "error": "Record not found"}`)
	} else {
		doneDate := r.FormValue("done_date")
		log.WithFields(log.Fields{"id": id, "Done Date": doneDate}).Info("Updating activity")
		t, _ := time.Parse("2006-01-02", doneDate)
		todo := &TodoItemModel{}
		db.First(&todo, id)
		activity := &ItemActivityModel{Date: t, Order: id, TodoItemId: id}
		db.Create(&activity)
		log.WithFields(log.Fields{"id": todo.Id, "Description": todo.Description}).Info("Updating activity")
		todo.DoneDates = append(todo.DoneDates, *activity)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"updated": true}`)
	}
}

func getDoneActivities(w http.ResponseWriter, r *http.Request) {
	// This method returns the list of done activities within a date range
	t1 := r.FormValue("start_date")
	t2 := r.FormValue("end_date")
	startDate, _ := time.Parse("2006-01-02", t1)
	endDate, _ := time.Parse("2006-01-02", t2)
	activities := []ItemActivityModel{}
	log.WithFields(log.Fields{"start_date": startDate, "end_date": endDate}).Info("found activity")
	db.Where("Date >= ? and Date <= ?", startDate, endDate).Find(&activities)
	// db.Find(&activities)
	response := []string{}
	for _, activity := range activities {
		todo := &TodoItemModel{}
		db.First(&todo, activity.TodoItemId)
		response = append(response, todo.Description)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	err := getItemByID(id)
	if err == false {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"updated": false, "error": "Record not found"}`)
	} else {
		completed, _ := strconv.ParseBool(r.FormValue("completed"))
		log.WithFields(log.Fields{"id": id, "Completed": completed}).Info("Updating activity")
		todo := &TodoItemModel{}
		db.First(&todo, id)
		todo.Completed = completed
		db.Save(&todo)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"updated": true}`)
	}
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	err := getItemByID(id)
	if err == false {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"deleted": false, "error": "record not found"}`)
	} else {
		log.WithFields(log.Fields{"id": id}).Info("Deleting the item")
		todo := &TodoItemModel{}
		db.First(&todo, id)
		db.Delete(&todo)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"deleted": true}`)
	}
}

func getTodoItems(completed bool) interface{} {
	var todos []TodoItemModel
	TodoItems := db.Where("completed = ?", completed).Find(&todos).Value
	return TodoItems
}

func getCompletedItems(w http.ResponseWriter, r *http.Request) {
	log.Info("Get completed activities")
	completedTodoItems := getTodoItems(true)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(completedTodoItems)
}

func getIncompleteItems(w http.ResponseWriter, r *http.Request) {
	log.Info("Get incomplete activities")
	incompleteTodoItems := getTodoItems(false)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(incompleteTodoItems)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	description := r.FormValue("description")
	log.WithFields(log.Fields{"description": description}).Info("Add new activity, saving to db")
	todo := &TodoItemModel{Description: description, Completed: false}
	db.Create(&todo)
	result := db.Last(&todo)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result.Value)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	log.Info("API Health is OK")
	w.Header().Set("Content-Type", "application/json")
	// `` declares raw strings inside which it is legal to have any characters
	io.WriteString(w, `{"alive": true}`)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/healthz", healthz).Methods("GET")
	r.HandleFunc("/app", handler).Methods("GET")
	r.HandleFunc("/todo", createItem).Methods("POST")
	r.HandleFunc("/todo-complete", getCompletedItems).Methods("Get")
	r.HandleFunc("/todo-incomplete", getIncompleteItems).Methods("Get")
	r.HandleFunc("/todo/{id}", updateItem).Methods("POST")
	r.HandleFunc("/todo/{id}", deleteItem).Methods("DELETE")
	r.HandleFunc("/todo-activity/{id}", addDoneDate).Methods("POST")
	r.HandleFunc("/done-activities", getDoneActivities).Methods("POST")
	return r

}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

func main() {
	defer db.Close()

	// db.Debug().DropTableIfExists(&TodoItemModel{}) // I doubt if this is a good practice
	db.Debug().AutoMigrate(&TodoItemModel{})

	// db.Debug().DropTableIfExists(&ItemActivityModel{}) // I doubt if this is a good practice
	db.Debug().AutoMigrate(&ItemActivityModel{})
	log.Info("Starting AJA API server")
	r := newRouter()
	handler := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "DELETE", "PATCH", "OPTIONS"},
	}).Handler(r)
	http.ListenAndServe(":8080", handler)
}
