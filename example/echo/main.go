package main

import (
	"context"
	"net/http"

	"github.com/annopkomol/logrusfmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func main() {
	var (
		env = "production"
		e   = echo.New()
		log = logrus.New()
	)

	echo.NotFoundHandler = func(c echo.Context) error {
		// render your 404 page
		return c.String(http.StatusNotFound, "not found page")
	}
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
	e.Use(echo.WrapMiddleware(logrusfmt.RequestHTTPMiddleware))
	e.Use(echo.WrapMiddleware(logrusfmt.LoggingHTTPMiddleware(log)))
	//recover must be place after logging
	e.Use(middleware.Recover())

	e.GET("/test", func(c echo.Context) error {
		var ctx context.Context = c.Request().Context()
		//
		//
		//some biz logic
		//
		//
		//log error
		//panic("hoho")
		log.WithContext(ctx).
			WithFields(logrus.Fields{
				"from":   "account A",
				"amount": 500,
			}).Error("couldn't transfer to account B")

		return c.String(http.StatusOK, "ok")
	})
	e.Logger.Fatal(e.Start(":80"))
}
