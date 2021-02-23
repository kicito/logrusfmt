package main

import (
	"fmt"
	"github.com/annopkomol/logrusfmt"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetFormatter(logrusfmt.ProductionFormatter)
	log.ExitFunc = func(i int) {
		//do nothing
	}
	//without stack trace
	logTest(log)
	//with stack tracing
	fmt.Println()
	log.SetReportCaller(true)
	logTest(log)
}

func logTest(log *logrus.Logger) {
	log.WithField("hello", "world").WithField("foo", "bar").Trace("trace")
	log.Debug("debug")
	log.Info("info")
	log.WithField("hello", "world").Warn("warning")
	log.Error("error")
	log.Fatal("fatal")
}
