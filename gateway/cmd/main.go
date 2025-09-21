package main

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {

	log.WithField(
		"time", time.Now,
	).Info("Started api gateway")

	// INIT producer

}
