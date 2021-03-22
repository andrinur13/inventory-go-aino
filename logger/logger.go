package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	_ "twc-ota-api/docs"
)

var version = "1.6.1"

//Logger : logging function
func Logger() {
	// path := "logs"
	// if _, err := os.Stat(path); os.IsNotExist(err) {
	// 	os.Mkdir(path, 0777)
	// }

	// if err := os.Chmod(path, 0777); err != nil {
	// 	fmt.Println(err)
	// }

	var filename = "logs/" + time.Now().Format("2006-01-02") + ".log"
	// Create the log file if doesn't exist. And append to it if it already exists.
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	// err := os.Chmod(filename, 0777)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	Formatter := new(logrus.TextFormatter)
	// You can change the Timestamp format. But you have to use the same date and time.
	// "2006-02-02 15:04:06" Works. If you change any digit, it won't work
	// ie "Mon Jan 2 15:04:05 MST 2006" is the reference time. You can't change it
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	logrus.SetFormatter(Formatter)
	if err != nil {
		// Cannot open log file. Logging to stderr
		fmt.Println(err)
	} else {
		logrus.SetOutput(f)
	}

}

func Info(message string, code string, status bool, resp string, req string) (response *string) {
	Logger()
	logrus.WithFields(logrus.Fields{
		"request":  req,
		"response": resp,
		"success":  status,
		"code":     code,
		"version":  version,
	}).Info(message)

	return nil
}

func Warning(message string, code string, status bool, req string) (response *string) {
	Logger()
	logrus.WithFields(logrus.Fields{
		"request": req,
		"success": status,
		"code":    code,
		"version": version,
	}).Warning(message)

	return nil
}
