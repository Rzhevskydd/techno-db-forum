package main

import (
	"github.com/Rzhevskydd/techno-db-forum/project/app/app"
)

func main() {
	// TODO аргументы командной строки
	cfg := app.Config{
		Port:   "5000",
		Addr:   "",
		DbHost: "localhost",
		//DbName: "forum_func",
		DbName: "forum",
		DbPort: "5432",
		DbUser: "forum",
		//DbUser: "danil",
		DbPwd:  "forum",
		//DbPwd:  "password",
	}

	a := app.App{}
	a.Initialize(cfg)

	a.Run(":" + cfg.Port)
}