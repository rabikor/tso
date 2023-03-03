package main

import (
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
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

	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			e.Logger.Fatal("Connection to mysql was not closed.")
		}
	}(db)

	rtr := e.Group("/api")

	rtr.GET("/ping", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, echo.Map{"status": true, "message": "pong"})
	})

	p := server.NewPagination(env)

	dr := database.NewDrugsRepository(db)
	drugs := server.NewDrugsHandler(p, dr)
	drugs.AddRoutes(rtr)

	ir := database.NewIllnessesRepository(db)
	illnesses := server.NewIllnessesHandler(env, ir)
	illnesses.AddRoutes(rtr)

	pr := database.NewProceduresRepository(db)
	procedures := server.NewProceduresHandler(env, pr)
	procedures.AddRoutes(rtr)

	sr := database.NewSchemesRepository(db)
	schemes := server.NewSchemesHandler(env, sr)
	schemes.AddRoutes(rtr)

	sdr := database.NewSchemeDaysRepository(db)
	schemeDays := server.NewSchemeDaysHandler(env, sdr)
	schemeDays.AddRoutes(rtr)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", env.Server.Port)))
}
