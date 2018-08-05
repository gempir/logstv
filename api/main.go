// package main

// import (
// 	"fmt"
// 	"net/http"
// )

// func main() {
// 	e := echo.New()
// 	e.HideBanner = true

// 	DefaultCORSConfig := middleware.CORSConfig{
// 		Skipper:      middleware.DefaultSkipper,
// 		AllowOrigins: []string{"*"},
// 		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
// 	}
// 	e.Use(middleware.CORSWithConfig(DefaultCORSConfig))

// 	e.GET("/", func(c echo.Context) error {
// 		return c.String(http.StatusOK, "Welcome to streamlogs!")
// 	})
// 	e.GET("/channel/:channel/user/:username", s.getCurrentUserLogs)

// 	fmt.Println("starting streamlogs API on port :8010")
// 	e.Logger.Fatal(e.Start(":8010"))
// }
