package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/ykpythemind/funho/config"
	"github.com/ykpythemind/funho/model"
)

func main() {
	config := config.Load()
	db, err := gorm.Open("mysql", config.DBAddr())
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	createUsersRecord(db)
}

func createUsersRecord(db *gorm.DB) {

	users := []string{
		"test1",
		"test2",
		"test3",
	}

	for _, userName := range users {
		user := model.User{Name: userName, Password: "password"}
		db.Create(&user)
	}

}
