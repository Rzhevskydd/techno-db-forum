package main

import (
	"github.com/Rzhevskydd/techno-db-forum/project/app/app"
)

func main() {
	// TODO аргументы командной строки
	cfg := app.Config{
		Port:   "",
		Addr:   "",
		DbHost: "",
		DbName: "",
		DbPort: "",
		DbUser: "",
		DbPwd:  "",
	}

	a := app.App{}
	a.Initialize(cfg)

	a.Run(":" + cfg.Port)
}