package db

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
	"twc-ota-api/config"
	"twc-ota-api/logger"

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

const (
	maxRetries    = 10
	retryInterval = 3000 * time.Millisecond
)

func WithRetry(model string, req string, db *gorm.DB, query func(db *gorm.DB) error) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = query(db)
		if err == nil {
			logger.Info(fmt.Sprintf("Query for %s executed successfully on attempt %d...", model, i+1), "200", true, "", req)
			return nil
		}

		if errors.Is(err, gorm.ErrInvalidTransaction) || !isRetryableError(err) {
			logger.Error(fmt.Sprintf("Non-retryable error occurred for %s on attempt %d...", model, i+1), "500", false, req, err)
			break
		}

		if strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint") {
			logger.Info(fmt.Sprintf("Duplicate key error encountered, skipping retry for %s", model), "200", false, "", req)
			return err
		}

		logger.Error(fmt.Sprintf("Retryable for %s error occurred on attempt %d...", model, i+1), "500", false, req, err)
		time.Sleep(retryInterval)
	}
	return err
}

func WithTransactionRetry(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		tx := db.Begin()
		if tx.Error != nil {
			logger.Error("Error starting transaction", "500", false, "", tx.Error)
			return tx.Error
		}

		err = fn(tx)
		if err == nil {
			if commitErr := tx.Commit().Error; commitErr == nil {
				logger.Info(fmt.Sprintf("Transaction committed successfully on attempt %d...", i+1), "200", false, "", "")
				return nil
			} else {
				err = commitErr
			}
		} else {
			tx.Rollback()
			logger.Error(fmt.Sprintf("Transaction rolled back on attempt %d: %s", i+1, err), "500", false, "", err)
		}

		if errors.Is(err, gorm.ErrInvalidTransaction) || !isRetryableError(err) {
			logger.Error(fmt.Sprintf("Non-retryable error occurred on attempt %d...", i+1), "500", false, "", err)
			break
		}

		if strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint") {
			logger.Info("Duplicate key error encountered, skipping retry", "200", false, "", "")
			return err
		}

		logger.Error(fmt.Sprintf("Retryable error occurred on attempt %d: %s", i+1, err), "500", false, "", err)
		reconnectDatabase()
		time.Sleep(retryInterval)
	}
	return err
}

func isRetryableError(err error) bool {
	return reflect.TypeOf(err).String() == "*net.OpError" || err.Error() == "driver: bad connection" || strings.Contains(err.Error(), "connection refused")
}

func reconnectDatabase() {
	logger.Info("Reconnecting to the database...", "200", false, "", "")
	Init()
}
