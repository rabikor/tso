package main

import (
	"fmt"
	"net/http"

	"treatment-scheme-organizer/config"
	"treatment-scheme-organizer/database"
	"treatment-scheme-organizer/database/migration"
	"treatment-scheme-organizer/server"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	server.NewErrorHandler(e)
	server.NewRouter(e)

	env, err := config.NewEnv("./.env")
	if err != nil {
		e.Logger.Fatal(err)
	}

	db, err := database.Open(env)
	if err != nil {
		e.Logger.Fatal(err)
	}

	if err := migration.Migrate(db); err != nil {
		e.Logger.Fatal(err)
	}

	defer func(db *database.DB) {
		err := db.Close()
		if err != nil {
			e.Logger.Fatal("Connection to mysql was not closed.")
		}
	}(db)

	rtr := e.Group("/api")

	rtr.GET("/ping", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, echo.Map{"status": true, "message": "pong"})
	})

	drugs := server.NewDrugsHandler(env, db)
	drugs.AddRoutes(rtr)

	illnesses := server.NewIllnessesHandler(env, db)
	illnesses.AddRoutes(rtr)

	procedures := server.NewProceduresHandler(env, db)
	procedures.AddRoutes(rtr)

	schemes := server.NewSchemesHandler(env, db)
	schemes.AddRoutes(rtr)

	schemeDays := server.NewSchemeDaysHandler(env, db)
	schemeDays.AddRoutes(rtr)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", env.Server.Port)))
}
