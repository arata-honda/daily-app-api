package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
    "log"
	"net/http"
	"time"
)

func main() {
	i := Impl{}
	i.InitDB()
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Post("/daily", i.PostDaily),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}

type Impl struct {
	DB *gorm.DB
}

type Article struct {
    ID                int64  `json:"id"`
    UserID            int64  `json:"userId"`
    Title             string `json:"title"`
    Content           string `json:"content"`
	CreatedAt		  time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}


func (i *Impl) InitDB() {
    var err error
    i.DB, err = gorm.Open("mysql", "daily:daily@(mysql)/daily?charset=utf8&parseTime=True&loc=Local")
    if err != nil {
        log.Fatalf("Got error when connect database, the error is '%v'", err)
    }
    i.DB.LogMode(true)
}

func (i *Impl) PostDaily(w rest.ResponseWriter, r *rest.Request) {
	reminder := Articles{}
	if err := r.DecodeJsonPayload(&reminder); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := i.DB.Save(&reminder).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&reminder)
}