package main

import (
	"assignment-2/config"
	"assignment-2/controller"

	"github.com/labstack/echo/v4"
)

func init() {
	// start connect DB at init
	config.StartDB()
}

func main() {
	e := echo.New()

	e.POST("/order", controller.Create)
	e.PUT("/order", controller.Update)
	e.GET("/order/:id", controller.Get)
	e.DELETE("/order/:id", controller.Delete)

	e.Logger.Fatal(e.Start(":9090"))
}
