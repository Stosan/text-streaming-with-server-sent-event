package main

import (
	"fmt"

	"streamer/route"

"github.com/labstack/echo"
"github.com/labstack/echo/middleware"
)
func main() {

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/stream", route.StreamText)

	fmt.Println("Server listening on :8080")
	e.Logger.Fatal(e.Start(":8080"))
}

