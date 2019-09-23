package playerapp

import (
	"fmt"
	"os"

	"github.com/fernandoocampo/thepingthepong/common/logging"
	"github.com/sirupsen/logrus"
)

var log *logging.Handle

func init() {
	var err error
	log, err = logging.NewLogger(logging.Options{LogLevel: "Info", LogFormat: "json", LogFields: logrus.Fields{"pkg": "playerapp", "srv": "thepingthepong"}})
	if err != nil {
		fmt.Printf("cant load logger: %v", err)
		os.Exit(1)
	}
}
