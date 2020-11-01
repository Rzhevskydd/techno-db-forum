package app

import (
	"database/sql"
	"fmt"
	"github.com/Rzhevskydd/techno-db-forum/project/app/units"
	userDelivery "github.com/Rzhevskydd/techno-db-forum/project/app/units/user/delivery"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

const ApiPrefix = "/api"

type Config struct {
	Port string
	Addr string
	DbHost string
	DbName string
	DbPort string
	DbUser string
	DbPwd string
}

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(cfg Config) {
	connectionString :=
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable ",
			cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPwd, cfg.DbName)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.DB.SetMaxOpenConns(100)
	a.DB.SetMaxIdleConns(30)
	a.DB.SetConnMaxLifetime(time.Hour)

	a.Router = mux.NewRouter().PathPrefix(ApiPrefix).Subrouter()

	a.initializeApplication()
}

func (a *App) Run(addr string) {
	defer func() {
		if err := a.DB.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeApplication() {
	repos := units.NewRepositories(a.DB)
	useCase := units.NewUseCase(repos)

	userRouter := a.Router.PathPrefix("/user").Subrouter()
	userDelivery.HandleUserRoutes(userRouter, useCase)

	//forumRouter := a.Router.PathPrefix("/forum").Subrouter()
	//forumDelivery.HandleForumRoutes(forumRouter, useCase)
	//userRouter := a.Router.PathPrefix("/user").Subrouter()
	//threadRouter := a.Router.PathPrefix("/thread").Subrouter()
	//postRouter := a.Router.PathPrefix("/post").Subrouter()
	//serviceRouter := a.Router.PathPrefix("/service").Subrouter()

	//a.Router.HandleFunc("/products", a.getProducts).Methods("GET")
	//a.Router.HandleFunc("/product", a.createProduct).Methods("POST")
	//a.Router.HandleFunc("/product/{id:[0-9]+}", a.getProduct).Methods("GET")
	//a.Router.HandleFunc("/product/{id:[0-9]+}", a.updateProduct).Methods("PUT")
	//a.Router.HandleFunc("/product/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")
}