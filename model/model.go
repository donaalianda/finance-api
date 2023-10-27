package model

import (
	"finance-api/config/env"
	"finance-api/model/mysql"

	"github.com/jinzhu/gorm"
)

var conn *gorm.DB

// Info ...
type Info struct {
	Config env.Config
}

// Connect ...
func (i *Info) Connect() *gorm.DB {
	var err error

	mysqlInfo := mysql.Info{
		Debug:    i.Config.GetBool("app.debug"),
		Hostname: i.Config.GetString("database.mysql.host"),
		Username: i.Config.GetString("database.mysql.user"),
		Password: i.Config.GetString("database.mysql.password"),
		Database: i.Config.GetString("database.mysql.database"),
		Port:     i.Config.GetInt("database.mysql.port"),
	}

	if conn != nil {
		return conn
	}

	conn, err = mysqlInfo.Connect()

	if err != nil {
		panic(err)
	}

	conn.LogMode(false)

	return conn
}
