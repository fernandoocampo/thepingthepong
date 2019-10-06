package repository

import (
	"fmt"
	"os"

	"github.com/fernandoocampo/thepingthepong/common/logging"
	"github.com/fernandoocampo/thepingthepong/domain"
	"github.com/sirupsen/logrus"
)

var log *logging.Handle

// InitLog initializes log configuration for this module.
func InitLog(data domain.LogData) {
	var err error
	log, err = logging.NewLogger(
		logging.Options{
			LogLevel:  data.Level,
			LogFormat: data.Format,
			LogFields: logrus.Fields{"pkg": "repository", "srv": "thepingthepong"},
		})
	if err != nil {
		fmt.Printf("cant load repository logger: %v", err)
		os.Exit(1)
	}
}
