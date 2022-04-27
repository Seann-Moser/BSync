package configuration

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

func connectToDB(conf *Config) error {
	dbConf := mysql.Config{
		AllowNativePasswords: true,
		User:                 conf.DBUser,
		Passwd:               conf.DBPassword,
		Net:                  "tcp",
		Addr:                 conf.DBHost + ":" + conf.DBPort,
		CheckConnLiveness:    true,
	}
	fmt.Println(dbConf.FormatDSN())

	db, err := sqlx.Open("mysql", dbConf.FormatDSN())
	if err != nil {
		return err
	}
	err = db.Ping()
	db.SetMaxOpenConns(10) //max open connections

	for retries := 0; retries < 20 && err != nil; retries++ {
		fmt.Println("Attempting to connect to db: ", retries)
		time.Sleep(5 * time.Second)
		err = db.Ping()
	}
	if err != nil {
		return err
	}
	conf.DB = db

	return nil
}
