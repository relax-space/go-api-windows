// Copyright 2015 Daniel Theophanes.
// Use of this source code is governed by a zlib-style
// license that can be found in the LICENSE file.

// simple does nothing except block while running the service.
package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/astaxie/beego/logs"
	"github.com/kardianos/service"
)

var logger *logs.BeeLogger

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	logger = logs.NewLogger()
	logger.SetLogger("file", `{"filename":"./logs/test.log"}`)

	e := echo.New()
	e.Use(middleware.CORS())

	e.GET("/ping", func(c echo.Context) error {
		log.Println("hello file1")
		logger.Error("hello file2")
		return c.String(http.StatusOK, "pong")
	})

	e.Logger.Fatal(e.Start(":8082"))

}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "go-windows-api",
		DisplayName: "go-windows-api",
		Description: "This is an example Go service.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	var logger service.Logger
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
