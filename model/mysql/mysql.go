package mysql

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Info ...
type Info struct {
	Debug    bool
	Hostname string
	Database string
	Username string
	Password string
	Port     int
}

// Connect ...
func (i *Info) Connect() (*gorm.DB, error) {
	connString := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		i.Username,
		i.Password,
		i.Hostname,
		i.Port,
		i.Database,
	)

	db, err := gorm.Open("mysql", connString)
	if err != nil {
		panic(err)
	}
	db.DB().SetMaxIdleConns(0)
	db.LogMode(false)

	return db, err
}
