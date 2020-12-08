package config

import (
	"io/ioutil"

	"github.com/buger/jsonparser"
	log "github.com/sirupsen/logrus"
)

// appConf : struct for attributes needed for application config
type appConf struct {
	ServerPort string
	Env        string
	GtHost     string
}

// mailConf : struct for attributes needed for email config
type mailConf struct {
	Host     string
	Port     int
	Email    string
	Username string
	Password string
}

// Database : struct for attributes needed for Database config
type Database struct {
	Name              string
	DriverName        string
	ConnectionString  string
	MaxConnectionOpen int
}

// App : store appConfig
var App = appConf{}

// Mail : store mailConfig
var Mail = mailConf{}

// Databases : for storing database config
var Databases = []Database{}

// Init : init config
// mapping json config
// params : dev or prod or other config name
func Init(env string) {
	if env == "local" || env == "dev" || env == "stg" || env == "prod" || env == "sbx" {
		App.Env = env
		//read configuration file
		cfgBlob, errReadCfg := ioutil.ReadFile("./config/config." + env + ".json")
		if errReadCfg != nil {
			log.Fatal(errReadCfg)
		}

		var err error
		var serverport []byte
		serverport, _, _, err = jsonparser.Get(cfgBlob, "app", "serverport")
		if err != nil {
			log.Fatal(err)
		} else {
			App.ServerPort = string(serverport)
		}

		var gtHost []byte
		gtHost, _, _, err = jsonparser.Get(cfgBlob, "app", "gt_host")
		if err != nil {
			log.Fatal(err)
		} else {
			App.GtHost = string(gtHost)
		}

		var mailHost []byte
		mailHost, _, _, err = jsonparser.Get(cfgBlob, "email", "host")
		if err != nil {
			log.Fatal(err)
		} else {
			Mail.Host = string(mailHost)
		}

		var mailPort int64
		mailPort, err = jsonparser.GetInt(cfgBlob, "email", "port")
		if err != nil {
			log.Fatal(err)
		} else {
			Mail.Port = int(mailPort)
		}

		var mailEmail []byte
		mailEmail, _, _, err = jsonparser.Get(cfgBlob, "email", "email")
		if err != nil {
			log.Fatal(err)
		} else {
			Mail.Email = string(mailEmail)
		}

		var mailUsername []byte
		mailUsername, _, _, err = jsonparser.Get(cfgBlob, "email", "username")
		if err != nil {
			log.Fatal(err)
		} else {
			Mail.Username = string(mailUsername)
		}

		var mailPass []byte
		mailPass, _, _, err = jsonparser.Get(cfgBlob, "email", "password")
		if err != nil {
			log.Fatal(err)
		} else {
			Mail.Password = string(mailPass)
		}

		jsonparser.ArrayEach(cfgBlob, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			if err != nil {
				log.Fatal(err)
			}

			var dbname, dbdriver, dbconnstr []byte
			var dbmaxconnopen int64

			dbname, _, _, err = jsonparser.Get(value, "name")
			dbdriver, _, _, err = jsonparser.Get(value, "driver_name")
			dbconnstr, _, _, err = jsonparser.Get(value, "connection_string")
			dbmaxconnopen, err = jsonparser.GetInt(value, "max_connection_string")

			var database = Database{
				Name:              string(dbname),
				DriverName:        string(dbdriver),
				ConnectionString:  string(dbconnstr),
				MaxConnectionOpen: int(dbmaxconnopen)}

			Databases = append(Databases, database)
		}, "databases")

	} else {
		log.Fatal("cannot initialize config")
	}
}
