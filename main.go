package main

import (
	"fmt"
	"log"
	"net/http"

	"treatment-scheme-organizer/config"
	"treatment-scheme-organizer/handler"
	"treatment-scheme-organizer/module"
	"treatment-scheme-organizer/route"
	"treatment-scheme-organizer/store"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

func main() {
	server := echo.New()

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4",
		config.Env.DB.User,
		config.Env.DB.Password,
		config.Env.DB.Host,
		config.Env.DB.Port,
		config.Env.DB.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to the Database", err)
	}

	if err := db.AutoMigrate(
		&module.Drug{},
	); err != nil {
		log.Fatal("Failed to migrate models' changes to Database", err)
	}

	apiRouterGroup := server.Group("/api")

	apiRouterGroup.GET("/ping", func(context echo.Context) error {
		return context.JSON(http.StatusOK, echo.Map{"status": true, "message": "pong"})
	})

	ds := store.NewDrugStore(db)
	dh := handler.NewDrugHandler(ds)
	dr := route.NewDrugRouter(dh)

	dr.AddRoutes(apiRouterGroup)

	server.Logger.Fatal(server.Start(fmt.Sprintf(":%d", config.Env.Server.Port)))
}
