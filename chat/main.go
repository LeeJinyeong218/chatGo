package main

import (
	"chat/dto"
	"chat/router"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

// InitConfig -> config
func InitConfig() (*dto.Config, error) {
	appEnv := flag.String("env", os.Getenv("API_ENV"), "API Environment")
	flag.Parse()

	var config dto.Config
	configType, configName, configPath := "yaml", "config", "."

	if *appEnv == "" {
		*appEnv = "production"
	}

	viper.SetConfigType(configType)
	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey(*appEnv, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// InitDatabase -> database engine 만들기
func InitDatabase(dbConfig dto.ConfigDatabase) (*gorm.DB, error) {
	// account from env
	dbUsername := os.Getenv("DB_USERNAME")
	if dbUsername == "" {
		return nil, fmt.Errorf("DB Username ERROR")
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		return nil, fmt.Errorf("DB Password ERROR")
	}
	dbAccount := fmt.Sprintf("%s:%s", dbUsername, dbPassword)

	// log
	dbLogFile, err := os.Create(dbConfig.Log)
	if err != nil {
		return nil, err
	}

	newLogger := logger.New(
		log.New(dbLogFile, "\r\n", log.Ldate|log.Ltime|log.Lshortfile), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	// create db
	dbUri := fmt.Sprintf("%s@tcp(%s:%d)/%s", dbAccount, dbConfig.Host, dbConfig.Port, dbConfig.Name)
	dbEngine, err := gorm.Open(mysql.Open(dbUri), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	return dbEngine, nil
}

// InitLogger -> log 만들기
func InitLogger(serverConfig dto.ConfigServer) (*logrus.Entry, error) {
	newLogger := logrus.New()
	logLevel := logrus.ErrorLevel
	if serverConfig.Debug {
		logLevel = logrus.DebugLevel
	}
	newLogger.SetLevel(logLevel)
	file, err := os.OpenFile(serverConfig.Log, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to log to file, using default stderr")
	} else {
		newLogger.SetOutput(file)
	}
	newLogger.SetFormatter(&logrus.TextFormatter{
		DisableColors:   false,
		TimestampFormat: time.RFC3339Nano,
		FullTimestamp:   true,
	})

	return logrus.NewEntry(newLogger), nil
}

// main
func main() {
	config, err := InitConfig()
	if err != nil {
		panic(err)
	}

	db, err := InitDatabase(config.Database)
	if err != nil {
		panic(err)
	}

	loggerEntry, err := InitLogger(config.Server)
	if err != nil {
		panic(err)
	}

	router.Router(config, db, loggerEntry)
}
