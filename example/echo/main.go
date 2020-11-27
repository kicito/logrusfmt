package main

import (
	"context"
	"github.com/annopkomol/logrusfmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	var (
		env = "production"
		e   = echo.New()
		log = logrus.New()
	)
	//enable stack tracing
	log.SetReportCaller(true)

	//set formatter
	switch env {
	case "local":
		log.SetLevel(logrus.TraceLevel)
		log.SetFormatter(logrusfmt.LocalFormatter)
	case "staging":
	case "production":
		log.SetFormatter(logrusfmt.ProductionFormatter)
	}

	//add middleware
	e.Use(logrusfmt.AddRequestCtxMiddleware)
	e.Use(logrusfmt.LoggingMiddleware(log))

	e.GET("/", func(c echo.Context) error {
		var ctx context.Context = c.Request().Context()
		//
		//
		//some biz logic
		//
		//
		//log error
		log.WithContext(ctx).
			WithFields(logrus.Fields{
				"from":   "account A",
				"amount": 500,
			}).Error("couldn't transfer to account B")

		return c.String(http.StatusOK, "ok")
	})
	e.Logger.Fatal(e.Start(":80"))
}
