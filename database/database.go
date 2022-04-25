package database

import (
	"baal/config"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	maxLifeTime  = 10
	maxOpenConns = 10
	maxIdleConns = 10
)

// New Connent to datebase
func New() (*gorm.DB, error) {
	conf := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	s := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
		config.Secret.Database.Mysql.Username,
		config.Secret.Database.Mysql.Password,
		config.Secret.Database.Mysql.Addr,
		config.Secret.Database.Mysql.Port,
		config.Secret.Database.Mysql.Database,
		config.Secret.Database.Mysql.Charset,
	)

	conn, err := gorm.Open(mysql.Open(s), conf)
	if err != nil {
		return nil, err
	}

	db, err := conn.DB()
	if err != nil {
		// fmt.Fprintln(os.Stderr, err)
		// os.Exit(1)
		return nil, err
	}

	db.SetConnMaxLifetime(time.Duration(maxLifeTime) * time.Second)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)

	return conn, nil
}
