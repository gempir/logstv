package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gocql/gocql"

	"github.com/gempir/logstv/common"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var cassandra *gocql.Session
var userCache = make(map[string]int64)

func main() {
	common.LoadEnv()

	var err error
	cassandra, err = common.NewDatabaseSession(strings.Split(common.GetEnv("DBHOSTS"), ","))
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.HideBanner = true

	DefaultCORSConfig := middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET},
	}
	e.Use(middleware.CORSWithConfig(DefaultCORSConfig))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to api.logs.tv!")
	})
	e.GET("/channel/:channel/user/:username", getUserLogs)

	fmt.Println("starting streamlogs API on port :8010")
	e.Logger.Fatal(e.Start(":8010"))
}
