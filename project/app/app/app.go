package app

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Rzhevskydd/techno-db-forum/project/app/units"
	forumDelivery "github.com/Rzhevskydd/techno-db-forum/project/app/units/forum/delivery"
	postDelivery "github.com/Rzhevskydd/techno-db-forum/project/app/units/post/delivery"
	serviceDelivery "github.com/Rzhevskydd/techno-db-forum/project/app/units/service/delivery"
	threadDelivery "github.com/Rzhevskydd/techno-db-forum/project/app/units/thread/delivery"
	userDelivery "github.com/Rzhevskydd/techno-db-forum/project/app/units/user/delivery"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

func LogRequestsMiddleware(h http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout,h)
}

func (a *App) Initialize(cfg Config) {
	var err error
	err = a.initializeDatabase(cfg)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter().PathPrefix(ApiPrefix).Subrouter()
	a.Router.Use(LogRequestsMiddleware)

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

	forumRouter := a.Router.PathPrefix("/forum").Subrouter()
	forumDelivery.HandleForumRoutes(forumRouter, useCase)

	threadRouter := a.Router.PathPrefix("/thread").Subrouter()
	threadDelivery.HandleThreadRoutes(threadRouter, useCase)

	postRouter := a.Router.PathPrefix("/post").Subrouter()
	postDelivery.HandlePostRoutes(postRouter, useCase)

	serviceRouter := a.Router.PathPrefix("/service").Subrouter()
	serviceDelivery.HandleServiceRoutes(serviceRouter, useCase)
}

func (a *App) initializeDatabase(cfg Config) (err error) {
	if a.DB != nil {
		return errors.New("db already initialized")
	}

	connectionString :=
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable ",
			cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPwd, cfg.DbName)
	a.DB, err = sql.Open("postgres", connectionString)

	a.DB.SetMaxOpenConns(100)
	a.DB.SetMaxIdleConns(30)
	a.DB.SetConnMaxLifetime(time.Hour)

	if query, err := ioutil.ReadFile("db/db_init_tables.sql"); err != nil {
		return err
	} else {
		_, err = a.DB.Exec(string(query))
		return err
	}

}
