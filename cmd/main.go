package main

import (
	"chat_server/handler"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	ngHandler := handler.NewNoticeHandler()
	e.GET("/ws", ngHandler.WsHandler)
	if err := e.Start(":8888"); err != nil {
		fmt.Println(err)
	}
}
