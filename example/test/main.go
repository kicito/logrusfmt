package main

import (
	"github.com/annopkomol/logrusfmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	var (
		e      = echo.New()
		logger = logrus.New()
	)

	e.Use(logrusfmt.AddRequestCtxMiddleware)
	e.Use(logrusfmt.LoggingMiddleware(logger))
	logger.SetLevel(logrus.TraceLevel)
	logger.SetReportCaller(true)
	logger.SetFormatter(logrusfmt.ProductionFormatter)
	logger.Info("info text")
	logger.Warn("warning")
	logger.Debug("debug text")

	e.GET("/", func(c echo.Context) error {
		logger.WithContext(c.Request().Context()).WithField("field1", "value1").Error("error text")
		return c.String(http.StatusOK, "ok")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
