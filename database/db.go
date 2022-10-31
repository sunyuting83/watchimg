package database

import (
	"fmt"
	"log"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	Eloquent *gorm.DB
)

// InitDB init db
func InitDB(d string) {
	dbName := strings.Join([]string{d, "data/data.sqlite"}, "/")
	// fmt.Println(dbName)
	Eloquent, _ = gorm.Open("sqlite3", dbName)
	if !Eloquent.HasTable(&ImgList{}) {
		if err := Eloquent.CreateTable(&ImgList{}).Error; err != nil {
			log.Fatal(err)
		}
	}
	if !Eloquent.HasTable(&User{}) {
		if err := Eloquent.CreateTable(&User{}).Error; err != nil {
			log.Fatal(err)
		}
	}
	tableName := "imageData"
	c := 0
	sql := `select * from sqlite_master where name='` + tableName + `' and sql like '%new_status%'`
	row := Eloquent.Raw(sql).Row()
	err := row.Scan(&c)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			sql := `alter table ` + tableName + ` add new_status int default '0'not null`
			fmt.Println(sql)
			Eloquent.Exec(sql)
		}
	}

	sql1 := `select * from sqlite_master where name='` + tableName + `' and sql like '%user_id%'`
	row1 := Eloquent.Raw(sql1).Row()
	err = row1.Scan(&c)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			sql := `alter table ` + tableName + ` add user_id int default '1'`
			fmt.Println(sql)
			Eloquent.Exec(sql)
		}
	}

	sql2 := `select * from sqlite_master where name='` + tableName + `' and sql like '%multiple%'`
	row2 := Eloquent.Raw(sql2).Row()
	err = row2.Scan(&c)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			sql := `alter table ` + tableName + ` add multiple int default '0'`
			fmt.Println(sql)
			Eloquent.Exec(sql)
		}
	}

	sql3 := `select * from sqlite_master where name='` + tableName + `' and sql like '%updatetime%'`
	row3 := Eloquent.Raw(sql3).Row()
	err = row3.Scan(&c)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			sql := `alter table ` + tableName + ` add updatetime int default '0'`
			fmt.Println(sql)
			Eloquent.Exec(sql)
		}
	}

	sql4 := `select * from sqlite_master where name='` + tableName + `' and sql like '%password%'`
	row4 := Eloquent.Raw(sql4).Row()
	err = row4.Scan(&c)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			sql := `alter table ` + tableName + ` add password VARCHAR`
			fmt.Println(sql)
			Eloquent.Exec(sql)
		}
	}

	Eloquent.SingularTable(true)
}
