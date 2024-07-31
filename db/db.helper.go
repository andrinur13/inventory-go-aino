package db

import (
	"twc-ota-api/config"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmgorm/v2"
	_ "go.elastic.co/apm/module/apmgorm/v2/dialects/postgres"
)

var DB []*gorm.DB

func Init() {
	for i := range config.Databases {
		var database = &config.Databases[i]

		db, err := apmgorm.Open(database.DriverName, database.ConnectionString)
		if err != nil {
			fmt.Printf("%v : %v \n", err, database.ConnectionString)
			// i := 0
			for i := 0; i < 4; i++ {
				// i++
				db, err = apmgorm.Open(database.DriverName, database.ConnectionString)
				if err != nil {
					fmt.Printf("RECONECT(%d)%v : %v \n", i, err, database.ConnectionString)
					time.Sleep(3 * time.Second)
					continue
				}
				break
				// return reDB
				// panic("failed to connect database")
			}
			// return db
			if err != nil {
				panic("failed to connect database")
			}
			// panic("failed to connect database")
		}

		db.DB().SetMaxOpenConns(database.MaxConnectionOpen)
		if config.App.Env == "dev" || config.App.Env == "stg" || config.App.Env == "local" {
			db.LogMode(true)
		}

		log.WithFields(log.Fields{
			"config": database,
		}).Info("Connected to database")

		gormMigration(database.Name, db)
		//append database to array
		DB = append(DB, db)
	}
}

// register entity for created table
func gormMigration(dbName string, db *gorm.DB) {
	// if dbName == "gosample_db" {
	// db.AutoMigrate(&entities.Example{}, &entities.User{})
	// }
}
