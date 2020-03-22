package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/ant0ine/go-json-rest/rest"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"google.golang.org/api/option"
)

func main() {
	i := Impl{}
	i.InitDB()
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Post("/daily", i.PostDaily),
		rest.Get("/healthcheck", i.GetHealthCheck),
		rest.Post("/auth", i.PostAuth),
		rest.Post("/create_user", i.CreateUser),
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
	UserID    int64     `json:"userId"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateUserParameter struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Profile        string `json:"profile"`
	ProfileImgPath string `json:"profileImgPath"`
	HeaderImgPath  string `json:"headerImaPath"`
	PassWord       string `json:"passWord"`
}

type User struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Profile        string `json:"profile"`
	ProfileImgPath string `json:"profileImgPath"`
	HeaderImgPath  string `json:"headerImaPath"`
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
	reminder := Article{}
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

func (i *Impl) GetHealthCheck(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(map[string]string{"Body": "OK"})
}

func (i *Impl) CreateUser(w rest.ResponseWriter, r *rest.Request) {
	// Firebase SDK のセットアップ
	opt := option.WithCredentialsFile("./envfiles/admin_sdk_firebase.json")
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	createUserParameter := CreateUserParameter{}

	if err := r.DecodeJsonPayload(&createUserParameter); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := (&auth.UserToCreate{}).
		Email(createUserParameter.Email).
		EmailVerified(false).
		Password(createUserParameter.PassWord).
		DisplayName(createUserParameter.Name).
		Disabled(false)
	u, err := client.CreateUser(ctx, params)

	dbparams := User{
		Name:           createUserParameter.Name,
		Email:          createUserParameter.Email,
		Profile:        createUserParameter.Profile,
		ProfileImgPath: createUserParameter.ProfileImgPath,
		HeaderImgPath:  createUserParameter.HeaderImgPath}

	if err != nil {
		log.Fatalf("error creating user: %v\n", err)
	}
	log.Printf("Successfully created user: %v\n", u)

	if err := i.DB.Save(dbparams).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(map[string]string{"Body": "Created"})
}

func (i *Impl) PostAuth(w rest.ResponseWriter, r *rest.Request) {
	// Firebase SDK のセットアップ
	opt := option.WithCredentialsFile(os.Getenv("FIREBASE_ADMIN_SDK_FILENAME"))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	auth, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	// クライアントから送られてきた JWT 取得
	authHeader := r.Header.Get("Authorization")
	idToken := strings.Replace(authHeader, "Bearer ", "", 1)

	// JWT の検証
	token, err := auth.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		// JWT が無効なら Handler に進まず別処理
		log.Printf("error verifying ID token: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
	}
	log.Printf("Verified ID token: %v\n", token)
}
